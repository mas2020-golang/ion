/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
Special thanks to this web article:
https://www.socketloop.com/tutorials/golang-secure-file-deletion-with-wipe-example
*/
package file

import (
	"fmt"
	"os"

	out "github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/file"
	"github.com/spf13/cobra"
)

var (
	verbose, dryrun, deep bool
	spaceFree             int64
)

func NewRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.MinimumNArgs(1),
		Use:  "rm <path>[...path]",
		Example: `# point the file to delete
$ ion rm test.txt
# delete a folder and all the subfolders
$ ion rm folder
# delete more objects
$ ion delete folder file1 file2`,
		Short: "Removes specific files or folders",
		Long: `The command removes the files or folders given as an input. Use a space as a separator between the files (or folders).
The command returns 0 in case of success and 1 in case something went wrong (plus an error onto the standard error).`,
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				if verbose {
					out.SubActivity(fmt.Sprintf("read path for %s", arg))
				}
				// new eraser object
				eraser := file.NewEraser(dryrun, verbose, deep)
				if err := eraser.Delete(arg); err != nil {
					out.Error("", err.Error())
					os.Exit(1)
				} else {
					out.InfoBox(fmt.Sprintf("removed %d files from '%s' (searched on %d dirs)", eraser.Files, arg, eraser.Dirs))
					out.InfoBox(fmt.Sprintf("total space retrieved is %d KB", eraser.SpaceFree/1024))
				}
			}
		},
	}
	cmd.GroupID = "sec"
	// flags
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode active")
	cmd.Flags().BoolVarP(&dryrun, "dry-run", "d", false, "if dry-run nothing is deleted")
	cmd.Flags().BoolVarP(&deep, "deep-clean", "c", false, "run a deep clean to remove the files in depth (takes longer but is more secure)")
	return cmd
}
