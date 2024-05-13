/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/file"
	"github.com/spf13/cobra"
)

var (
	replaceSub, replacePattern string
	replaceVerbose, replaceAll bool
)

// NewReplaceCmd represents the replace command
func NewReplaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		//Args: cobra.MinimumNArgs(1),
		Use:  "replace [flags] <file|pipe|standard-input>",
		Args: cobra.MinimumNArgs(1),
		Example: `$ ion replace ...

-- ...
$ ...
`,
		Short: "Replace the content of the given file",
		Long: `The replace command replaces the content of a given file with the specified patter. Replace command
also contains some other options to alter the input file in a bunch of different ways.
Take a look at the example section for more explanations.`,
		Run: func(cmd *cobra.Command, args []string) {
			replacer := file.NewReplacer(replaceVerbose, replaceAll, replacePattern, replaceSub)
			err := replacer.Replace(args[0])
			output.CheckErrorAndExit("file.NewReplaceCmd", "replacing error", err)
		},
	}
	cmd.GroupID = "file"

	// flags
	cmd.Flags().StringVarP(&replaceSub, "substitution", "s", replaceSub, "substitution pattern")
	cmd.Flags().StringVarP(&replacePattern, "pattern", "p", replacePattern, "regexp pattern")
	cmd.Flags().BoolVarP(&replaceVerbose, "verbose", "v", false, "verbosity mode")
	cmd.Flags().BoolVarP(&replaceAll, "all", "a", false, "if all the search is done on the entire file, otherwise only the first occurrence is taken")

	return cmd
}
