package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	Version, GitCommit string
	templateFuncs      = template.FuncMap{
		"trim":                    strings.TrimSpace,
		"trimRightSpace":          trimRightSpace,
		"trimTrailingWhitespaces": trimRightSpace,
	}
)

func init() {
	Version = "0.2.0-dev"
}

// GetBytesFromPipe reads from the pipe and return the buffer of bytes of the given argument
func GetBytesFromPipe() *os.File {
	//var bs []byte
	//buf := bytes.NewBuffer(bs)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		//scanner := bufio.NewScanner(os.Stdin)
		//
		//for scanner.Scan() {
		//	buf.Write(scanner.Bytes())
		//	fmt.Print(scanner.Text())
		//}
		//
		//if err := scanner.Err(); err != nil {
		//	log.Fatal(err)
		//}
		return os.Stdin
	}
	//fmt.Printf("number of bytes from the pipe are %d\n", len(buf.Bytes()))
	return nil
}

// ReadPassword reads the standard input in hidden mode
func ReadPassword(text string) (string, error) {
	fmt.Print(text)
	buf, err := term.ReadPassword(0)
	return string(buf), err
}

// Check checks if an error and exit
func Check(err error) {
	if err != nil {
		fmt.Printf("%s %v\n", RedS("Error:"), err)
		os.Exit(1)
	}
}

// TODO: templating functions to move into another file.go or package

// GetHelpFunction prints the output on the screen based on a specific template
func GetHelpFunction(t string) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		printCommandHelp(os.Stdout, t, c)
	}
}

// getCommandHelp return the help for the single command
func printCommandHelp(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	t.Execute(w, data)
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}
