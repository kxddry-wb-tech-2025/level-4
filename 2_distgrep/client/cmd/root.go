package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "distgrep",
	Short: "A distributed grep tool",
	Long:  `A minimal grep implementation using Cobra CLI framework`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(grepCmd)
}
