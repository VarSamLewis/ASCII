package cmd

import (
	"log/slog"

	asciiLog "ascii/ascii/src/internal/log"

	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "ascii",
	Short: "Convert images to ASCII art",
	Long:  "A tool to convert images to ASCII art, outputting to terminal or saving as an image file.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		asciiLog.Init(debug)
		slog.Debug("debug logging enabled")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
}
