package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "autolang",
	Short: "Write in Korean. Get English-quality AI responses.",
	Long: `AutoLang is a transparent local proxy that sits between your CLI and AI APIs.
It automatically detects Korean input, translates it to English before sending,
and returns responses in Korean — with zero changes to your workflow.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	rootCmd.AddCommand(
		newStartCmd(),
		newStopCmd(),
		newStatusCmd(),
		newStatsCmd(),
		newInitCmd(),
	)
}
