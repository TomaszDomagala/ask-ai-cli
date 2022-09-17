package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// configSetCmd represents the set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("set called")

	},
}

func init() {
	configCmd.AddCommand(configSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configSetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configSetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
