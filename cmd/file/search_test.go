package file

import (
	"bytes"
	"testing"
)

func TestSearch(t *testing.T) {
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
			"",
			"Error: no files as argument\n",
			-1,
			true,
			nil,
		},
		{
			"../../test/test-files/search.txt",
			"Flags",
			"on '../../test/test-files/search.txt':\nFlags are:\n",
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
			//fmt.Println(c)
			cmd.SetArgs([]string{c.file})
			cmd.Execute()
			res := string(w.Bytes())
			if res != c.expected {
				t.Errorf("error: %v", res)
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
