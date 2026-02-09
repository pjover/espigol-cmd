package cmd

import (
	"errors"
	"fmt"

	"github.com/pjover/espigol/cfg"
	"github.com/spf13/cobra"
)

var csvPath string

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

func init() {
	rootCmd.AddCommand(importEntitiesCmd)
	importEntitiesCmd.AddCommand(importPartnersCmd)

	importPartnersCmd.Flags().StringVar(&csvPath, "csv", "", "path to csv file")
}
