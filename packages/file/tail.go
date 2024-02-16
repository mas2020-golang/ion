package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/ion/packages/utils"
)

type Tail struct {
}

func NewTail() *Tail {
	return &Tail{}
}

func (t *Tail) Tail(args []string, r int) (lines []string, err error) {
	var (
		f *os.File = utils.GetBytesFromPipe()
	)
	if f == nil {
		if len(args) == 0 {
			return nil, fmt.Errorf("no file argument")
		}
		// load the file into the buffer
		f, err = os.Open(args[0])
		return nil, err
	}
	return t.getLines(f, r)
}

// getLines returns the n lines in buf
func (t *Tail) getLines(f *os.File, n int) (lines []string, err error) {
	var (
		b      []byte
		l      [][]byte
		pos, i int
	)

	defer func() {
		err := f.Close()
		output.Error("Tail.getLines()", err.Error())
	}()

	buf, err := ioutil.ReadAll(f)
	//fmt.Printf("file content (bytes): %v\n", buf)
	output.CheckErrorAndExit("", "", err)
	pos = len(buf) - 1

	for i < n && pos >= 0 {
		// \n
		if buf[pos] == 10 {
			if len(b) == 0 {
				b = append(b, buf[pos])
				pos--
			} else {
				// create a new line
				l = append(l, t.reverseSlice(b))
				b = make([]byte, 0)
				i++
			}
		} else {
			b = append(b, buf[pos])
			pos--
		}
		// edge case: beginning of the file
		if pos == -1 && len(b) > 0 {
			l = append(l, t.reverseSlice(b))
		}
	}
	// create the slice of lines
	for i := len(l) - 1; i >= 0; i-- {
		lines = append(lines, string(l[i]))
	}
	return
}

func (t *Tail) reverseSlice(b []byte) []byte {
	var (
		r []byte
	)
	if len(b) == 1 {
		return b
	}
	for i := len(b) - 1; i >= 0; i-- {
		r = append(r, b[i])
	}
	return r
}
