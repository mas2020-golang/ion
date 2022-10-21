/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/mas2020-golang/goutils/output"
	out "github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

var (
	nocolors, countLines, countPattern, onlyMatch, invert bool
	insensitive, onlyResult                               bool
	cmd                                                   *cobra.Command
	matchLines, matchPattern, after, before, currentLine  int
	matchElems, lines                                     []string
	prLines                                               map[int]bool // save the printed lines
	prAfter                                               map[int]bool // lines to print after the match
)

func NewSearchCmd() *cobra.Command {
	cmd = &cobra.Command{
		Args: cobra.MinimumNArgs(1),
		Use:  "search <PATTERN> <PATH> [...PATH]",
		Example: `# search this in the demo-file
$ ion search "this" demo-file`,
		Short: "Search for the given pattern into the standard input or one or more files",
		Long: `The command searches for the pattern given as a first parameter. The command can search
directly from the standard input or one or more files passed an argument. The pattern is highlighted with red.`,
		Run: func(cmd *cobra.Command, args []string) {
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
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode active") //TODO
	return cmd
}

func search(args []string) {
	out.TraceLog("", "search starting...")
	prLines = make(map[int]bool)
	prAfter = make(map[int]bool)
	if len(args) == 0 {
		//out.Error("", "the pattern is missing")
		cmd.PrintErr("the pattern is missing")
		// cmd.Fprint(cmd.OutOrStderr(), "the pattern is missing")
		return
	}
	var (
		f *os.File = utils.GetBytesFromPipe()
	)
	if f == nil || os.Getenv("ION_DEBUG") == "true" {
		if len(args) <= 1 {
			cmd.PrintErrln("Error: no files as argument")
			os.Exit(1)
		}
		// load the file into the buffer
		for i := 1; i < len(args); i++ {
			f, err := os.Open(args[i])
			out.CheckErrorAndExit("", "opening the file as an argument", err)

			if !countLines && !countPattern && len(args) > 2 {
				if nocolors {
					cmd.Printf("=> on '%s':\n", args[i])
				} else {
					cmd.Printf("=> on '%s':\n", out.YellowBoldS(args[i]))
				}
			}
			//err = readLines(cmd, args[0], f)
			matchElems = []string{}
			matchLines, matchPattern = 0, 0
			err = checkLine(args[0], f)
			out.CheckErrorAndExit("", "", err)
		}
	} else {
		// read from the standard input
		err := checkLine(args[0], f)
		out.CheckErrorAndExit("", "", err)
	}
}

// checkLine checks any line to find the pattern matching
func checkLine(pattern string, f *os.File) error {
	var found int
	// remember to close the file at the end of the program
	defer f.Close()
	if insensitive {
		pattern = "(?i)" + pattern
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		currentLine++
		s := scanner.Text()
		// save the lines in case -B or -A is given
		if before > 0 || after > 0 {
			lines = append(lines, s)
		}
		results := r.FindAllStringIndex(s, -1)
		if results != nil {
			// at the first match, if onlyResult exits
			if onlyResult {
				found = 1
				break
			}

			// if invert do not print the match
			if invert {
				continue
			}
			out.TraceLog("", fmt.Sprintf("line num. %d, line => %s, results => %v\n", currentLine, s, results))
			// there is at least one match
			printResults(results, s)
			matchLines++
			matchPattern += len(results)
		} else {
			// show the entire line non in match with the patter
			if invert {
				cmd.Println(s)
			}
		}
	}
	// --invert is not compatible with onlyMatch, countLines, countPattern, onlyResult
	if invert {
		return nil
	}
	if onlyResult {
		cmd.Println(found)
		return nil
	}
	// print only the matches patterns
	if onlyMatch {
		for _, el := range matchElems {
			PrintColor(el + "\n")
		}
	}
	// print only the lines in match with the pattern
	if countLines {
		cmd.Println(matchLines)
	}
	// print only the number of times the pattern is matching
	if countPattern {
		cmd.Println(matchPattern)
	}
	return nil
}

func printResults(results [][]int, line string) {
	if countLines || countPattern {
		return
	}
	start := 0
	for _, el := range results {
		matchElems = append(matchElems, line[el[0]:el[1]])
		if onlyMatch {
			continue
		}
		printBefore()
		if el[0] > start {
			Print(line[start:el[0]])
			PrintColor(line[el[0]:el[1]])
		} else {
			PrintColor(line[el[0]:el[1]])
		}
		start = el[1]
	}
	if !onlyMatch {
		// print the rest of the line
		if start < len(line) {
			Print(line[start:])
			cmd.Println()
			prLines[currentLine] = true
		}
		if start == len(line) {
			cmd.Println()
		}
		setAfter()
		printAfter()
	}
}

func Print(text string) {
	cmd.Print(text)
}

func PrintColor(text string) {
	if !nocolors {
		cmd.Printf("%s", out.RedS(text))
	} else {
		Print(text)
	}
}

func printBefore() {
	start, end := 0, 0
	if currentLine == 1 {
		return
	}
	end = currentLine - 1
	if currentLine-before > 0 {
		start = currentLine - before
	}
	for i := start; i <= end; i++ {
		// print only if the was not already printed
		if ok, _ := prLines[i]; !ok {
			cmd.Println(lines[i])
			prLines[i] = true
			out.TraceLog("printBefore", fmt.Sprintf("save line %d", i))
		}
	}
}

// set what to print before the match
func setAfter() {
	if after > 0 {
		for i := 1; i <= after; i++ {
			prAfter[currentLine+i] = true
		}
	}
	output.TraceLog("", fmt.Sprintf("after lines: %v", prAfter))
}

// set what to print after the match
func printAfter() {
	for k, _ := range prLines {
		// print only if the was not already printed
		if ok, _ := prAfter[k]; !ok {
			cmd.Println(lines[k])
			prLines[k] = true
			out.TraceLog("printAfter", fmt.Sprintf("save line %d", k))
		}
	}
}
