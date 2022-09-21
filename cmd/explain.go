package cmd

import (
	"errors"
	"fmt"

	"github.com/TomaszDomagala/ask-ai-cli/pkg/config"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/errs"
	"github.com/TomaszDomagala/ask-ai-cli/pkg/openai"

	"github.com/spf13/cobra"
)

var (
	// errNoCommandArg is returned when no command argument is provided
	errNoCommandArg = errors.New("no command argument provided")
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain <command>",
	Short: "Explain provided command",
	Long: `This command will explain the provided command.

Example:
	$ aai explain "ls -l"
	List the contents of the current directory in long format
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errs.New(errNoCommandArg, "Please provide a command to explain")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var explainer Explainer

		switch globalConfig.Provider.Get() {
		case "openai":
			var openaiCfg openai.Config
			err := config.Decode(globalConfig, &openaiCfg)
			if err != nil {
				return fmt.Errorf("failed to decode config: %w", err)
			}
			explainer = openai.NewClient(openaiCfg)
		default:
			return fmt.Errorf("unknown provider: %v", globalConfig.Provider.Get())
		}

		command := args[0]
		response, err := explainer.Explain(command)
		if err != nil {
			return fmt.Errorf("failed to explain a command: %w", err)
		}
		fmt.Println(response)

		return nil
	},
}

type Explainer interface {
	// Explain returns an explanation for a given command.
	Explain(command string) (string, error)
}

func init() {
	rootCmd.AddCommand(explainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
