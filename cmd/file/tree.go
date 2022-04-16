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
	"path/filepath"
	"strings"
)

var (
	levels, dirs, files int = 5, 0, 0
	summary, colorize   bool
)

// NewTreeCmd represents the tree command
func NewTreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tree <FOLDER>",
		Example: `$ ion tree --levels 2 test.txt
$ ion tree -l 10 .`,
		Short: "Show the file system in a tree graphical representation",
		Long: `Show the file system structure in a tree graphical representation. The tree shows the folder
and the files as a hierarchy.`,
		Run: func(cmd *cobra.Command, args []string) {
			carriage := ""
			for i, arg := range args {
				if i > 0 {
					carriage = "\n"
				}
				fmt.Printf("%s%s\n", carriage, utils.BoldS(arg))
				err := tree(arg, 0, "")
				if err != nil {
					fmt.Println(err)
					continue
				}
				if summary {
					fmt.Printf("\ndirectories: %d, files: %d\n", dirs, files)
				}
				dirs, files = 0, 0
			}
			cmd.OutOrStdout()
		},
	}

	// flags
	cmd.Flags().IntVarP(&levels, "levels", "l", levels, "number of levels for the hierarchy")
	cmd.Flags().BoolVarP(&summary, "summary", "s", false, "for printing the numbers of directories and files nested to the specified path")
	cmd.Flags().BoolVarP(&colorize, "color", "c", false, "for colorizing the folders")
	return cmd
}

func tree(path string, level int, symbol string) error {
	//var symbol string = elemChar
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat error for the path %s, error: %v", path, err)
	}

	// Do not print if it is the first level (error if the level 0 is not a directory)
	if level == 0 && !fi.IsDir() {
		return fmt.Errorf("Error: the given path %s is not a directory", path)
	}

	if !fi.IsDir() {
		if level != 0 {
			files++
			fmt.Printf("%s\n", fi.Name())
		}
		return nil
	} else {
		if level != 0 {
			dirs++
			if colorize {
				fmt.Printf("%s\n", utils.BlueS(fi.Name()))
			} else {
				fmt.Printf("%s\n", fi.Name())
			}
		}
	}

	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("ReadDir error for the path %s, error: %v", path, err)
	}

	for i, fi := range fis {
		if level >= levels {
			break
		}
		if len(fi.Name()) > 0 && strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		fmt.Printf(symbol)
		line := ""
		// is the last element?
		if i == len(fis)-1 {
			fmt.Printf("└── ")
			line = "   "
		} else {
			line = "│  "
			fmt.Printf("├── ")
		}
		if err := tree(filepath.Join(path, fi.Name()), level+1, symbol+line); err != nil {
			return err
		}
	}

	return nil
}
