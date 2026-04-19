package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/Kyeong6/autolang/internal/config"
	"github.com/Kyeong6/autolang/internal/proxy"
	"github.com/Kyeong6/autolang/internal/translate"
)

func newStartCmd() *cobra.Command {
	var daemon bool

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the AutoLang proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if daemon {
				return startDaemon(cfg)
			}

			tr, err := translate.New(cfg)
			if err != nil {
				fmt.Printf("⚠ Translation disabled: %v\n", err)
				fmt.Println("  Set AUTOLANG_PROVIDER and AUTOLANG_API_KEY to enable translation.")
				tr = nil
			} else {
				fmt.Printf("✓ Translation provider: %s\n", tr.Name())
			}

			fmt.Printf("AutoLang proxy starting on port %d\n", cfg.Proxy.Port)
			fmt.Printf("Set: export ANTHROPIC_BASE_URL=http://localhost:%d\n", cfg.Proxy.Port)
			fmt.Println("Press Ctrl+C to stop.")

			p := proxy.New(cfg, tr)
			return p.Start()
		},
	}

	cmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "Run proxy in background")
	return cmd
}

// startDaemon re-execs the current binary in a detached process and writes a PID file.
func startDaemon(cfg *config.Config) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot locate executable: %w", err)
	}

	dir, err := autolangDir()
	if err != nil {
		return err
	}

	logPath := filepath.Join(dir, "autolang.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("cannot open log file: %w", err)
	}

	proc, err := os.StartProcess(exe, []string{exe, "start"}, &os.ProcAttr{
		Files: []*os.File{os.Stdin, logFile, logFile},
		Sys:   &syscall.SysProcAttr{Setsid: true},
	})
	if err != nil {
		return fmt.Errorf("cannot start background process: %w", err)
	}

	pidPath := filepath.Join(dir, "autolang.pid")
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(proc.Pid)), 0o644); err != nil {
		return fmt.Errorf("cannot write PID file: %w", err)
	}

	// Release the child so it keeps running after the parent exits.
	proc.Release() //nolint:errcheck

	fmt.Printf("AutoLang proxy started (PID %d, port %d)\n", proc.Pid, cfg.Proxy.Port)
	fmt.Printf("Log: %s\n", logPath)
	return nil
}

func autolangDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot locate home directory: %w", err)
	}
	dir := filepath.Join(home, ".autolang")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create ~/.autolang: %w", err)
	}
	return dir, nil
}
