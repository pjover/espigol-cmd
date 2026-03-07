package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/spf13/cobra"
)

type startCmd struct {
	cmd *cobra.Command
}

func NewStartCmd() cli.Cmd {
	c := &startCmd{}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the REST API server",
		RunE:  c.run,
	}
	c.cmd = cmd
	return c
}

func (c *startCmd) Cmd() *cobra.Command {
	return c.cmd
}

func (c *startCmd) run(cmd *cobra.Command, args []string) error {
	// Check if already running
	if pid, err := readPidFile(); err == nil {
		// Check if process exists
		if process, err := os.FindProcess(pid); err == nil {
			if err := process.Signal(syscall.Signal(0)); err == nil {
				return fmt.Errorf("server is already running with PID %d", pid)
			}
		}
	}

	pid := os.Getpid()
	if err := writePidFile(pid); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	fmt.Printf("Starting espigol server (PID: %d)...\n", pid)

	// Block until signal is received to simulate daemon/server running
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// In Phase 2, this is where the HTTP server starts.
	// For now, it's just blocking.

	<-sigChan
	fmt.Println("\nStopping espigol server...")

	removePidFile()
	return nil
}
