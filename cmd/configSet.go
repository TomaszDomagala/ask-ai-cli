package cmd

import (
	"errors"
	"fmt"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/config"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/config/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// configSetCmd represents the set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config values in config file using flags",

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
			// nothing to do
			return nil
		}

		// Write config to file
		if configSetCmdConfig.File.Changed() {
			// Write to file specified by flag, create directories if needed
			file := configSetCmdConfig.File.Get()
			fs := GetFs(cmd.Context())
			if err = fs.MkdirAll(filepath.Dir(file), 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(file), err)
			}
			err = cfg.WriteConfigAs(configSetCmdConfig.File.Get())
		} else if cfg.ConfigFileUsed() != "" {
			err = cfg.WriteConfig()
		} else {
			return errors.New("no config file to write to")
		}

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
	File flags.Flag[string]
}

var configSetCmdConfig ConfigSetCmdConfig

func init() {
	var err error
	configCmd.AddCommand(configSetCmd)

	configSetCmdConfig = ConfigSetCmdConfig{
		File: flags.StringP(configSetCmd.Flags(), "file", "f", "", "file to write config to"),
	}
	if err = configSetCmd.MarkFlagFilename("file", "yaml"); err != nil {
		panic(err)
	}
}
