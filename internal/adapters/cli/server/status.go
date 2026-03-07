package server

import (
	"fmt"
	"os"
	"syscall"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	cmd *cobra.Command
}

func NewStatusCmd() cli.Cmd {
	c := &statusCmd{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check if the REST API server is running",
		Run:   c.run,
	}
	c.cmd = cmd
	return c
}

func (c *statusCmd) Cmd() *cobra.Command {
	return c.cmd
}

func (c *statusCmd) run(cmd *cobra.Command, args []string) {
	pid, err := readPidFile()
	if err != nil {
		fmt.Println("Server is not running (PID file not found).")
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Server is not running (PID %d not found).\n", pid)
		removePidFile()
		return
	}

	// Send signal 0 to check if process exists
	if err := process.Signal(syscall.Signal(0)); err != nil {
		fmt.Printf("Server is not running (PID %d not found or accessible).\n", pid)
		removePidFile()
		return
	}

	fmt.Printf("Server is running with PID: %d\n", pid)
}
