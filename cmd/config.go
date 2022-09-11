package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
	"strings"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display or change current configuration",

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		ctx := cmd.Context()
		globalConfig := GetGlobalConfig(ctx)

		if globalConfig.ConfigFileUsed() == "" {
			fmt.Println("No config file found")
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

// configToString returns a string representation of the config
// viper config is
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

func init() {
	rootCmd.AddCommand(configCmd)
}
