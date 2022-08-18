/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
Special thanks to this web article:
https://www.socketloop.com/tutorials/golang-secure-file-deletion-with-wipe-example
*/
package file

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	out "github.com/mas2020-golang/goutils/output"
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
		Short: "Removes the files or folders given as an input",
		Long: `The command removes the files or folders given as an input. Use a space as a separator between the files (or folders).
The command returns 0 in case of success and 1 in case something went wrong (plus an error onto the standard error).`,
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				if verbose {
					out.SubActivity(fmt.Sprintf("read path for %s", arg))
				}
				if err := readPath(arg); err != nil {
					out.Error("", err.Error())
					os.Exit(1)
				} else {
					out.InfoBox(fmt.Sprintf("removed %d files from %s (searched on %d dirs)", files, arg, dirs))
					out.InfoBox(fmt.Sprintf("total space retrieved is %d KB", spaceFree/1024))
				}
			}
		},
	}

	// flags
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode active")
	cmd.Flags().BoolVarP(&dryrun, "dry-run", "d", false, "if dry-run nothing is deleted")
	cmd.Flags().BoolVarP(&deep, "deep-clean", "c", false, "run a deep clean to remove the files in depth (takes longer but is more secure)")
	return cmd
}

// readPath reads the path and in case it is a file delete it, otherwise search into the nested
// folders
func readPath(path string) error {
	//var symbol string = elemChar
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat error for the path %s, error: %v", path, err)
	}

	if fi.IsDir() {
		dirs++
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("ReadDir error for the path %s, error: %v", path, err)
		}

		for _, fi := range fis {
			if err := readPath(filepath.Join(path, fi.Name())); err != nil {
				return err
			}
		}
		deleteFolder(path)
	} else {
		spaceFree += fi.Size()
		if err := deleteFile(path); err != nil {
			return err
		}
	}
	return nil
}

// deleteFile deletes the given file
func deleteFile(f string) error {
	files++
	if dryrun {
		message(out.RedS(fmt.Sprintf("the file %s would be deleted", f)))
		return nil
	}

	// file deletion

	if deep {
		if err := deepClean(f); err != nil {
			return fmt.Errorf("error deleting the file %s, error: %v", f, err)
		}
		message(out.RedS(fmt.Sprintf("the file %s has been wiped out", f)))
	}

	if err := os.Remove(f); err != nil {
		return fmt.Errorf("error deleting the file %s, error: %v", f, err)
	}
	message(out.RedS(fmt.Sprintf("the file %s has been deleted", f)))

	return nil
}

// deleteFolder deletes the given folder
func deleteFolder(f string) error {
	if dryrun {
		message(out.BlueS(fmt.Sprintf("the folder %s would be deleted", f)))
		return nil
	}

	// folder deletion
	if err := os.Remove(f); err != nil {
		return err
	}

	message(out.BlueS(fmt.Sprintf("deleted the folder %s", f)))
	return nil
}

func message(t string) {
	if verbose || dryrun {
		fmt.Println(t)
	}
}

// deepClean removes the file in depth
func deepClean(path string) error {
	// open the file
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	defer func(*os.File) {
		if err := f.Close(); err != nil {
			out.Error("deepClean", err.Error())
		}
	}(f)

	if err != nil {
		return fmt.Errorf("error opening the file %s: %s", path, err.Error())
	}

	// stat the file
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error opening the file %s: %s", path, err.Error())
	}

	// slice for the size of the file
	var fSize int64 = fi.Size()
	const chunk = 2 * (1 << 20) // 2 MB

	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fSize) / float64(chunk)))

	lastPosition := 0

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(chunk, float64(fSize-int64(i*chunk))))
		partZeroBytes := make([]byte, partSize)

		// fill out the part with zero value
		copy(partZeroBytes[:], "0")

		// over write every byte in the chunk with 0
		_, err := f.WriteAt([]byte(partZeroBytes), int64(lastPosition))

		if err != nil {
			return fmt.Errorf("error wiping 0 to the file %s: %s", path, err.Error())
		}

		// update last written position
		lastPosition = lastPosition + partSize
	}

	return nil
}
