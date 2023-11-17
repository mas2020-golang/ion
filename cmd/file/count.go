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
	words bool
)

// NewCountCmd represents the wc command
func NewCountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Use: "count <file|pipe|standard-input>",
		Example: `# set the file to read
$ ion count test.txt

# read from the standard input
$ ion ion count < test.txt

# read from the pipe redirection
$ cat test.txt | ion count`,
		Short: "Show the lines or the words of the given input",
		Long: `The count command shows the lines or the words of the given input
The command can read the standard input, a file, the result of a pipe redirection and
return the corresponding words or lines into the standard output`,
		Run: func(cmd *cobra.Command, args []string) {
			counter := file.NewCounter(words)
			c, err := counter.Count(args[0])
			output.CheckErrorAndExit("", "", err)
			fmt.Println(c)
		},
	}

	// flags
	cmd.Flags().BoolVarP(&words, "words", "w", words, "number of words contained in the file/standard input")
	return cmd
}
