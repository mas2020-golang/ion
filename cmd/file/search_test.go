package file

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	cases := []struct {
		file     string
		pattern  string
		expected string
		count    int
		err      bool
	}{
		{
			"../../test/test-files/search.txt",
			"",
			"Error: no files as argument\n",
			-1,
			true,
		},
		{
			"../../test/test-files/search.txt",
			"Flags",
			"-- on '../../test/test-files/search.txt':\nFlags are:",
			-1,
			false,
		},
	}

	cmd := NewSearchCmd()
	w := bytes.NewBuffer(nil)
	cmd.SetErr(w)
	for _, c := range cases {
		cmd.SetArgs([]string{c.file})
		cmd.Execute()
		if c.err {
			res := string(w.Bytes())
			fmt.Println(res)
			fmt.Println(c.expected)
			if res != c.expected {
				t.Errorf("error: %v", res)
			}
		} else {
			// TODO: continue from here
		}
	}
}
