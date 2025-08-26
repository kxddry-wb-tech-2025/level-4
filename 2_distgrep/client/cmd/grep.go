package cmd

import (
	"client/internal/models"
	"client/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	after        int
	before       int
	contextLines int
	invert       bool
	ignorecase   bool
	countOnly    bool
	fixedstring  bool
	printNumbers bool
	addrs        []string
)

var grepCmd = &cobra.Command{
	Use:   "grep [PATTERN] [FILE...]",
	Short: "Parse grep-like flags and addresses",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]
		files := args[1:]

		if len(files) == 0 {
			fmt.Println("No files provided")
			return
		}

		// Handle context flag which sets both before and after
		beforeCtx := before
		afterCtx := after
		if contextLines > 0 {
			beforeCtx = contextLines
			afterCtx = contextLines
		}

		flags := models.GrepFlags{
			FixedString:  fixedstring,
			PrintNumbers: printNumbers,
			IgnoreCase:   ignorecase,
			Invert:       invert,
			After:        afterCtx,
			Before:       beforeCtx,
			CountOnly:    countOnly,
		}

		if err := service.Run(pattern, files, addrs, flags); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	grepCmd.Flags().IntVarP(&after, "after", "A", 0, "Print NUM lines of trailing context after matching lines")
	grepCmd.Flags().IntVarP(&before, "before", "B", 0, "Print NUM lines of leading context before matching lines")
	grepCmd.Flags().IntVarP(&contextLines, "context", "C", 0, "Print NUM lines of output context")
	grepCmd.Flags().BoolVarP(&invert, "invert", "v", false, "Invert the sense of matching, to select non-matching lines")
	grepCmd.Flags().BoolVarP(&ignorecase, "ignore-case", "i", false, "Ignore case distinctions in patterns and data")
	grepCmd.Flags().BoolVarP(&countOnly, "count", "c", false, "Print only a count of selected lines per FILE")
	grepCmd.Flags().BoolVarP(&fixedstring, "fixed-string", "F", false, "Interpret PATTERN as a fixed string, not a regular expression")
	grepCmd.Flags().BoolVarP(&printNumbers, "print-numbers", "n", false, "Print line numbers with output lines")
	grepCmd.Flags().StringSliceVar(&addrs, "addrs", nil, "Comma-separated list of server addresses")
}
