package utils

import (
	"fmt"
	"os"

	"github.com/mas2020-golang/goutils/output"
	"golang.org/x/term"
)

var (
	Version, GitCommit, BuildDate, GoVersion string
)

func init() {
	Version = "0.4.0-dev"
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

// getReader loads the *os.File from the pipe or from arg file
func GetReader(args []string) (f *os.File, err error) {
	if len(os.Getenv("ION_DEBUG")) == 0 {
		f = GetBytesFromPipe()
	}
	if f == nil { // means that there is no pipe value
		if len(args) == 0 {
			return nil, fmt.Errorf("no file argument (remember that you can also redirect the pipe or the standard input)")
		}
		// load the file into the buffer
		f, err = os.Open(args[0])

		if err != nil {
			return nil, err
		}
		output.TraceLog("utils.GetReader", fmt.Sprintf("file %s, opened", args[0]))
	}
	return
}
