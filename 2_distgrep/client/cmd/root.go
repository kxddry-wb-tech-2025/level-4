package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grep [PATTERN] [FILE...]",
	Short: "A distributed grep tool",
	Long:  `A minimal grep implementation using Cobra CLI framework`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runGrep(args)
	},
}

// Execute is the main entry point for the CLI application
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(grepCmd)
	rootCmd.Flags().IntVarP(&after, "after", "A", 0, "Print NUM lines of trailing context after matching lines")
	rootCmd.Flags().IntVarP(&before, "before", "B", 0, "Print NUM lines of leading context before matching lines")
	rootCmd.Flags().IntVarP(&contextLines, "context", "C", 0, "Print NUM lines of output context")
	rootCmd.Flags().BoolVarP(&invert, "invert", "v", false, "Invert the sense of matching, to select non-matching lines")
	rootCmd.Flags().BoolVarP(&ignorecase, "ignore-case", "i", false, "Ignore case distinctions in patterns and data")
	rootCmd.Flags().BoolVarP(&countOnly, "count", "c", false, "Print only a count of selected lines per FILE")
	rootCmd.Flags().BoolVarP(&fixedstring, "fixed-string", "F", false, "Interpret PATTERN as a fixed string, not a regular expression")
	rootCmd.Flags().BoolVarP(&printNumbers, "print-numbers", "n", false, "Print line numbers with output lines")
	rootCmd.Flags().StringSliceVar(&addrs, "addrs", nil, "Comma-separated list of server addresses")
	rootCmd.Flags().IntVar(&quorum, "quorum", 0, "Quorum of successful servers required (default: majority)")
}
