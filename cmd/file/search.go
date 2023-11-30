/*
Copyright © 2020 @mas2020 andrea.genovesi@gmail.com
*/
package file

import (
	"github.com/mas2020-golang/ion/packages/file"
	"github.com/spf13/cobra"
)

var (
	nocolors, countLines, countPattern, onlyMatch, invert bool
	insensitive, onlyResult, onlyFilename, isDir          bool
	recursive                                             bool
	cmd                                                   *cobra.Command
	before, after, argsN, currentLine, level              int
	prLines                                               map[int]bool // save the printed lines
	mLines                                                map[int]string
	mLinesMatch                                           map[int][][]int
)

func NewSearchCmd() *cobra.Command {
	cmd = &cobra.Command{
		Args: cobra.MinimumNArgs(1),
		Use:  "search <PATTERN> <PATH> [...PATH]",
		Example: `# search this in the demo-file
$ ion search "this" demo-file`,
		Short: "Search for the given pattern into the standard input or one or more files",
		Long: `The command searches for the given pattern. The command can read
directly from the standard input, one or more files or directories passed an argument. The pattern is highlighted with the red color.`,
		Run: func(cmd *cobra.Command, args []string) {
			argsN = len(args)
			output := cmd.OutOrStdout()
			searcher := file.NewSearcher(countLines, countPattern, onlyMatch, nocolors, invert, insensitive, onlyResult,
				onlyFilename, recursive, before, after, output)
			searcher.Search(args)
		},
	}
	// cmd.SetOut(os.Stdout)
	// flags
	cmd.Flags().BoolVarP(&countLines, "count-lines", "l", false, "shows only how many lines match with the pattern")
	cmd.Flags().BoolVarP(&countPattern, "count-pattern", "p", false, "shows only how many time a pattern is in match")
	cmd.Flags().BoolVarP(&onlyMatch, "only-match", "m", false, "shows only the substring that match, not the entire line")
	cmd.Flags().BoolVarP(&nocolors, "no-colors", "n", false, "no colors are printed onto the standard output")
	cmd.Flags().BoolVarP(&invert, "invert", "t", false, "shows the lines that doesn't match with the pattern")
	cmd.Flags().BoolVarP(&insensitive, "insensitive", "i", false, "to match with no case sensitivity")
	cmd.Flags().BoolVarP(&onlyResult, "only-result", "r", false, "if there is at least one match it returns 1, otherwise 0")
	cmd.Flags().IntVarP(&before, "before", "B", 0, "shows also the NUMBER of lines before the match")
	cmd.Flags().IntVarP(&after, "after", "A", 0, "shows also the NUMBER of lines after the match")
	cmd.Flags().BoolVarP(&onlyFilename, "only-filename", "f", false, "shows only the filename if a pattern matches one or several times. If the pattern doesn't match, no output is given.")
	cmd.Flags().BoolVarP(&recursive, "recursive", "d", false, "if true and the PATH is a folder searches in all the sub folders")
	return cmd
}
