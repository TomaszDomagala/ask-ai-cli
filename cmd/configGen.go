package cmd

import (
	"github.com/spf13/cobra"
)

// configGenCmd represents the configGen command
var configGenCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate config",

	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	configCmd.AddCommand(configGenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configGenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configGenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
