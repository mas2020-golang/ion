package utils

import (
	"fmt"
	"io"
)

const (
	Reset        = "\033[0m"
	Bold         = "\033[1m"
	Blue         = "\033[34m"
	LightBlue   = "\033[94m"
	Orange       = "\033[38;5;167m"
	OkFlag       = "\u001B[92m✔\u001B[0m"
	ErrorFlag    = "\033[91m✘\033[0m"
	ErrorCmd     = "[\033[91mERROR\033[0m]"
)

// BoldS returns a string bold
func BoldS(t string) string {
	return fmt.Sprintf("%s%s%s", Bold, t, Reset )
}

// BoldOut bolds the passed argument
func BoldOut(t string, w io.Writer) {
	fmt.Fprintf(w, "%s%s%s", Bold, t, Reset )
}

// BlueS returns a blue string
func BlueS(t string) string {
	return fmt.Sprintf("%s%s%s", Blue, t, Reset )
}

// BlueOut write a blue text into the given writer
func BlueOut(t string, w io.Writer) {
	fmt.Fprintf(w, "%s%s%s", Blue, t, Reset )
}



