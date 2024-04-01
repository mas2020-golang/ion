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
			fmt.Printf(`%-12s%s
%-12s%s
%-12s%s
%-12s%s`, 
"Version:",utils.Version, 
"Git commit:", utils.GitCommit, 
"Built on: ", utils.BuildDate, 
"Go version:", utils.GoVersion)
			fmt.Println()
		},
	}
	return newCmd
}
