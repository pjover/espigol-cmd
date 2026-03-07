package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandManager mocks the ports.CommandManager
type MockCommandManager struct {
	mock.Mock
}

func (m *MockCommandManager) AddCommand(cmd interface{}) {
	m.Called(cmd)
}

func (m *MockCommandManager) Execute() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// MockCmd mocks the cli.Cmd interface
type MockCmd struct {
	mock.Mock
	cmd *cobra.Command
}

func (m *MockCmd) Cmd() *cobra.Command {
	m.Called()
	return m.cmd
}

func TestNewServerCmd(t *testing.T) {
	mockStartCmd := new(MockCmd)
	mockStartCmd.cmd = &cobra.Command{Use: "start"}
	mockStartCmd.On("Cmd").Return(mockStartCmd.cmd)

	mockStopCmd := new(MockCmd)
	mockStopCmd.cmd = &cobra.Command{Use: "stop"}
	mockStopCmd.On("Cmd").Return(mockStopCmd.cmd)

	mockStatusCmd := new(MockCmd)
	mockStatusCmd.cmd = &cobra.Command{Use: "status"}
	mockStatusCmd.On("Cmd").Return(mockStatusCmd.cmd)

	// Create command
	sCmd := NewServerCmd(mockStartCmd, mockStopCmd, mockStatusCmd)

	cobraCmd := sCmd.Cmd()
	assert.NotNil(t, cobraCmd)
	assert.Equal(t, "server", cobraCmd.Use)

	// Validate subcommands were added
	assert.True(t, cobraCmd.HasSubCommands())
	subCmds := cobraCmd.Commands()
	assert.Len(t, subCmds, 3)

	var subCmdNames []string
	for _, sc := range subCmds {
		subCmdNames = append(subCmdNames, sc.Use)
	}

	assert.Contains(t, subCmdNames, "start")
	assert.Contains(t, subCmdNames, "stop")
	assert.Contains(t, subCmdNames, "status")

	// Ensure our mock cmdManager registers it
	mockCM := new(MockCommandManager)
	mockCM.On("AddCommand", sCmd).Return()

	mockCM.AddCommand(sCmd)
	mockCM.AssertExpectations(t)
}
