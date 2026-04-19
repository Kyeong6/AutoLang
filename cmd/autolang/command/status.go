package command

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Kyeong6/autolang/internal/config"
)

func newStatusCmd() *cobra.Command {
	var quiet bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show AutoLang proxy status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			running, pid := proxyStatus(cfg.Proxy.Port)

			if quiet {
				if !running {
					return fmt.Errorf("proxy not running")
				}
				return nil
			}

			if running {
				fmt.Printf("✓ AutoLang proxy is running (PID %d, port %d)\n", pid, cfg.Proxy.Port)
			} else {
				fmt.Println("✗ AutoLang proxy is not running.")
				fmt.Println("  Run: autolang start")
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Exit code only (for shell scripts)")
	return cmd
}

// proxyStatus checks if the proxy is reachable via its /health endpoint
// and reads the PID from the PID file.
func proxyStatus(port int) (running bool, pid int) {
	// Read PID from file
	home, _ := os.UserHomeDir()
	pidPath := filepath.Join(home, ".autolang", "autolang.pid")
	if data, err := os.ReadFile(pidPath); err == nil {
		pid, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}

	// Confirm via HTTP health check
	client := &http.Client{Timeout: 500 * time.Millisecond}
	url := fmt.Sprintf("http://localhost:%d/health", port)
	resp, err := client.Get(url)
	if err != nil {
		return false, pid
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK, pid
}
