package cli

import (
	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/spf13/cobra"
)

type importCmd struct {
	cmd *cobra.Command
}

func NewImportCmd(importPartnersCmd, importExpenseForecastsCmd cli.Cmd) cli.Cmd {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data from a file",
		Aliases: []string{
			"i",
		},
	}

	cmd.AddCommand(importPartnersCmd.Cmd())
	cmd.AddCommand(importExpenseForecastsCmd.Cmd())

	return &importCmd{cmd: cmd}
}

func (i *importCmd) Cmd() *cobra.Command {
	return i.cmd
}
