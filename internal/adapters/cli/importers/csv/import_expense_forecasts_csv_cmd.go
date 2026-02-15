package cli

import (
	"fmt"

	"github.com/pjover/espigol/internal/adapters/cli"
	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/cobra"
)

type importExpenseForecastsCsvCmd struct {
	importService ports.ImportService
	filePath      string
}

func NewImportExpenseForecastsCsvCmd(importService ports.ImportService) cli.Cmd {
	return &importExpenseForecastsCsvCmd{
		importService: importService,
	}
}

func (i *importExpenseForecastsCsvCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expense-forecasts",
		Short: "Import expense forecasts from a CSV file",
		Aliases: []string{
			"ef",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := i.importService.ImportExpenseForecasts(i.filePath)
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
