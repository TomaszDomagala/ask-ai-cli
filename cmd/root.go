package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TomaszDomagala/ask-ai-cli/pkg/config"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/config/flags"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/openai"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// defaultConfigPaths is a list of paths to check for config files
	defaultConfigPaths = []string{
		"/etc/aai/",
		"$HOME/.aai/",
		".",
	}
	// configFileName is the name of the config file, without extension
	configFileName = "config"
	// configFileType is the type of the config file
	configFileType = "yaml"
	// firstDefaultConfigFile is the path to the first default config file
	firstDefaultConfigFile = filepath.Join(defaultConfigPaths[0], configFileName+"."+configFileType)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   `aai [query]`,
	Short: "Ask AI to suggest a command",
	Long: `
ask-ai-cli (aai)
A command line tool that helps you find a command you need.
It uses AI to suggest a command based on your query.

It is advised to not use suggestions blindly,
but rather to read the documentation and understand
what the command does before running it.

Example:
    aai "show files with size greater than 1MB"
`,

	Args: cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg := GetGlobalConfig(cmd.Context())

		if globalConfig.ConfigFile.Get() != "" {
			cfg.SetConfigFile(globalConfig.ConfigFile.Get())
		} else {
			cfg.SetConfigName(configFileName)
			cfg.SetConfigType(configFileType)
			for _, path := range defaultConfigPaths {
				cfg.AddConfigPath(path)
			}
		}

		err = cfg.ReadInConfig()
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		// Setup config by attaching viper to the config struct
		err = config.Attach(cfg, &globalConfig)

		if err != nil {
			return fmt.Errorf("failed to attach config: %w", err)
		}

		// Setup logs
		loglevel, err := zerolog.ParseLevel(globalConfig.LogLevel.Get())
		if err != nil {
			return fmt.Errorf("failed to parse log level: %v", err)
		}
		// From now on, we can use loggers
		zerolog.SetGlobalLevel(loglevel)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var suggester Suggester

		switch globalConfig.Provider.Get() {
		case "openai":
			var openaiCfg openai.Config
			err := config.Decode(globalConfig, &openaiCfg)
			if err != nil {
				return fmt.Errorf("failed to decode config: %w", err)
			}
			suggester = openai.NewClient(openaiCfg)
		default:
			return fmt.Errorf("unknown provider: %v", globalConfig.Provider.Get())
		}

		query := args[0]
		response, err := suggester.Suggest(query)
		if err != nil {
			return fmt.Errorf("failed to suggest a command: %w", err)
		}
		fmt.Println(response)

		return nil
	},
}

type Suggester interface {
	// Suggest returns a suggestion for a given query.
	Suggest(query string) (string, error)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	global := viper.New()
	fs := afero.NewOsFs()

	ctx := context.Background()
	ctx = SetGlobalConfig(ctx, global)
	ctx = SetFs(ctx, fs)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type GlobalConfig struct {
	ConfigFile flags.Flag[string]
	Provider   config.Value[string]
	LogLevel   config.Value[string]

	OpenAiConfig
}

type OpenAiConfig struct {
	// ApiKey for OpenAI API
	ApiKey config.Value[string]

	// OpenAi request settings.
	// Description: https://beta.openai.com/docs/api-reference/completions/create

	Model            config.Value[string] // Model
	Temperature      config.Value[float64]
	MaxTokens        config.Value[int]
	TopP             config.Value[float64]
	FrequencyPenalty config.Value[float64]
	PresencePenalty  config.Value[float64]
}

var globalConfig GlobalConfig

func init() {
	// define global config
	globalConfig = GlobalConfig{
		ConfigFile: flags.String(rootCmd.PersistentFlags(), "config", "", fmt.Sprintf("config file - if not provided, the following paths will be checked: %v", fmt.Sprint(defaultConfigPaths))),
		Provider:   config.String("provider", config.WithFlag(rootCmd.PersistentFlags(), "provider", "openai", "provider to use for suggestions")),
		LogLevel:   config.String("loglevel", config.WithFlag(rootCmd.PersistentFlags(), "loglevel", "disabled", "log level (zerolog)")),

		OpenAiConfig: OpenAiConfig{
			ApiKey:           config.String("openai.apikey", config.WithFlag(rootCmd.PersistentFlags(), "openai-apikey", "", "openai api key")),
			Model:            config.String("openai.model", config.WithFlag(rootCmd.PersistentFlags(), "openai-model", "text-davinci-002", "openai model to use for completion")),
			Temperature:      config.Float64("openai.temperature", config.WithFlag(rootCmd.PersistentFlags(), "openai-temperature", 0.2, "temperature")),
			MaxTokens:        config.Int("openai.maxtokens", config.WithFlag(rootCmd.PersistentFlags(), "openai-maxtokens", 100, "max tokens")),
			TopP:             config.Float64("openai.topp", config.WithFlag(rootCmd.PersistentFlags(), "openai-topp", 1.0, "top p")),
			FrequencyPenalty: config.Float64("openai.frequencypenalty", config.WithFlag(rootCmd.PersistentFlags(), "openai-frequencypenalty", 0.0, "frequency penalty")),
			PresencePenalty:  config.Float64("openai.presencepenalty", config.WithFlag(rootCmd.PersistentFlags(), "openai-presencepenalty", 0.0, "presence penalty")),
		},
	}

}
