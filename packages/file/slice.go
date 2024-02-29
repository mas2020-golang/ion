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

func (t *Slice) Slice(args []string, sliceBytes string, sliceChars string, sliceCols string, delim string) (sliceVal []string, err error) {
	// reads the arg input file
	f, err := utils.GetReader(args)
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
		op = "fields"
	}
	output.TraceLog("file.Slice", fmt.Sprintf("op = %s, sliceInterval = %s", op, sliceInterval))
	startSl, end, err := t.getIntervals(sliceInterval)
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
			sliceVal = append(sliceVal, t.getBytes(s, startSl, end))
		case "chars":
			sliceVal = append(sliceVal, t.getChars(s, startSl, end))
		case "fields":
			val, err := t.getFields(s, delim, startSl, end)
			if err != nil {
				return nil, err
			}
			sliceVal = append(sliceVal, val)
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("error during the scan of the file: %s", err)
	}

	return sliceVal, nil
}

// getIntervals accepts the arg to check for returning start and end
func (t *Slice) getIntervals(arg string) ([]int, int, error) {
	// trimming first
	arg = strings.Trim(arg, " ")
	// any :?
	if strings.Contains(arg, ":") {
		// edge case, only : is given
		if arg == ":" {
			return nil, -1, utils.ErrMalformed
		}
		elems := strings.Split(arg, ":")
		output.TraceLog("file.getIntervals", fmt.Sprintf("elems: %#v", elems))
		// 2 elems?
		if len(elems) != 2 {
			return nil, -1, utils.ErrMalformed
		}
		start, err := strconv.Atoi(elems[0])
		if err != nil {
			return []int{start}, -1, utils.ErrMalformed
		}
		end := 0
		// check if end is empty value (end is empty in case of giving 1:)
		if len(strings.Trim(elems[1], " ")) == 0 {
			end = math.MaxInt32 // due to this in the getSliceBytes func will be set end = len(s)
		} else {
			end, err = strconv.Atoi(elems[1])
			if err != nil {
				return []int{start}, end, utils.ErrMalformed
			}
		}
		// start and end are done, is end < start?
		if end < start || start < 0 || end < 0 {
			return []int{start}, end, utils.ErrMalformed
		}
		// no errors
		return []int{start}, end, nil
	} else if strings.Contains(arg, ",") {
		// splitting by comma, we have multiple starts with the same end each one
		starts := []int{}
		elems := strings.Split(arg, ",")
		output.TraceLog("file.getIntervals", fmt.Sprintf("elems: %#v", elems))
		// convert the elems to int
		for _, s := range elems {
			intVal, err := strconv.Atoi(s)
			if err != nil {
				return nil, -1, utils.ErrMalformed
			}
			starts = append(starts, intVal)
		}
		return starts, -1, nil
	} else {
		// single start point
		start, err := strconv.Atoi(arg)
		if err != nil || start < 0 {
			return nil, -1, utils.ErrMalformed
		}
		return []int{start}, start, err
	}
}

// Returns the slices for the --bytes (e.g. -b 1,2,3)
func (t *Slice) getBytes(s string, startSl []int, end int) string {
	output.TraceLog("file.getBytes", fmt.Sprintf("s: %q, startSl: %#v, end: %d", s, startSl, end))
	if len(s) == 0 {
		return ""
	}
	result := ""
	// to check if change the end value or not, in case we have multiple starts means we have the comma as separator
	// and the end val has to be equals to start
	multipleStartVals := len(startSl) > 1

	for _, start := range startSl {
		if multipleStartVals {
			end = start
		}
		// edge cases: start > len(s)
		if start > len(s) {
			continue
		}
		if end > len(s) {
			end = len(s)
		}

		start--
		result += s[start:end]
		output.TraceLog("file.getBytes", fmt.Sprintf("computed start: %d, end: %d, result: %q", start, end, result))
	}
	return result
}

// Returns the slices for the --fields
func (t *Slice) getFields(s string, d string, startSl []int, end int) (string, error) {
	// pre checks
	if len(s) == 0 {
		return "", nil
	}

	if len(d) > 1 {
		// trimming from delimiter
		d = strings.Trim(d, " ")
		if len(d) > 1 {
			return "", utils.ErrSepMalformed
		}
	}
	output.TraceLog("file.getFields", fmt.Sprintf("computed delimiter: %q", d))

	// to check if change the end value or not, in case we have multiple starts means we have the comma as separator
	// and the end val has to be equals to start
	multipleStartVals := len(startSl) > 1
	result := ""
	slices := strings.Split(s, d)

	// cycle the start points
	for _, start := range startSl {
		if start > len(slices) {
			continue
		}
		if multipleStartVals {
			end = start
		}
		start--
		if end > len(slices) {
			end = len(slices)
		}

		for start < end {
			// if delimiter only the slices[start] == ""
			if slices[start] != "" {
				if len(result) > 0 {
					result += d + slices[start]
				} else {
					result += slices[start]
				}

			}
			start++
		}
	}

	return result, nil
}

// Returns the slices for the --chars
func (t *Slice) getChars(s string, startSl []int, end int) string {
	if len(s) == 0 {
		return ""
	}
	// to check if change the end value or not, in case we have multiple starts means we have the comma as separator
	// and the end val has to be equals to start
	multipleStartVals := len(startSl) > 1
	result := ""
	// runes creation
	runes := make([]int32, 0, len(s))
	for _, c := range s {
		runes = append(runes, c)
	}
	// cycle the start points
	for _, start := range startSl {
		if start > len(s) {
			continue
		}
		if multipleStartVals {
			end = start
		}

		start--
		if end > len(runes) {
			end = len(runes)
		}

		for start < end {
			result += fmt.Sprintf("%c", runes[start])
			start++
		}
	}

	output.TraceLog("file.getChars", fmt.Sprintf("final is '%s'", result))
	return result
}
