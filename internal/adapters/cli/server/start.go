package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/cobra"
)

type startCmd struct {
	cmd    *cobra.Command
	server ports.Server
}

func NewStartCmd(server ports.Server) cli.Cmd {
	c := &startCmd{
		server: server,
	}
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

	// Start HTTP server in a separate goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- c.server.Start()
	}()

	// Block until signal is received or server fails
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		removePidFile()
		return fmt.Errorf("server crashed: %w", err)
	case <-sigChan:
		fmt.Println("\nStopping espigol server...")
		if err := c.server.Stop(context.Background()); err != nil {
			fmt.Printf("Error stopping server: %v\n", err)
		}
	}

	removePidFile()
	return nil
}
