package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/TomaszDomagala/ask-ai-cli/pkg/config"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/openai"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// defaults
var (
	defaultConfigPaths = []string{
		"etc/aai/",
		"$HOME/.aai/",
		".",
	}
)

type GlobalConfig struct {
	Provider string `name:"provider" value:"openai" usage:"provider to use" validate:"required"`
	LogLevel string `name:"loglevel" value:"disabled" usage:"log level (zerolog)"`
}

var providerConfig *viper.Viper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   `aai [query]`,
	Short: "Ask AI to suggest a command",
	Long: `
ask-ai-cli (aai)
A command line tool that helps you to find a command you need.
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
		var conf GlobalConfig

		globalCfg := GetGlobalConfig(cmd.Context())

		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("failed to get config path flag: %w", err)
		}

		if configPath != "" {
			globalCfg.SetConfigFile(configPath)
		} else {
			globalCfg.SetConfigName("config")
			globalCfg.SetConfigType("yaml")
			for _, path := range defaultConfigPaths {
				globalCfg.AddConfigPath(path)
			}
		}

		err = globalCfg.ReadInConfig()
		if err != nil {
			// Ignore config file not found error, we will use defaults
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				// Another error occurred
				return fmt.Errorf("failed to read config: %v", err)
			}
		}

		if err = globalCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}
		if err = config.FillWithConfig(globalCfg, &conf); err != nil {
			return fmt.Errorf("failed to fill config: %w", err)
		}

		providerConfig = globalCfg.Sub(fmt.Sprintf("providers.%s", conf.Provider))
		if err = config.BindProviderFlags(providerConfig, cmd.PersistentFlags(), conf.Provider); err != nil {
			return fmt.Errorf("failed to bind provider flags: %w", err)
		}
		cmd.SetContext(SetProviderConfig(cmd.Context(), providerConfig))

		loglevel, err := zerolog.ParseLevel(conf.LogLevel)
		if err != nil {
			return fmt.Errorf("failed to parse log level: %v", err)
		}
		// From now on, we can use loggers
		zerolog.SetGlobalLevel(loglevel)

		log.Debug().Msgf("using config file: %v", viper.ConfigFileUsed())

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var globalCfg GlobalConfig
		var suggester Suggester

		ctx := cmd.Context()

		if err = config.FillWithConfig(GetGlobalConfig(ctx), &globalCfg); err != nil {
			return fmt.Errorf("failed to fill global config: %w", err)
		}

		switch globalCfg.Provider {
		case "openai":
			var openaiCfg openai.Config
			if err = config.FillWithConfig(GetProviderConfig(ctx), &openaiCfg); err != nil {
				return fmt.Errorf("failed to fill openai config: %w", err)
			}
			suggester = openai.NewClient(openaiCfg)
		default:
			return fmt.Errorf("unknown provider: %v", globalCfg.Provider)
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

	ctx := context.Background()
	ctx = SetGlobalConfig(ctx, global)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags independent of provider
	rootCmd.PersistentFlags().String("config", "", fmt.Sprintf("config file - if not provided, the following paths will be checked: %v", fmt.Sprint(defaultConfigPaths)))
	config.FlagsFromStruct(rootCmd.PersistentFlags(), GlobalConfig{}, "")

	// Provider specific flags
	config.FlagsFromStruct(rootCmd.PersistentFlags(), openai.Config{}, "openai")
}
