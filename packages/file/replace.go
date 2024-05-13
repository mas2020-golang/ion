package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/utils"
)

type Replacer struct {
	all          bool
	verbose      bool
	pattern      string
	substitution string
}

func NewReplacer(verbose, all bool, pattern, substitution string) *Replacer {
	return &Replacer{
		all:          all,
		verbose:      verbose,
		pattern:      pattern,
		substitution: substitution,
	}
}

func (r *Replacer) Replace(path string) error {
	var oriLine, newLine, verboseLine string
	var replaced bool

	utils.Verbose(output.YellowS("=================================================\n"), r.verbose)
	utils.Verbose(output.YellowS("Yellow, green and red colors represent verbosity information.\nLegend:\n"), r.verbose)
	utils.Verbose(output.GreenS("IT IS SOMETHING ADDED\n"), r.verbose)
	utils.Verbose(output.RedS("IT IS SOMETHING REMOVED\n"), r.verbose)
	utils.Verbose(output.YellowS("=================================================\n"), r.verbose)

	// Open the input file for reading
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// // Open the output file for writing
	// outputFile, err := os.Create(path + ".repl")
	// if err != nil {
	// 	return err
	// }
	// defer outputFile.Close()

	// Create a reader for the input file
	reader := bufio.NewReader(inputFile)

	// Create a writer for the output file
	// writer := bufio.NewWriter(outputFile)

	// Read the file line by line
	if len(r.pattern) > 0 {
		for {
			oriLine, err = reader.ReadString('\n')
			// utils.Verbose(line, r.verbose)
			if err != nil && err != io.EOF {
				return err
			}

			// Replace the matched strings using the replace function
			replaced, newLine, verboseLine = r.replacePattern(oriLine)
			// Write the line to the output file
			// _, writeErr := writer.WriteString(line)
			// if writeErr != nil {
			// 	return writeErr
			// }
			fmt.Print(newLine)
			if replaced {
				utils.Verbose(verboseLine, r.verbose)
			}

			if err == io.EOF {
				break
			}
			if !r.all && replaced {
				break
			}
		}
	}

	// Flush the writer to ensure all data is written to the output file
	// err = writer.Flush()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (r *Replacer) replacePattern(text string) (bool, string, string) {
	// Compile the regular expression pattern.
	regexp := regexp.MustCompile(r.pattern)

	// Find all matches of the regular expression pattern in the text.
	matches := regexp.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return false, text, ""
	}
	pointer := 0
	result := strings.Builder{}
	verboseResult := strings.Builder{}

	utils.Verbose(output.YellowS(fmt.Sprintf("matches: %v\n", matches)), r.verbose)
	// in case of verbosity the substitution pattern is green
	if r.verbose {
		r.substitution = output.GreenS(r.substitution)
	}
	for _, pairs := range matches {
		start := pairs[0]
		end := pairs[1]
		result.WriteString(text[pointer:start])
		result.WriteString(r.substitution)
		if r.verbose {
			verboseResult.WriteString(text[pointer:start])
			verboseResult.WriteString(output.RedS(text[start:end]))
		}
		pointer = end
	}
	if pointer < len(text) {
		result.WriteString(text[pointer:])
		if r.verbose {
			verboseResult.WriteString(text[pointer:])
		}
	}
	return true, result.String(), verboseResult.String()
}
