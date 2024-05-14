package file

import (
	"testing"
)

func TestReplace(t *testing.T) {
	cases := []struct {
		file     string
		pattern  string
		expected string
		count    int
		flags    []string
	}{
		{
			"../../test-files/search.txt",
			"Flags",
			"Flags are:\n",
			-1,
			[]string{"-p '\bthe\b'"},
		},
	}
}
