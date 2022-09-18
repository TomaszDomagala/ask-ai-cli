package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/TomaszDomagala/ask-ai-cli/pkg/config"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/config/flags"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configSetCmd represents the set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config values in config file using flags",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetGlobalConfig(cmd.Context())
		fs := GetFs(cmd.Context())

		if !configSetCmdConfig.Create.Get() {
			return nil
		}
		if cfg.ConfigFileUsed() != "" {
			return nil
		}

		fileExists, err := afero.Exists(fs, firstDefaultConfigFile)
		if err != nil {
			return fmt.Errorf("failed to check if file %s exists: %w", firstDefaultConfigFile, err)
		}
		if !fileExists {
			err = fs.Mkdir(defaultConfigPaths[0], 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", defaultConfigPaths[0], err)
			}
			_, err = fs.OpenFile(firstDefaultConfigFile, os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", firstDefaultConfigFile, err)
			}
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg := GetGlobalConfig(cmd.Context())
		changes := false

		err = config.Traverse(globalConfig, func(value config.AnyValue) error {
			if value.IsSet() {
				value.SetAny(value.GetAny())
				changes = true
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to traverse config: %w", err)
		}
		if !changes {
			return nil
		}

		err = cfg.WriteConfig()
		if err != nil {
			var notFoundErr viper.ConfigFileNotFoundError
			if errors.As(err, &notFoundErr) {
				return fmt.Errorf("Error: Config file does not exist.\nUse --create to create it in the default path %s\nUse --create and --config to specify custom config file destination\n\n%w", firstDefaultConfigFile, err)
			}
			return fmt.Errorf("failed to write config: %w", err)
		}

		if err != nil {
			return fmt.Errorf("failed to traverse config: %w", err)
		}
		return nil
	},
}

type ConfigSetCmdConfig struct {
	Create flags.Flag[bool]
}

var configSetCmdConfig ConfigSetCmdConfig

func init() {
	configCmd.AddCommand(configSetCmd)

	configSetCmdConfig = ConfigSetCmdConfig{
		Create: flags.Bool(configSetCmd.Flags(), "create", false, "Create config file"),
	}
}
