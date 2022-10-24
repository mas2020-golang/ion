/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	out "github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

var (
	nocolors, countLines, countPattern, onlyMatch, invert bool
	insensitive, onlyResult, onlyFilename, isDir          bool
	recursive                                             bool
	cmd                                                   *cobra.Command
	before, after, argsN, currentLine, level              int
	prLines                                               map[int]bool // save the printed lines
	mLines                                                map[int]string
	mLinesMatch                                           map[int][][]int
)

func NewSearchCmd() *cobra.Command {
	cmd = &cobra.Command{
		Args: cobra.MinimumNArgs(1),
		Use:  "search <PATTERN> <PATH> [...PATH]",
		Example: `# search this in the demo-file
$ ion search "this" demo-file`,
		Short: "Search for the given pattern into the standard input or one or more files",
		Long: `The command searches for the given pattern. The command can search
directly from the standard input, one or more files or directories passed an argument. The pattern is highlighted with the red color.`,
		Run: func(cmd *cobra.Command, args []string) {
			argsN = len(args)
			search(args)
		},
	}
	cmd.SetOut(os.Stdout)
	// flags
	cmd.Flags().BoolVarP(&countLines, "count-lines", "l", false, "shows only how many lines match with the pattern")
	cmd.Flags().BoolVarP(&countPattern, "count-pattern", "p", false, "shows only how many time a pattern is in match")
	cmd.Flags().BoolVarP(&onlyMatch, "only-match", "m", false, "shows only the substring that match, not the entire line")
	cmd.Flags().BoolVarP(&nocolors, "no-colors", "n", false, "no colors are printed onto the standard output")
	cmd.Flags().BoolVarP(&invert, "invert", "t", false, "shows the lines that doesn't match with the pattern")
	cmd.Flags().BoolVarP(&insensitive, "insensitive", "i", false, "to match with no case sensitivity")
	cmd.Flags().BoolVarP(&onlyResult, "only-result", "r", false, "if there is at least one match it returns 1, otherwise 0")
	cmd.Flags().IntVarP(&before, "before", "B", 0, "shows also the NUMBER of lines before the match")
	cmd.Flags().IntVarP(&after, "after", "A", 0, "shows also the NUMBER of lines after the match")
	cmd.Flags().BoolVarP(&onlyFilename, "only-filename", "f", false, "shows only the filename when a pattern matches one or several times") //TODO: to implement
	cmd.Flags().BoolVarP(&recursive, "recursive", "d", false, "if the PATH is a folder searches in the sub folders too and not only in its first level")
	return cmd
}

func search(args []string) {
	out.TraceLog("", "search starting...")
	prLines = make(map[int]bool)
	if len(args) == 0 {
		//out.Error("", "the pattern is missing")
		cmd.PrintErr("the pattern is missing")
		// cmd.Fprint(cmd.OutOrStderr(), "the pattern is missing")
		return
	}
	var (
		// read from the standard input
		f *os.File = utils.GetBytesFromPipe()
	)
	if f == nil || os.Getenv("ION_DEBUG") == "true" {
		if len(args) <= 1 {
			cmd.PrintErrln("Error: no files as argument")
			os.Exit(1)
		}
		// open the files
		for i := 1; i < len(args); i++ {
			level = 0
			// search
			err := searchInFiles(args[0], args[i])
			out.CheckErrorAndExit("", "", err)
		}
	} else {
		// search
		err := startSearching(args[0], f, "")
		out.CheckErrorAndExit("", "", err)
	}
}

// searchInFiles opens each file and start for searching
func searchInFiles(pattern string, path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	// case path is a dir
	if fi.IsDir() {
		level++
		// if not recursive stop at level 1
		if level > 1 && !recursive {
			return nil
		}
		isDir = true
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("ReadDir error for the path %s, error: %v", path, err)
		}

		for _, fi := range fis {
			err = searchInFiles(pattern, filepath.Join(path, fi.Name()))
			if err != nil {
				return err
			}
		}
	}

	// case path is a file
	f, err := os.Open(path)
	defer f.Close()
	out.CheckErrorAndExit("", "opening the file as an argument", err)
	// search
	err = startSearching(pattern, f, path)
	if err != nil {
		return err
	}
	return nil
}

// Start the searching
func startSearching(pattern string, f *os.File, filename string) error {
	err := readLines(pattern, f)
	if err != nil {
		return err
	}
	// print the name of the files only if there is a match and the file passed are more that one
	if len(mLinesMatch) > 0 {
		if argsN > 2 || isDir {
			if nocolors {
				cmd.Printf("> '%s':\n", filename)
			} else {
				cmd.Printf("> '%s':\n", out.YellowBoldS(filename))
			}
		}
	}

	// check the flags and print the result
	checkFlags()
	out.TraceLog("search", fmt.Sprintf("matched lines: %v", mLinesMatch))
	return nil
}

// readLines reads the file, save the content and save the lines in match
// with the pattern
func readLines(pattern string, f *os.File) error {
	mLines = make(map[int]string)
	mLinesMatch = make(map[int][][]int)
	currentLine = 0

	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	// check --insentive flag
	if insensitive {
		pattern = "(?i)" + pattern
	}

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		s := scanner.Text()
		// add the line
		mLines[currentLine] = s
		results := r.FindAllStringIndex(s, -1)
		if results != nil {
			out.TraceLog("", fmt.Sprintf("line num. %d, line => %s, results => %v\n", currentLine, s, results))
			mLinesMatch[currentLine] = results
			// the lines in match are set as to be already printed (needed for the printAfter and the printBefore)
			//prLines[currentLine] = true
		}
		currentLine++
	}
	return nil
}

// checkFlags elaborates the result based on some flags that have special behaviours:
// --count-lines, --count-pattern, --only-match
// Returns true in case no more output is needed
func checkFlags() {
	// check --only-result
	if onlyResult {
		if len(mLinesMatch) > 0 {
			cmd.Println(1)
			return
		}
		cmd.Println(0)
		return
	}
	// print only the lines
	if countLines {
		cmd.Println(len(mLinesMatch))
		return
	}
	// print only the number of elements in match with the pattern
	if countPattern {
		i := 0
		for k, _ := range mLinesMatch {
			i += len(mLinesMatch[k])
		}
		cmd.Println(i)
		return
	}

	// print only the match elements
	if onlyMatch {
		printMatch(true)
		return
	}

	// check --invert
	if invert {
		printInvert()
		return
	}

	// print the lines in match
	printLinesInMatch()
}

// printLinesInMatch prints the single line in match
func printLinesInMatch() {
	// cycle the match lines
	for i := 0; i < currentLine; i++ {
		// continue if no match or it is not already printed
		if _, ok := mLinesMatch[i]; !ok {
			continue
		}
		l := mLines[i] // line to print
		printBefore(i)
		printMatchLine(l, mLinesMatch[i], i)
		// print After
		printAfter(i)
	}
}

// printMatchLine prints the line in match
func printMatchLine(l string, elems [][]int, indexL int) {
	p := 0
	// skip it is already printed
	if _, ok := prLines[indexL]; ok {
		return
	}
	for _, m := range elems {
		if m[0] > p {
			cmd.Print(l[p:m[0]])
		}
		// print the match
		if nocolors {
			cmd.Print(l[m[0]:m[1]])
		} else {
			cmd.Print(out.RedBoldS(l[m[0]:m[1]]))
		}
		p = m[1]
	}
	// print the end of the sprint (eventually)
	if p <= len(l) {
		cmd.Println(l[p:])
	}
	prLines[indexL] = true
}

// Print the matching elements on the standard output
func printMatch(newLine bool) {
	for k, v := range mLinesMatch {
		s := mLines[k]
		for _, m := range v {
			if newLine {
				cmd.Println(s[m[0]:m[1]])
			} else {
				cmd.Print(s[m[0]:m[1]])
			}
		}
	}
}

// printInvert prints the line that does not contain any
// match with the pattern
func printInvert() {
	// cycle line by line
	for i := 0; i < currentLine; i++ {
		if _, ok := mLinesMatch[i]; !ok {
			cmd.Println(mLines[i])
		}
	}
}

// printAfter prints the lines after the matches. It happens after
// the print of each match pattern. It prints the line only if it has been not
// already printed. currentPrintL contains the value of the current line that
// has been printed on the screen.
func printAfter(currentPrintL int) {
	out.TraceLog("printAfter", fmt.Sprintf("start currentPrintL: %d, currentLine: %d, after: %d", currentPrintL, currentLine, after))
	p := after
	// currentPrintL is the end of the file
	if currentPrintL == currentLine || p == 0 {
		return
	}
	// move forward to the next line to print
	currentPrintL++
	for currentPrintL <= currentLine && p > 0 {
		out.TraceLog("printAfter", fmt.Sprintf("currentPrintL: %d, currentLine: %d, p: %d", currentPrintL, currentLine, p))
		// if the lines has been already printed move on
		if _, ok := prLines[currentPrintL]; ok {
			currentPrintL++
			p--
			continue
		}
		// if the line match the pattern print highlighting the matches else it print normally
		if _, ok := mLinesMatch[currentPrintL]; ok {
			printMatchLine(mLines[currentPrintL], mLinesMatch[currentPrintL], currentPrintL)
		} else {
			cmd.Println(mLines[currentPrintL])
			prLines[currentPrintL] = true
		}
		// move to the next line and decrease the pointer
		currentPrintL++
		p--
	}
}

// printBefore prints the n lines before the line in match with the pattern
func printBefore(currentPrintL int) {
	if currentPrintL == 0 || before == 0 {
		return
	}
	p := currentPrintL - before
	if currentPrintL-before < 0 {
		p = 0
	}

	for currentPrintL > p {
		out.TraceLog("printBefore", fmt.Sprintf("currentPrintL: %d, currentLine: %d, p: %d", currentPrintL, currentLine, p))
		// if the lines has been already printed move on
		if _, ok := prLines[p]; ok {
			p++
			continue
		}
		// print the line
		cmd.Println(mLines[p])
		prLines[p] = true
		p++
	}
}
