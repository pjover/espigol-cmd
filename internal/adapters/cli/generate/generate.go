package generate

import (
	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/spf13/cobra"
)

type generateCmd struct {
	cmd *cobra.Command
}

// NewGenerateCmd creates the parent "generate" command.
func NewGenerateCmd(subCmds ...cli.Cmd) cli.Cmd {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate PDF reports",
		Long:  "Commands to generate PDF reports for cooperative data.",
	}
	for _, sub := range subCmds {
		if sub != nil {
			cmd.AddCommand(sub.Cmd())
		}
	}
	return &generateCmd{cmd: cmd}
}

func (g *generateCmd) Cmd() *cobra.Command {
	return g.cmd
}
