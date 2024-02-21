/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"fmt"

	out "github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/file"
	"github.com/spf13/cobra"
)

var (
	sliceBytes, sliceChars, sliceCols string
	delimiter                         = ""
)

// NewTailCmd represents the tail command
func NewSliceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.MinimumNArgs(1),
		Use:  "slice [flags] <file|pipe|standard-input>",
		Example: `$ ion slice -b 1:3 test.txt

-- read from the standard input
$ ion slice -b 1:3 < test.txt

-- read from the pipe
$ cat test.txt | ion slice -b 1:3

-- extract the bytes from start to end expressed as start:end:
$ ion slice -b 1:3 test.txt

You can specify start: to start from that point to the end of the input line or simply start to get a single char.
In case you need separated char you can use comma, for example: 1,2,6 gives you the corresponding bytes as string.
If the single byte is not an ascii char, specify the colon to get the right char (usually for UTF8 econded files).

-- extract the chars from the start to end expressed as start:end:
$ ion slice -c 1:3 test.txt

You can specify a single char or a chunk of chars or the beginning char till the end as: -c 10:.

-- extract by columns using a delimiter expressed by the -d option:
$ ion slice -f 3 -d " "

In the example above we are cutting by space.
You can specify the intervals as already seen.
`,
		Short: "Slice the provided input",
		Long: `The slice command slices a line and extracts the text. The input can be cut
		by byte position, chars or fields.`,
		Run: func(cmd *cobra.Command, args []string) {
			slice := file.NewSlice()
			s, err := slice.Slice(args[0], sliceBytes, sliceChars, sliceCols)
			out.CheckErrorAndExit("", "", err)
			fmt.Println(s)
		},
	}

	// flags
	cmd.Flags().StringVarP(&sliceBytes, "bytes", "b", sliceBytes, "interval of bytes to slice (start:end)")
	cmd.Flags().StringVarP(&sliceChars, "chars", "c", sliceChars, "interval of chars to slice (start:end)")
	cmd.Flags().StringVarP(&sliceCols, "fields", "f", sliceCols, "interval of fields to slice (start:end)")
	cmd.Flags().StringVarP(&delimiter, "delimiter", "d", delimiter, "delimiter for determine the fields slicing (space is the default)")

	return cmd
}
