package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show token savings for current session",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Task 09 — read stats from proxy via socket or file
			fmt.Println("No session data available.")
			fmt.Println("Start a session with: autolang start")
			return nil
		},
	}
}
