package file

import (
	"bytes"
	"testing"
	// "github.com/sirupsen/logrus"
)

func TestReplace(t *testing.T) {
	//logrus.SetLevel(logrus.TraceLevel)
	cases := []struct {
		file     string
		expected string
		flags    []string
	}{
		{
			"../../test-files/replace.txt",
			"--no-colors: no highlight colors in CHANGE TEXT output",
			[]string{"-p", "\\bthe\\b", "-s", "CHANGE TEXT", "-o"},
		},
	}

	for _, c := range cases {
		cmd := NewReplaceCmd()
		w := bytes.NewBuffer(nil)
		cmd.SetOut(w)
		cmd.SetErr(w)
		args := []string{c.file}
		if c.flags != nil {
			args = append(args, c.flags...)
		}
		cmd.SetArgs(args)
		cmd.Execute()
		res := w.String()
		if res != c.expected {
			t.Errorf("\nexpected:\n%v\ngot:\n%v", c.expected, res)
			t.Logf("len res: %d, len expected: %d\n", len(res), len(c.expected))
			t.Logf("bytes res:\n%v\nbytes expected:\n%v\n", []byte(res), []byte(c.expected))
		}
	}
}
