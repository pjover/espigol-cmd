package server

import (
	"fmt"
	"os"
	"syscall"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/spf13/cobra"
)

type stopCmd struct {
	cmd *cobra.Command
}

func NewStopCmd() cli.Cmd {
	c := &stopCmd{}
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the REST API server",
		RunE:  c.run,
	}
	c.cmd = cmd
	return c
}

func (c *stopCmd) Cmd() *cobra.Command {
	return c.cmd
}

func (c *stopCmd) run(cmd *cobra.Command, args []string) error {
	pid, err := readPidFile()
	if err != nil {
		return fmt.Errorf("server is not running or PID file not found")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		removePidFile()
		return fmt.Errorf("failed to find server process (PID %d)", pid)
	}

	// Send SIGTERM
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop server (PID %d): %w", pid, err)
	}

	fmt.Printf("Sent stop signal to server (PID %d).\n", pid)

	// Optional: wait a moment or let daemon's trap cleanup the file. We'll simply remove it here as a fallback or trust daemon does it.
	removePidFile()

	return nil
}
