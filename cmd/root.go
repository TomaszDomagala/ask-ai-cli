package cmd

import (
	"fmt"
	"os"

	"github.com/TomaszDomagala/ask-ai-cli/pkg/openai"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// defaults
const (
	defaultProvider = "openai"
)

// defaults
var (
	defaultConfigPaths = []string{
		"etc/aai/",
		"$HOME/.aai/",
		".",
	}
)

type ConfigFields struct {
	Provider   string
	ApiKey     string
	ConfigPath string
	LogLevel   string
}

var cfgFields = ConfigFields{
	Provider:   "provider",
	ApiKey:     "apikey",
	ConfigPath: "config",
	LogLevel:   "loglevel",
}

var providerConfig *viper.Viper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   `aai [query]`,
	Short: "Ask AI to suggest a command",
	Long: `ask-ai-cli (aai) is a command line tool that helps you to find a command you need.
It uses AI to suggest a command based on your query.

Example:
    aai "show files with size greater than 1MB"
`,

	Args: cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Init config
		var err error

		configPath, err := cmd.Flags().GetString(cfgFields.ConfigPath)
		if err != nil {
			return fmt.Errorf("failed to get config path flag: %w", err)
		}

		if configPath != "" {
			viper.SetConfigFile(configPath)
		} else {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			for _, path := range defaultConfigPaths {
				viper.AddConfigPath(path)
			}
		}

		err = viper.ReadInConfig()
		if err != nil {
			// Ignore config file not found error, we will use defaults
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				// Another error occurred
				return fmt.Errorf("failed to read config: %v", err)
			}
		}

		if err = viper.BindPFlag(cfgFields.Provider, cmd.Flags().Lookup(cfgFields.Provider)); err != nil {
			return fmt.Errorf("failed to bind provider flag: %w", err)
		}
		provider := viper.GetString(cfgFields.Provider)
		providerConfig = viper.Sub(fmt.Sprintf("providers.%s", provider))

		err = providerConfig.BindPFlags(cmd.PersistentFlags())
		if err != nil {
			return fmt.Errorf("failed to bind flags: %v", err)
		}

		loglevel, err := zerolog.ParseLevel(viper.GetString(cfgFields.LogLevel))
		if err != nil {
			return fmt.Errorf("failed to parse log level: %v", err)
		}
		// From now on, we can use loggers
		zerolog.SetGlobalLevel(loglevel)

		log.Debug().Msgf("using config file: %v", viper.ConfigFileUsed())

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		client := openai.NewClient(providerConfig.GetString(cfgFields.ApiKey), openai.DefaultCompletionConfig)

		response, err := client.Suggest(query)
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

	err := rootCmd.Execute()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to execute root command")
		os.Exit(1)
	}
}

func init() {
	// Local flags
	rootCmd.PersistentFlags().String("config", "", fmt.Sprintf("config file - if not provided, the following paths will be checked: %v", fmt.Sprint(defaultConfigPaths)))

	// Persistent flags
	rootCmd.PersistentFlags().String("apikey", "", "api key for the provider")
	rootCmd.PersistentFlags().StringP("provider", "p", defaultProvider, "provider to use")
	rootCmd.PersistentFlags().String("loglevel", "DISABLED", "zerolog log level")
}
