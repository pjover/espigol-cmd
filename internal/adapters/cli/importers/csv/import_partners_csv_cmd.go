package cli

import (
	"fmt"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/cobra"
)

type importPartnersCsvCmd struct {
	importService ports.ImportService
	filePath      string
}

func NewImportPartnersCsvCmd(importService ports.ImportService) cli.Cmd {
	return &importPartnersCsvCmd{
		importService: importService,
	}
}

func (i *importPartnersCsvCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "partners",
		Short: "Import partners from a CSV file",
		Aliases: []string{
			"p",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := i.importService.ImportPartners(i.filePath)
			if err != nil {
				return err
			}
			fmt.Println(msg)
			return nil
		},
	}

	cmd.Flags().StringVarP(&i.filePath, "file", "f", "", "Path to the CSV file")
	cmd.MarkFlagRequired("file")

	return cmd
}
