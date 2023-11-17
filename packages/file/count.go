package file

import (
	"fmt"
	"io"
	"os"

	"github.com/mas2020-golang/ion/packages/utils"
)

type Counter struct {
	words bool
}

func NewCounter(words bool) *Counter {
	return &Counter{words: words}
}

func (c *Counter) Count(path string) (count int, err error) {
	var (
		f *os.File = utils.GetBytesFromPipe()
	)
	if f == nil { // no standard input, file name is expected
		if len(path) == 0 {
			return -1, fmt.Errorf("no file argument")
		}
		// load the file into the buffer
		f, err = os.Open(path)
		if err != nil {
			return -1, err
		}
	}
	return c.getCount(f)
}

// getLines returns the number of lines in f or the words depending to
// the path passed to the command
func (c *Counter) getCount(f *os.File) (count int, err error) {
	var (
		b     []byte
		space bool
	)

	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	space = true // first iteration starts with a space
	// read 1024 each time
	b = make([]byte, 1024)
	for {
		n, err := f.Read(b)
		if n > 0 {
			for i := 0; i < n; i++ {
				// case --word
				if c.words {
					if !isSpace(b[i]) && space {
						count++
					}
					// save if the current byte is a space to check for the
					// next iteration
					space = isSpace(b[i])
				} else {
					// case line count
					if b[i] == 10 {
						count++
					}
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			count = -1
			err = fmt.Errorf("read %d bytes: %v", n, err)
			break
		}
	}
	return count, err
}

// isSpace returns true if the byte passed is considered a char delimiter.
// The following chars are considered space: \n \t space.
func isSpace(b byte) bool {
	return b == 9 || b == 32 || b == 10
}
