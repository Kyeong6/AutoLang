package command

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/Kyeong6/autolang/internal/config"
)

func newStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show token savings for current session",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			url := fmt.Sprintf("http://localhost:%d/stats", cfg.Proxy.Port)
			client := &http.Client{Timeout: 500 * time.Millisecond}
			resp, err := client.Get(url)
			if err != nil {
				fmt.Println("AutoLang proxy is not running.")
				fmt.Println("  Run: autolang start")
				return nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("cannot read stats: %w", err)
			}
			fmt.Println(string(body))
			return nil
		},
	}
}
