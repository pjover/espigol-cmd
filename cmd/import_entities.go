package cmd

import (
	"errors"
	"fmt"

	"github.com/pjover/espigol/cfg"
	"github.com/spf13/cobra"
)

var csvPath string
var expenseCSV string

var importEntitiesCmd = &cobra.Command{
	Use:   "import",
	Short: "Import entities",
}

var importPartnersCmd = &cobra.Command{
	Use:   "partners",
	Short: "Import partners from CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		if csvPath == "" {
			return errors.New("missing --csv flag")
		}

		csvImporter := cfg.DI().Importer()
		if err := csvImporter.ImportPartners(csvPath); err != nil {
			return fmt.Errorf("import partners: %w", err)
		}

		return nil
	},
}

var importExpenseForecastsCmd = &cobra.Command{
	Use:   "expense-forecasts",
	Short: "Import expense forecasts from CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		if expenseCSV == "" {
			return errors.New("missing --csv flag")
		}

		csvImporter := cfg.DI().Importer()
		if err := csvImporter.ImportExpenseForecasts(expenseCSV); err != nil {
			return fmt.Errorf("import expense forecasts: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importEntitiesCmd)
	importEntitiesCmd.AddCommand(importPartnersCmd)
	importEntitiesCmd.AddCommand(importExpenseForecastsCmd)

	importPartnersCmd.Flags().StringVar(&csvPath, "csv", "", "path to csv file")
	importExpenseForecastsCmd.Flags().StringVar(&expenseCSV, "csv", "", "path to csv file")
}
