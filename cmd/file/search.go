/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	out "github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/utils"
	"github.com/spf13/cobra"
)

var (
	nocolors, countLines, countPattern, onlyMatch, invert bool
	cmd                                                   *cobra.Command
)

// TODO: put the cmd variable as a private variable for the package so that we can use fmt.Print, Printf etc...
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
	cmd.Flags().BoolVarP(&invert, "invert", "i", false, "shows the lines that doesn't match with the pattern")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode active")
	return cmd
}

func search(args []string) {
	if len(args) == 0 {
		//out.Error("", "the pattern is missing")
		cmd.PrintErr("the pattern is missing")
		// fmt.Fprint(cmd.OutOrStderr(), "the pattern is missing")
		return
	}
	var (
		f *os.File = utils.GetBytesFromPipe()
	)
	if f == nil || os.Getenv("ION_DEBUG") == "true" {
		if len(args) <= 1 {
			cmd.PrintErr("Error: no files as argument\n")
			//os.Exit(1)
			return
		}
		// load the file into the buffer
		for i := 1; i < len(args); i++ {
			f, err := os.Open(args[i])
			out.CheckErrorAndExit("", "opening the file as an argument", err)
			if i > 1 {
				fmt.Println()
			}
			if !countLines && !countPattern {
				if nocolors {
					fmt.Print(fmt.Sprintf("on '%s':\n", args[i]))
				} else {
					fmt.Printf(fmt.Sprintf("on '%s':\n", out.YellowBoldS(args[i])))
				}
			}
			//err = readLines(cmd, args[0], f)
			err = checkLine(args[0], f)
			out.CheckErrorAndExit("", fmt.Sprintf("reading the file %s", args[i]), err)
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
			// there is at least one match
			printResults(results, s)
			fmt.Println()
		}
	}
	return nil
}

func printResults(results [][]int, line string) {
	start := 0
	for _, el := range results {
		if el[0] > start {
			Print(line[start:el[0]])
		} else {
			PrintColor(line[el[0]:el[1]])
		}
		start += (el[1] - el[0])
	}
	if start < len(line) {
		Print(line[start:])
	}
}

func Print(text string) {
	fmt.Print(text)
}

func PrintColor(text string) {
	if !nocolors {
		fmt.Printf("%s", out.RedS(text))
	} else {
		Print(text)
	}
}

// readLines scans the file line by line
func readLines(cmd *cobra.Command, pattern string, f *os.File) error {
	var l, n, mLines int
	// remember to close the file at the end of the program
	defer f.Close()
	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		match := r.FindString(scanner.Text())
		if len(match) != 0 {
			if !invert {
				mLines++
				if !countLines {
					n += searchLineInMatch(scanner.Text(), match)
				}
			}
		} else {
			mLines++
			if invert && !countLines {
				// lines that do not match with the patterns
				fmt.Println(scanner.Text())
			}
		}
		l++
	}

	if countLines {
		fmt.Println(mLines)
	}
	if countPattern && !invert {
		fmt.Println(n)
	}
	// verbose mode
	if verbose {
		out.InfoBox(fmt.Sprintf("lines read %d", l))
		out.InfoBox(fmt.Sprintf("matching patterns: %d", n))
		out.InfoBox(fmt.Sprintf("matching lines: %d", mLines))
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func searchLineInMatch(line string, match string) int {
	p, n := 0, 0
	var output string
	for {
		idx := strings.Index(line, match)
		//fmt.Printf("\n-- search in %q the match %q, idx %d\n", line, match, idx)
		if idx >= 0 {
			n++
			printLine(line[p:idx], false)
			if nocolors {
				output = line[idx : idx+len(match)]
			} else {
				output = out.RedS(line[idx : idx+len(match)])
			}
			if !countLines && !countPattern {
				fmt.Print(output)
				if onlyMatch {
					fmt.Println()
				}
			}
			line = line[idx+len(match):]
			//fmt.Printf("\nnew line is: %s", out.YellowS(line))
		} else {
			printLine(line, true)
			//fmt.Println("3")
			break
		}
	}
	return n
}

// printLine prints the content into the standard output
func printLine(text string, newLine bool) {
	if countLines || countPattern || onlyMatch {
		return
	}
	if newLine {
		fmt.Println(text)
	} else {
		fmt.Print(text)
	}
}
