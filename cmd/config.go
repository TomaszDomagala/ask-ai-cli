package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display or change current configuration",
}

func init() {
	rootCmd.AddCommand(configCmd)
}
