package file

import (
	"testing"

	"github.com/mas2020-golang/ion/packages/file"
	"github.com/mas2020-golang/ion/packages/utils"
)

/*
Test function for the slice command.
*/
func TestSlice(t *testing.T) {
	cases := []struct {
		file       string
		sliceBytes string
		sliceChars string
		sliceCols  string
		expected   []string
		err        error
	}{
		//TODO: add more use cases to the application
		{
			"../../test-files/slice.txt",
			"2", // --bytes testing
			"",
			"",
			[]string{"", "-", " ", "E"},
			nil,
		},
		{
			"../../test-files/slice.txt",
			"1:4", // --bytes testing
			"",
			"",
			[]string{"", "--3", "A B ", "TEST"},
			nil,
		},
		{
			"../../test-files/slice.txt",
			"1:-4", // --bytes testing
			"",
			"",
			[]string{"", "", "", ""},
			nil,
		},
		{
			"../../test-files/slice.txt",
			"-1:", // --bytes testing
			"",
			"",
			[]string{"", "", "", ""},
			nil,
		},
		{
			"../../test-files/slice.txt",
			"1:", // --bytes testing
			"",
			"",
			[]string{"", "--3", "A B C", "TEST"},
			nil,
		},
		{
			"../../test-files/slice.txt",
			":", // --bytes testing
			"",
			"",
			[]string{"", "--3", "A B C", "TEST"},
			utils.ErrMalformed,
		},
		{
			"../../test-files/slice.txt",
			"",
			"",
			"",
			[]string{"", "--3", "A B C", "TEST"},
			utils.ErrMalformed,
		},
		{
			"../../test-files/slice.txt",
			"-1", // --bytes testing
			"",
			"",
			[]string{"", "", "", ""},
			nil,
		},
	}
	for _, c := range cases {
		slice := file.NewSlice()
		values, err := slice.Slice(c.file, c.sliceBytes, c.sliceChars, c.sliceCols)
		if err != c.err {
			t.Errorf("with %q, got >>%v<<, expected >>%v<<",
				c.file, err, c.err)
		}

		for i, v := range values {
			if v != c.expected[i] {
				t.Errorf("with %q, got %q, expected %q",
					c.file, v, c.expected[i])
			}
		}

	}

}
