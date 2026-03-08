package cli

import (
	"fmt"

	"github.com/pjover/espigol/internal/domain"
	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/cobra"
)

type commandManager struct {
	configService ports.ConfigService
	rootCmd       *cobra.Command
}

func NewCommandManager(configService ports.ConfigService) ports.CommandManager {
	// RootCmd represents the base command when called without any subcommands
	title := fmt.Sprintf("sam v%s, Gestor de facturació de Hobbiton", domain.Version)
	rootCmd := &cobra.Command{
		Use:     "espigol",
		Short:   title,
		Long:    title + " (+ info: https://github.com/pjover/espigol)",
		Version: domain.Version,
	}

	cobra.OnInitialize(configService.Init)
	return commandManager{configService, rootCmd}
}

func (c commandManager) GetRootCmd() *cobra.Command {
	return c.rootCmd
}

func (c commandManager) AddCommand(cmd interface{}) {
	command := cmd.(Cmd)
	c.rootCmd.AddCommand(command.Cmd())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It only needs to happen once to the RootCmd.
func (c commandManager) Execute() []string {
	cobra.CheckErr(c.rootCmd.Execute())
	return []string{}
}
