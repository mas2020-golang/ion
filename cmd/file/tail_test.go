package file

import (
	"testing"
)

/*
Test function for the ls command of test module.
*/
func TestTail(t *testing.T) {
	cases := []struct {
		file     string
		rows     int
		expected int
	}{
		{
			"../../test/tail-1.txt",
			3,
			3,
		},
		{
			"../../test/tail-2.txt",
			1,
			1,
		},
		{
			"../../test/tail-2.txt",
			10,
			3,
		},
		{
			"../../test/tail-2.txt",
			0,
			0,
		},
		{
			"../../test/tail-3.txt",
			4,
			4,
		},
		{
			"../../test/tail-4.txt",
			3,
			2,
		},
		{
			"../../test/tail-4.txt",
			1,
			1,
		},
		{
			"../../test/tail-4.txt",
			2,
			2,
		},
		{
			"../../test/tail-6.txt",
			4,
			0,
		},
		{
			"../../test/tail-7.txt",
			2,
			2,
		},
		{
			"../../test/tail-7.txt",
			1,
			1,
		},
		{
			"../../test/tail-7.txt",
			5,
			3,
		},
	}

	for _, c := range cases {
		l, err := tail([]string{c.file}, c.rows)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if len(l) != c.expected {
			t.Errorf("with %s --rows %d, got %d, expected %d",
				c.file, c.rows, len(l), c.expected)
		}
	}

}
