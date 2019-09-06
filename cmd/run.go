package cmd

import (
	"github.com/spf13/cobra"
	"smart.com/weixin/smart/logp"
	"smart.com/weixin/smart/modules/infra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run main function",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

//Run default cmd
func Run() error {
	logger := logp.NewLogger(ModuleName)
	logger.Info("Enter Run")

	infra.Run()
	return nil
}
