package file

import (
	"testing"

	"github.com/mas2020-golang/ion/packages/file"
	"github.com/mas2020-golang/ion/packages/utils"
)

/*
Test function for the slice command.
*/
func TestBytesAndChars(t *testing.T) {
	casesNew := []struct {
		file     string
		input    [][]string
		expected [][]string
		err      []error
	}{
		{
			// --bytes testing
			"../../test-files/slice.txt",
			[][]string{
				{"2", "", "", ""},
				{"1:4", "", "", ""},
				{"1:-4", "", "", ""},
				{"-1", "", "", ""},
				{"1:", "", "", ""},
				{":", "", "", ""},
				{"", "", "", ""},
				{"-1", "2", "", ""},
				{"18,2", "", "", ""}},
			[][]string{
				{"", "-", " ", "E", "\xb8", "-"},
				{"", "--3 ", "A B ", "TEST", "世\xe7","---"},
				{"", "", "", "", "","---"},
				{"", "", "", "", "","---"},
				{"", "--3 -10", "A B C", "TEST", "世界 field3 44","---"},
				{"", "--3", "A B C", "TEST", "","---"},
				{"", "--3", "A B C", "TEST", "","---"},
				{"", "", "", "", "","---"},
				{"", "-", " ", "E", "\xb8","-"}},
			[]error{nil, nil, utils.ErrMalformed, utils.ErrMalformed, nil, utils.ErrMalformed, utils.ErrMalformed, utils.ErrMalformed, nil},
		},
		{
			// --chars testing
			"../../test-files/slice.txt",
			[][]string{
				{"", "2", "", ""},
				{"", "1:3", "", ""},
				{"", "1:-3", "", ""},
				{"", "100", "", ""},
				{"", "1:", "", ""},
				{"", "18,2", "", ""},
			},
			[][]string{
				{"", "-", " ", "E", "界","-"},
				{"", "--3", "A B", "TES", "世界 ","---"},
				{"", "", "", "", "","---"},
				{"", "", "", "", "",""},
				{"", "--3 -10", "A B C", "TEST", "世界 field3 44","---"},
				{"", "-", " ", "E", "界","-"},
			},
			[]error{nil, nil, utils.ErrMalformed, nil, nil, nil},
		},
	}

	// test cases execution
	for _, c := range casesNew {
		slice := file.NewSlice()
		for i := 0; i < len(c.input); i++ {
			values, err := slice.Slice([]string{c.file}, c.input[i][0], c.input[i][1], c.input[i][2], c.input[i][3])
			if err != c.err[i] {
				t.Errorf("with %q, got >>%v<<, expected >>%v<<",
					c.file, err, c.err[i])
			}

			for i2, v := range values {
				if v != c.expected[i][i2] {
					t.Errorf("with [-b %q, -c %q, -f %q], got '%s', expected '%s'",
						c.input[i][0], c.input[i][1], c.input[i][2], v, c.expected[i][i2])
				}
			}
		}
	}
}

func TestFields(t *testing.T) {
	casesNew := []struct {
		file     string
		input    [][]string
		expected [][]string
		err      []error
	}{
		{
			// --fields testing
			"../../test-files/slice.txt",
			[][]string{
				{"", "", "2", " "},
				{"", "", "2:3", "-"},
				{"", "", "2:100", " "},
			},
			[][]string{
				{"", "-10", "B", "", "field3",""},
				{"", "3 ", "", "", "",""},
				{"", "-10", "B C", "", "field3 44",""},
			},
			[]error{nil, nil, nil},
		},
	}

	// test cases execution
	for _, c := range casesNew {
		slice := file.NewSlice()
		for i := 0; i < len(c.input); i++ {
			values, err := slice.Slice([]string{c.file}, c.input[i][0], c.input[i][1], c.input[i][2], c.input[i][3])
			if err != c.err[i] {
				t.Errorf("with %q, got >>%v<<, expected >>%v<<",
					c.file, err, c.err[i])
			}

			for i2, v := range values {
				if v != c.expected[i][i2] {
					t.Errorf("with [-b %q, -c %q, -f %q, d %q], got '%s', expected '%s'",
						c.input[i][0], c.input[i][1], c.input[i][2], c.input[i][3], v, c.expected[i][i2])
				}
			}
		}
	}
}
