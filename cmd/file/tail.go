/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com

*/
package file

import (
	"fmt"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var (
	rows int = 10
)

// NewTailCmd represents the tail command
func NewTailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tail <file|pipe|standard-input>",
		Example: `$ ion tail --rows 10 test.txt
// read from the standard input
$ ion tail -r 10 < test.txt
// read from the pipe
$ cat test.txt | ion tail --rows 10`,
		Short: "Show the n latest rows from the input",
		Long: `The tail command shows the n latest rows from the given input.
The command can read the standard input, a file, the result of a pipe redirection and
return the corresponding rows.`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				f   *os.File = utils.GetBytesFromPipe()
				err error
			)
			if f == nil {
				if len(args) == 0 {
					utils.Check(fmt.Errorf("no file argument"))
				}
				// load the file into the buffer
				f, err = os.Open(args[0])
				utils.Check(err)
			}
			lines, err := getLines(f, rows)
			utils.Check(err)
			for _, l := range lines {
				fmt.Print(l)
			}
		},
	}

	// flags
	cmd.Flags().IntVarP(&rows, "rows", "r", rows, "number of rows to show starting from the end of the file")
	cmd.MarkFlagRequired("rows")
	return cmd
}

// getLines returns the n lines in buf
func getLines(f *os.File, n int) (lines []string, err error) {
	var (
		b        []byte
		l        [][]byte
		i, found int
		pos      int
	)
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	//fmt.Printf("file content (bytes): %v\n", buf)
	utils.Check(err)
	pos = len(buf) - 1

	for i < n && pos >= 0 {
		// \n
		if buf[pos] == 10 {
			found++
			if found > 1 {
				l = append(l, reverseSlice(b))
				b = make([]byte, 0)
				i++
				found = 0
				continue
			}
		}
		b = append(b, buf[pos])
		// in case it is the beginning add the last buffer
		if pos == 0{
			l = append(l, reverseSlice(b))
		}
		pos--
	}

	for i := len(l) - 1; i >= 0; i-- {
		lines = append(lines, string(l[i]))
	}
	return
}

func reverseSlice(b []byte) []byte {
	var (
		r []byte
	)
	if len(b) == 1 {
		return b
	}
	for i := len(b) - 1; i >= 0; i-- {
		r = append(r, b[i])
	}
	return r
}
