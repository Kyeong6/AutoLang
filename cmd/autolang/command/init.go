package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Output shell integration script",
		Long: `Output shell integration script for eval.

Add this to your ~/.zshrc or ~/.bashrc:

  echo 'eval "$(autolang init)"' >> ~/.zshrc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(shellInitScript)
			return nil
		},
	}
}

const shellInitScript = `
# AutoLang shell integration
_autolang_init() {
  if ! autolang status --quiet 2>/dev/null; then
    autolang start --daemon > /dev/null 2>&1
  fi
  export ANTHROPIC_BASE_URL="http://localhost:7878"
}

_autolang_init
`
