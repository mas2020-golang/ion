package file

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mas2020-golang/ion/packages/utils"
)

type Slice struct {
}

func NewSlice() *Slice {
	return &Slice{}
}

func (t *Slice) Slice(arg string, sliceBytes string, sliceChars string, sliceCols string) (sliceVal []string, err error) {
	// reads the arg input file
	f, err := utils.GetReader(arg)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// TODO: 1. read the file line by line; 2. invoke the right func; 3. return the value or the error
	sliceInterval := sliceBytes
	op := "bytes"
	if len(sliceInterval) == 0 {
		sliceInterval = sliceChars
		op = "chars"
	}
	if len(sliceInterval) == 0 {
		sliceInterval = sliceCols
		op = "cols"
	}
	start, end, err := getIntervals(sliceInterval)
	if err != nil {
		return nil, fmt.Errorf("error getting intervals: %s", err)
	}

	// read the input line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		// determine the operation to execute
		switch op {
		case "bytes":
			sliceVal = append(sliceVal, getSliceBytes(s, start, end))
		case "chars":
			fmt.Println("not implemented yet")
		case "cols":
			fmt.Println("not implemented yet")
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("error during the scan of the file: %s", err)
	}

	// sliceBytes takes the precedence on sliceChars and sliceCols
	return sliceVal, nil
}

// getIntervals accepts the arg to check for returning start and end
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
		start, err := strconv.Atoi(strings.Trim(arg, " "))
		return start, start, err
	}
}

func getSliceBytes(s string, start, end int) string {
	start--
	if start < 0 || start > len(s) || end < 0 || start > end {
		return ""
	}

	if end > len(s) {
		end = len(s)
	}

	return s[start:end]
}
