package file

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mas2020-golang/ion/packages/utils"
)

type Slice struct {
}

func NewSlice() *Slice {
	return &Slice{}
}

func (t *Slice) Slice(arg string, sliceBytes string, sliceChars string, sliceCols string) (sliceVal string, err error) {
	var f *os.File = utils.GetBytesFromPipe()
	if f == nil { // means that there is no pipe value
		if len(arg) == 0 {
			return "", fmt.Errorf("no file argument")
		}
		// load the file into the buffer
		f, err = os.Open(arg)
		if err != nil {
			return "", err
		}
	}
	// TODO: 1. read the file line by line; 2. invoke the right func; 3. return the value or the error
	sliceInterval := sliceBytes
	if len(sliceInterval) == 0 {
		sliceInterval = sliceChars
	}
	if len(sliceInterval) == 0 {
		sliceInterval = sliceCols
	}
	start, end, err := getIntervals(sliceInterval)
	if err != nil {
		return "", fmt.Errorf("error getting intervals: %s", err)
	}

	// sliceBytes takes the precedence on sliceChars and sliceCols
	return fmt.Sprintf("start: %d, end: %d", start, end), nil
}

func getIntervals(arg string) (int, int, error) {
	// any :?
	if strings.Contains(arg, ":") {
		elems := strings.Split(arg, ":")
		// 2 elems?
		if len(elems) != 2 {
			return -1, -1, errors.New("passing the colon you have to specify an interval, e.g. 3:5")
		}
		start, err := strconv.Atoi(strings.Trim(elems[0], " "))
		if err != nil {
			return start, -1, err
		}
		end, err := strconv.Atoi(strings.Trim(elems[1], " "))
		if err != nil {
			return start, end, err
		}

		// no errors
		return start, end, nil
	} else {
		start, err := strconv.Atoi(arg)
		return start, -1, err
	}
}

func getStringByBytes(s string, start, end int) string {
	start--
	if start < 0 || start > len(s) || end < 0 || start > end {
		return ""
	}

	if end > len(s) {
		end = len(s)
	}

	return s[start:end]
}
