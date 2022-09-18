package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configViewCmd represents the configView command
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Display config",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		ctx := cmd.Context()
		globalConfig := GetGlobalConfig(ctx)

		if globalConfig.ConfigFileUsed() == "" {
			return nil
		}

		readonlyConfig := viper.New()
		readonlyConfig.SetConfigFile(globalConfig.ConfigFileUsed())
		if err = readonlyConfig.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		fmt.Print(configToString(readonlyConfig))
		return nil
	},
}

func init() {
	configCmd.AddCommand(configViewCmd)
}

// configToString returns a string representation of the config
func configToString(cfg *viper.Viper) string {
	keys := cfg.AllKeys()
	sort.Strings(keys)

	return configLevelToString(keys, cfg, 0)
}

func configLevelToString(keys []string, cfg *viper.Viper, indent int) string {
	out := ""
	indentStr := strings.Repeat(" ", indent)

	keysByPrefix := make(map[string][]string)
	for _, key := range keys {
		if strings.Contains(key, ".") {
			prefix := strings.Split(key, ".")[0]
			keysByPrefix[prefix] = append(keysByPrefix[prefix], key)
		} else {
			keysByPrefix[""] = append(keysByPrefix[""], key)
		}
	}

	prefixes := make([]string, 0, len(keysByPrefix))
	for prefix := range keysByPrefix {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)

	for _, prefix := range prefixes {
		keys := keysByPrefix[prefix]
		if prefix == "" {
			for _, key := range keys {
				out += fmt.Sprintf("%s%s: %v\n", indentStr, key, cfg.Get(key))
			}
		} else {
			out += fmt.Sprintf("%s%s:\n", indentStr, prefix)
			trimmedKeys := make([]string, 0, len(keys))
			for _, key := range keys {
				trimmedKeys = append(trimmedKeys, strings.TrimPrefix(key, prefix+"."))
			}

			out += configLevelToString(trimmedKeys, cfg.Sub(prefix), indent+2)
		}
	}
	return out
}
