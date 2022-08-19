/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mas2020-golang/ion/cmd/file"
	"github.com/mas2020-golang/ion/cmd/security"
	"github.com/mas2020-golang/ion/packages/utils"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ion",
	Short: "Ion is your swiss knife for having with you a lot of useful commands",
	Long: `Ion is an all-in-one application to sum up a lot of useful tools in a single command.
The swiss knife for every SysAdmin/DevOps!. You can use the ion commands as you do with pipes,
standard input/output and a lot of other daily basis activities.
`,
	Example: `// to search some content into a file
$ ion search --regexp '(temp)' --color test.txt

// tail the last 10 rows
$ ion tail --rows 10 test.txt`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ion.yaml)")
	rootCmd.SetHelpFunc(utils.GetHelpFunction(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}`))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add the other commands
	rootCmd.AddCommand(file.NewTailCmd())
	rootCmd.AddCommand(file.NewTreeCmd())
	rootCmd.AddCommand(security.NewCryptoCmd())
	rootCmd.AddCommand(file.NewCountCmd())
	rootCmd.AddCommand(file.NewRmCmd())
}

// getHelp returns the help for the command
func getHelp() func(*cobra.Command, []string) {
	return func(*cobra.Command, []string) {
		fmt.Print(`Ion is an all-in-one application to sum up a lot of useful tools in a single command.
The swiss knife for every SysAdmin/DevOps! You can use the ion commands as you do with pipes,
standard input/output and a lot of other daily basis activities.

File Commands:
  count      Count the lines or the words of the given input
  rm         Remove the files or folders given as an input
  tail       Show the n latest rows from the given input
  tree       Show the file system in a tree graphical representation
  
Security Commands:
  encrypt    Encrypt file/s or folder/s using the AES algorithm
  decrypt    Decrypt file/s or folder/s using the AES algorithm

Flags:
  -h, --help     help for ion
  -t, --toggle   Help message for toggle

Use "ion [command] --help" for more information about a command.
`)
	}
}
