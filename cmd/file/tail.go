/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"fmt"

	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/file"
	"github.com/spf13/cobra"
)

var (
	rows int = 10
)

// NewTailCmd represents the tail command
func NewTailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tail <file|pipe|standard-input>",
		Example: `$ ion tail --rows 10 test.txt

// read from the standard input
$ ion tail -r 10 < test.txt

// read from the pipe
$ cat test.txt | ion tail --rows 10`,
		Short: "Show the n latest rows from the given input",
		Long: `The tail command shows the n latest rows from the given input.
The command can read the standard input or a given file and returns the corresponding rows.
If the --rows is not given, the command returns the last 10 rows.`,
		Run: func(cmd *cobra.Command, args []string) {
			t := file.NewTail()
			lines, err := t.Tail(args, rows)
			output.CheckErrorAndExit("", "", err)
			for _, l := range lines {
				fmt.Print(l)
			}
		},
	}

	cmd.GroupID = "file"

	// flags
	cmd.Flags().IntVarP(&rows, "rows", "r", rows, "number of rows to show starting from the end of the file")
	//err := cmd.MarkFlagRequired("rows")
	//utils.Check(err)
	return cmd
}
