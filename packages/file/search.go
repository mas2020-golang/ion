/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	out "github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/utils"
)

type Searcher struct {
	nocolors, countLines, countPattern, onlyMatch, invert bool
	insensitive, onlyResult, onlyFilename, isDir          bool
	recursive                                             bool
	before, after, argsN, currentLine, level              int
	prLines                                               map[int]bool // save the printed lines
	mLines                                                map[int]string
	mLinesMatch                                           map[int][][]int
	output                                                io.Writer
}

func NewSearcher(countLines, countPattern, onlyMatch, nocolors, invert, insensitive, onlyResult,
	onlyFilename, recursive bool, before, after int, output io.Writer) *Searcher {
	return &Searcher{
		countLines:   countLines,
		countPattern: countPattern,
		onlyMatch:    onlyMatch,
		nocolors:     nocolors,
		invert:       invert,
		insensitive:  insensitive,
		onlyResult:   onlyResult,
		onlyFilename: onlyFilename,
		recursive:    recursive,
		before:       before,
		after:        after,
		output:       output,
	}
}

func (s *Searcher) Search(args []string) error {
	out.TraceLog("", "search starting...")
	(*s).prLines = make(map[int]bool)
	if len(args) == 0 {
		//out.Error("", "the pattern is missing")
		return fmt.Errorf("the pattern is missing")
	}
	var (
		// read from the standard input
		f *os.File = utils.GetBytesFromPipe()
	)
	if f == nil || os.Getenv("ION_DEBUG") == "true" {
		if len(args) <= 1 {
			out.CheckErrorAndExit("", "", fmt.Errorf("no files as argument"))
		}
		// open the files
		for i := 1; i < len(args); i++ {
			s.level = 0
			// search
			err := s.searchInFiles(args[0], args[i])
			out.CheckErrorAndExit("", "", err)
		}
	} else {
		// search
		err := s.startSearching(args[0], f, "")
		out.CheckErrorAndExit("", "", err)
	}
	return nil
}

// searchInFiles opens each file and start for searching
func (s *Searcher) searchInFiles(pattern string, path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	// case path is a dir
	if fi.IsDir() {
		s.level++
		// if not recursive stop at level 1
		if s.level > 1 && !s.recursive {
			return nil
		}
		s.isDir = true
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("ReadDir error for the path %s, error: %v", path, err)
		}

		for _, fi := range fis {
			err = s.searchInFiles(pattern, filepath.Join(path, fi.Name()))
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
	err = s.startSearching(pattern, f, path)
	if err != nil {
		return err
	}
	return nil
}

// Start the searching
func (s *Searcher) startSearching(pattern string, f *os.File, filename string) error {
	err := s.readLines(pattern, f)
	if err != nil {
		return err
	}
	// print the name of the files only if there is a match and the file passed are more than one
	// if len(mLinesMatch) > 0 {
	if (s.argsN > 2 || s.isDir) && !s.onlyFilename {
		if s.nocolors {
			s.output.Write([]byte(fmt.Sprintf("> '%s':\n", filename)))
		} else {
			s.output.Write([]byte(fmt.Sprintf("> '%s':\n", out.YellowBoldS(filename))))
		}
	}

	// check the flags and print the result
	s.checkFlags(filename)
	out.TraceLog("search", fmt.Sprintf("matched lines: %v", s.mLinesMatch))
	return nil
}

// readLines reads the file, save the content and save the lines in match
// with the pattern
func (s *Searcher) readLines(pattern string, f *os.File) error {
	s.mLines = make(map[int]string)
	s.mLinesMatch = make(map[int][][]int)
	s.currentLine = 0

	// check --insentive flag
	if s.insensitive {
		pattern = "(?i)" + pattern
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		scan := scanner.Text()
		// add the line
		s.mLines[s.currentLine] = scan
		results := r.FindAllStringIndex(scan, -1)
		if results != nil {
			out.TraceLog("", fmt.Sprintf("line num. %d, line => %s, results => %v\n", s.currentLine, s, results))
			s.mLinesMatch[s.currentLine] = results
		}
		s.currentLine++
	}
	return nil
}

// checkFlags elaborates the result based on some flags that have special behaviours:
// --count-lines, --count-pattern, --only-match
// Returns true in case no more output is needed
func (s *Searcher) checkFlags(filename string) {
	// check --only-filename
	if s.onlyFilename {
		if len(s.mLinesMatch) > 0 {
			s.output.Write([]byte(fmt.Sprintf("%s\n", filename)))
		}
		return
	}
	// check --only-result
	if s.onlyResult {
		if len(s.mLinesMatch) > 0 {
			s.output.Write([]byte(fmt.Sprintf("%s\n", "1")))
			return
		}
		s.output.Write([]byte(fmt.Sprintf("%s\n", "0")))
		return
	}
	// print only the lines
	if s.countLines {
		s.output.Write([]byte(fmt.Sprintf("%d\n", len(s.mLinesMatch))))
		return
	}
	// print only the number of elements in match with the pattern
	if s.countPattern {
		i := 0
		for k, _ := range s.mLinesMatch {
			i += len(s.mLinesMatch[k])
		}
		s.output.Write([]byte(fmt.Sprintf("%d\n", i)))
		return
	}

	// print only the match elements
	if s.onlyMatch {
		s.printMatch(true)
		return
	}

	// check --invert
	if s.invert {
		s.printInvert()
		return
	}

	// print the lines in match
	s.printLinesInMatch()
}

// printLinesInMatch prints the single line in match
func (s *Searcher) printLinesInMatch() {
	// cycle the match lines
	for i := 0; i < s.currentLine; i++ {
		// continue if no match or it is not already printed
		if _, ok := s.mLinesMatch[i]; !ok {
			continue
		}
		l := s.mLines[i] // line to print
		s.printBefore(i)
		s.printMatchLine(l, s.mLinesMatch[i], i)
		// print After
		s.printAfter(i)
	}
}

// printMatchLine prints the line in match
func (s *Searcher) printMatchLine(l string, elems [][]int, indexL int) {
	p := 0
	// skip it is already printed
	if _, ok := s.prLines[indexL]; ok {
		return
	}
	for _, m := range elems {
		if m[0] > p {
			s.output.Write([]byte(fmt.Sprintf("%s", l[p:m[0]])))
		}
		// print the match
		if s.nocolors {
			s.output.Write([]byte(fmt.Sprintf("%s", l[m[0]:m[1]])))
		} else {
			s.output.Write([]byte(fmt.Sprintf("%s", out.RedBoldS(l[m[0]:m[1]]))))
		}
		p = m[1]
	}
	// print the end of the sprint (eventually)
	if p <= len(l) {
		s.output.Write([]byte(fmt.Sprintf("%s\n", l[p:])))
	}
	s.prLines[indexL] = true
}

// Print the matching elements on the standard output
func (s *Searcher) printMatch(newLine bool) {
	for k, v := range s.mLinesMatch {
		line := s.mLines[k]
		for _, m := range v {
			if newLine {
				s.output.Write([]byte(fmt.Sprintf("%s\n", line[m[0]:m[1]])))
			} else {
				s.output.Write([]byte(fmt.Sprintf("%s", line[m[0]:m[1]])))
			}
		}
	}
}

// printInvert prints the line that does not contain any
// match with the pattern
func (s *Searcher) printInvert() {
	// cycle line by line
	for i := 0; i < s.currentLine; i++ {
		if _, ok := s.mLinesMatch[i]; !ok {
			s.output.Write([]byte(fmt.Sprintf("%s\n", s.mLines[i])))
		}
	}
}

// printAfter prints the lines after the matches. It happens after
// the print of each match pattern. It prints the line only if it has been not
// already printed. currentPrintL contains the value of the current line that
// has been printed on the screen.
func (s *Searcher) printAfter(currentPrintL int) {
	out.TraceLog("printAfter", fmt.Sprintf("start currentPrintL: %d, currentLine: %d, after: %d", currentPrintL, s.currentLine, s.after))
	p := s.after
	// currentPrintL is the end of the file
	if currentPrintL == s.currentLine || p == 0 {
		return
	}
	// move forward to the next line to print
	currentPrintL++
	for currentPrintL <= s.currentLine && p > 0 {
		out.TraceLog("printAfter", fmt.Sprintf("currentPrintL: %d, currentLine: %d, p: %d", currentPrintL, s.currentLine, p))
		// if the lines has been already printed move on
		if _, ok := s.prLines[currentPrintL]; ok {
			currentPrintL++
			p--
			continue
		}
		// if the line match the pattern print highlighting the matches else it print normally
		if _, ok := s.mLinesMatch[currentPrintL]; ok {
			s.printMatchLine(s.mLines[currentPrintL], s.mLinesMatch[currentPrintL], currentPrintL)
		} else {
			s.output.Write([]byte(fmt.Sprintf("%s\n", s.mLines[currentPrintL])))
			s.prLines[currentPrintL] = true
		}
		// move to the next line and decrease the pointer
		currentPrintL++
		p--
	}
}

// printBefore prints the n lines before the line in match with the pattern
func (s *Searcher) printBefore(currentPrintL int) {
	if currentPrintL == 0 || s.before == 0 {
		return
	}
	p := currentPrintL - s.before
	if currentPrintL-s.before < 0 {
		p = 0
	}

	for currentPrintL > p {
		out.TraceLog("printBefore", fmt.Sprintf("currentPrintL: %d, currentLine: %d, p: %d", currentPrintL, s.currentLine, p))
		// if the lines has been already printed move on
		if _, ok := s.prLines[p]; ok {
			p++
			continue
		}
		// print the line
		s.output.Write([]byte(fmt.Sprintf("%s\n", s.mLines[p])))
		s.prLines[p] = true
		p++
	}
}
