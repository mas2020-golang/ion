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
		expected   string
	}{
		{
			"../../test-files/tail-1.txt",
			"2",
			"",
			"",
			"t",
		},
	}
	for _, c := range cases {
		slice := file.NewSlice()
		s, err := slice.Slice(c.file, c.sliceBytes, c.sliceChars, c.sliceCols)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if s != c.expected {
			t.Errorf("with %q, got %s, expected %s",
				c.file, s, c.expected)
		}
	}

}
