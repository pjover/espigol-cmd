package cli

import (
	"fmt"

	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/cobra"
	"github.com/pjover/espigol/internal/adapters/cli"
)

type importPartnersCsvCmd struct {
	importService ports.ImportService
}

func NewImportPartnersCsvCmd(importService ports.ImportService) cli.Cmd {
	return importPartnersCsvCmd{
		importService: importService,
	}
}

func (i importPartnersCsvCmd) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "importarSocis ruta/al/fitxer.csv",
		Short: "Importar socis d'un fitxer CSV",
		Aliases: []string{
			"is",
			"importarSocis",
			"importar-socis",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := i.importService.ImportPartners(args[0])
			if err != nil {
				return err
			}
			fmt.Println(msg)
			return nil
		},
	}
}
