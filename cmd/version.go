package cmd

import (
	"fmt"

	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}

func NewVersionCmd() *cobra.Command {
	// Show sub commands
	newCmd := &cobra.Command{
		Use: "version",
		//Args:  cobra.ExactArgs(1),
		Short: "Show the application version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			//utils.BuildDate = "2022-10-25"
			fmt.Printf(`%-12s%s
Git commit: %s
%-12s%s`, "Version:", utils.Version, utils.GitCommit, "Built on:", utils.BuildDate)
			fmt.Println()
		},
	}
	return newCmd
}
