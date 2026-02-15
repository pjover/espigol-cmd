package cli

import (
	"fmt"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/cobra"
)

type importExpenseForecastsCsvCmd struct {
	importService ports.ImportService
}

func NewImportExpenseForecastsCsvCmd(importService ports.ImportService) cli.Cmd {
	return importExpenseForecastsCsvCmd{
		importService: importService,
	}
}

func (i importExpenseForecastsCsvCmd) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "importarPrevisionsDespesa ruta/al/fitxer.csv",
		Short: "Importar previsions de despesa d'un fitxer CSV",
		Aliases: []string{
			"ipd",
			"importExpenseForecasts",
			"importarPrevisionsDespesa",
			"importar-previsions-despesa",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := i.importService.ImportExpenseForecasts(args[0])
			if err != nil {
				return err
			}
			fmt.Println(msg)
			return nil
		},
	}
}
