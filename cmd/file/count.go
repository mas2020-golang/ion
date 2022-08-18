/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com

*/
package file

import (
	"fmt"
	"io"
	"os"

	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

var (
	words bool
)

// NewCountCmd represents the wc command
func NewCountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "count <file|pipe|standard-input>",
		Example: `# point the file to read
$ ion wc test.txt
# read from the standard input
$ ion ion wc < test.txt 
# read from the pipe
$ cat test.txt | ion wc`,
		Short: "Show the lines or the words of the given input",
		Long: `The wc command shows the lines or the words of the given input
The command can read the standard input, a file, the result of a pipe redirection and
return the corresponding words or lines`,
		Run: func(cmd *cobra.Command, args []string) {
			c, err := wc(args)
			utils.Check(err)
			fmt.Println(c)
		},
	}

	// flags
	cmd.Flags().BoolVarP(&words, "words", "w", words, "number of words contained in the file/standard input")
	return cmd
}

func wc(args []string) (count int, err error) {
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
	return getCount(f)
}

// getLines returns the number of lines in f or the words depending to
// the args passed to the command
func getCount(f *os.File) (count int, err error) {
	var (
		b     []byte
		space bool
	)

	defer func() {
		err := f.Close()
		utils.Check(err)
	}()
	if err != nil {
		return -1, fmt.Errorf("error on file stat: '%v'", err)
	}
	space = true // first iteration starts with a space
	// read 1024 each time
	b = make([]byte, 1024)
	for {
		n, err := f.Read(b)
		if n > 0 {
			for i := 0; i < n; i++ {
				// case --word
				if words {
					if !isSpace(b[i]) && space {
						count++
					}
					// save if the current byte is a space to check for the
					// next iteration
					space = isSpace(b[i])
				} else {
					// case line count
					if b[i] == 10 {
						count++
					}
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return -1, fmt.Errorf("read %d bytes: %v", n, err)
		}
	}
	return count, nil
}

// isSpace returns true if the byte passed is considered a char delimiter.
// The following chars are considered space: \n \t space.
func isSpace(b byte) bool {
	return b == 9 || b == 32 || b == 10
}
