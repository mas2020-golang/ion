package file

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mas2020-golang/goutils/output"
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

	// sliceBytes takes the precedence over sliceChars and sliceCols
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
	start, end, err := t.getIntervals(sliceInterval)
	if err != nil {
		return nil, err
	}

	// read the input line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		// determine the operation to execute
		switch op {
		case "bytes":
			sliceVal = append(sliceVal, t.getSliceBytes(s, start, end))
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

	return sliceVal, nil
}

// getIntervals accepts the arg to check for returning start and end
func (t *Slice) getIntervals(arg string) (int, int, error) {
	// any :?
	if strings.Contains(arg, ":") {
		// edge case, only : is given
		if arg == ":" {
			return -1, -1, utils.ErrMalformed
		}
		elems := strings.Split(arg, ":")
		output.Trace("file.getIntervals()", fmt.Sprintf("elems: %#v", elems))
		// 2 elems?
		if len(elems) != 2 {
			return -1, -1, utils.ErrMalformed
		}
		start, err := strconv.Atoi(strings.Trim(elems[0], " "))
		if err != nil {
			return start, -1, err
		}
		end := 0
		// check if end is empty value (end is empty in case of giving 1:)
		if len(strings.Trim(elems[1], " ")) == 0 {
			end = math.MaxInt32 // due to this in the getSliceBytes func will be set end = len(s)
		} else {
			end, err = strconv.Atoi(strings.Trim(elems[1], " "))
		}

		if err != nil {
			return start, end, err
		}

		// no errors
		return start, end, nil
	} else {
		start, err := strconv.Atoi(strings.Trim(arg, " "))
		if err != nil {
			return -1, -1, utils.ErrMalformed
		}
		return start, start, err
	}
}

func (t *Slice) getSliceBytes(s string, start, end int) string {
	start--
	if start < 0 || start > len(s) || end < 0 || start > end {
		return ""
	}

	if end > len(s) {
		end = len(s)
	}

	return s[start:end]
}

// Returns the slices (it stops at the first invalid index)
func (t *Slice) getSlices(s string, d string, fields []uint8) []string {
	response := make([]string, 0)
	slices := strings.Split(s, d)

	for _, f := range fields {
		f--
		if f < 0 || int(f) >= len(slices) {
			return response
		}
		if int(f) < len(slices) {
			response = append(response, slices[f])
		}
	}
	return response
}
