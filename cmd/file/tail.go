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
		Use: "tail <file|pipe|standard-input>",
		Example: `$ ion tail --rows 10 test.txt
// read from the standard input
$ ion tail -r 10 < test.txt
// read from the pipe
$ cat test.txt | ion tail --rows 10`,
		Short: "Show the n latest rows from the given input",
		Long: `The tail command shows the n latest rows from the given input.
The command can read the standard input, a file, the result of a pipe redirection and
return the corresponding rows.`,
		Run: func(cmd *cobra.Command, args []string) {
			lines, err := tail(args, rows)
			utils.Check(err)
			for _, l := range lines {
				fmt.Print(l)
			}
		},
	}

	// flags
	cmd.Flags().IntVarP(&rows, "rows", "r", rows, "number of rows to show starting from the end of the file")
	err := cmd.MarkFlagRequired("rows")
	utils.Check(err)
	return cmd
}

func tail(args []string, r int) (lines []string, err error) {
	var (
		f *os.File = utils.GetBytesFromPipe()
	)
	if f == nil {
		if len(args) == 0 {
			utils.Check(fmt.Errorf("no file argument"))
		}
		// load the file into the buffer
		f, err = os.Open(args[0])
		utils.Check(err)
	}
	return getLines(f, r)
}

// getLines returns the n lines in buf
func getLines(f *os.File, n int) (lines []string, err error) {
	var (
		b      []byte
		l      [][]byte
		pos, i int
	)

	defer func() {
		err := f.Close()
		utils.Check(err)
	}()

	buf, err := ioutil.ReadAll(f)
	//fmt.Printf("file content (bytes): %v\n", buf)
	utils.Check(err)
	pos = len(buf) - 1

	for i < n && pos >= 0 {
		// \n
		if buf[pos] == 10 {
			if len(b) == 0 {
				b = append(b, buf[pos])
				pos--
			} else {
				// create a new line
				l = append(l, reverseSlice(b))
				b = make([]byte, 0)
				i++
			}
		} else {
			b = append(b, buf[pos])
			pos--
		}
		// edge case: beginning of the file
		if pos == -1 && len(b) > 0 {
			l = append(l, reverseSlice(b))
		}
	}
	// create the slice of lines
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
