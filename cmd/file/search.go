/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	out "github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

var (
	nocolors, countLines, countPattern, onlyMatch, invert bool
	cmd                                                   *cobra.Command
	matchLines, matchPattern                              int
	matchElems                                            []string
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

	// flags
	cmd.Flags().BoolVarP(&countLines, "count-lines", "l", false, "shows only how many lines match with the pattern")
	cmd.Flags().BoolVarP(&countPattern, "count-pattern", "p", false, "shows only how many time a pattern is in match")
	cmd.Flags().BoolVarP(&onlyMatch, "only-match", "o", false, "shows only the substring that match, not the entire line")
	cmd.Flags().BoolVarP(&nocolors, "no-colors", "n", false, "no colors on the standard output")
	cmd.Flags().BoolVarP(&invert, "invert", "i", false, "shows the lines that doesn't match with the pattern") //TODO
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode active")                               //TODO
	return cmd
}

func search(args []string) {
	out.TraceLog("", "search starting...")
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

			if !countLines && !countPattern {
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
	// remember to close the file at the end of the program
	defer f.Close()
	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		s := scanner.Text()
		results := r.FindAllStringIndex(s, -1)
		if results != nil {
			out.TraceLog("", fmt.Sprintf("line: %s, results: %v", s, results))
			// there is at least one match
			printResults(results, s)
			matchLines++
			matchPattern += len(results)
		}
	}
	// print only the matches patters
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
	if countPattern && !invert {
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
		if el[0] > start {
			Print(line[start:el[0]])
			PrintColor(line[el[0]:el[1]])
		} else {
			PrintColor(line[el[0]:el[1]])
		}
		start = el[1]
	}
	if !onlyMatch {
		if start < len(line) {
			Print(line[start:])
			cmd.Println()
		}
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
