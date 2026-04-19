package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the AutoLang proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := autolangDir()
			if err != nil {
				return err
			}

			pidPath := filepath.Join(dir, "autolang.pid")
			data, err := os.ReadFile(pidPath)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("AutoLang proxy is not running.")
					return nil
				}
				return fmt.Errorf("cannot read PID file: %w", err)
			}

			pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
			if err != nil {
				return fmt.Errorf("invalid PID file: %w", err)
			}

			proc, err := os.FindProcess(pid)
			if err != nil {
				return fmt.Errorf("cannot find process %d: %w", pid, err)
			}

			if err := proc.Signal(syscall.SIGTERM); err != nil {
				return fmt.Errorf("cannot stop process %d: %w", pid, err)
			}

			os.Remove(pidPath) //nolint:errcheck
			fmt.Printf("AutoLang proxy stopped (PID %d).\n", pid)
			return nil
		},
	}
}
