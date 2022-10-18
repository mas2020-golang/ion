package file

import (
	"bytes"
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
		flags    []string
	}{
		{
			"../../test/test-files/search.txt",
			"Flags",
			"Flags are:\n",
			-1,
			[]string{"--no-colors"},
		},
		{
			"../../test/test-files/search.txt",
			"<NUMBER>",
			"--after <NUMBER>: shows also the NUMBER of lines after the match\n--before <NUMBER>: shows also the NUMBER of lines before the match\n",
			-1,
			[]string{"--no-colors"},
		},
		{
			"../../test/test-files/search.txt",
			"<NUMBER>",
			"2\n",
			-1,
			[]string{"--no-colors", "--count-lines"},
		},
		{
			"../../test/test-files/search.txt",
			"<NUMBER>",
			"2\n",
			-1,
			[]string{"--no-colors", "--count-pattern"},
		},
		{
			"../../test/test-files/search.txt",
			"",
			"23\n",
			-1,
			[]string{"--no-colors", "--count-lines"},
		},
		{
			"../../test/test-files/search.txt",
			"number",
			"5\n",
			-1,
			[]string{"--no-colors", "-i", "-p"},
		},
	}

	//cmd := NewSearchCmd()
	for _, c := range cases {
		cmd := NewSearchCmd()
		w := bytes.NewBuffer(nil)
		cmd.SetOut(w)
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
			t.Logf("bytes res:\n%v\nbytes expected:\n%v\n", []byte(res), []byte(c.expected))
		}
	}
}
