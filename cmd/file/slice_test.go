package file

import (
	"testing"

	"github.com/mas2020-golang/ion/packages/file"
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
	}{
		//TODO: add more use cases to the application
		{
			"../../test-files/slice.txt",
			"2", // --bytes testing
			"",
			"",
			[]string{"", "-", " ", "E"},
		},
	}
	for _, c := range cases {
		slice := file.NewSlice()
		values, err := slice.Slice(c.file, c.sliceBytes, c.sliceChars, c.sliceCols)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		for i, v := range values {
			if v != c.expected[i] {
				t.Errorf("with %q, got %q, expected %q",
					c.file, v, c.expected[i])
			}
		}

	}

}
