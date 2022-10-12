package file

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	os.Setenv("ION_LOGLEVEL", "6")
	cases := []struct {
		file     string
		pattern  string
		expected string
		count    int
		err      bool
		flags    []string
	}{
		{
			"../../test/test-files/search.txt",
			"Flags",
			"=> on '../../test/test-files/search.txt':\nFlags are:\n",
			-1,
			false,
			[]string{"--no-colors"},
		},
	}

	cmd := NewSearchCmd()
	for _, c := range cases {
		if c.err {
			w := bytes.NewBuffer(nil)
			cmd.SetErr(w)
			fmt.Println(c)
			cmd.SetArgs([]string{c.file})
			err := cmd.Execute()
			fmt.Println(err)
			res := string(w.Bytes())
			if res != c.expected {
				t.Errorf("error: %v", res)
				t.Logf("len res: %d, len expected: %d\n", len(res), len(c.expected))
			}
		} else {
			w := bytes.NewBuffer(nil)
			cmd.SetOut(w)
			//fmt.Println(c)
			args := []string{c.pattern, c.file}
			if c.flags != nil {
				for _, f := range c.flags {
					args = append(args, f)
				}
			}
			cmd.SetArgs(args)
			cmd.Execute()
			res := string(w.Bytes())
			if res != c.expected {
				t.Errorf("expected:\n%v\ngot:\n%v", c.expected, res)
				t.Logf("len res: %d, len expected: %d\n", len(res), len(c.expected))
			}
		}
	}
}
