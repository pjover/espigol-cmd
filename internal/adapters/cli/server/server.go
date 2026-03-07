package server

import (
	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/spf13/cobra"
)

type serverCmd struct {
	cmd *cobra.Command
}

func NewServerCmd(startCmd, stopCmd, statusCmd cli.Cmd) cli.Cmd {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage the REST API server",
		Long:  "Commands to start, stop, and check the status of the REST API server.",
	}

	if startCmd != nil {
		cmd.AddCommand(startCmd.Cmd())
	}
	if stopCmd != nil {
		cmd.AddCommand(stopCmd.Cmd())
	}
	if statusCmd != nil {
		cmd.AddCommand(statusCmd.Cmd())
	}

	return &serverCmd{cmd: cmd}
}

func (s *serverCmd) Cmd() *cobra.Command {
	return s.cmd
}
