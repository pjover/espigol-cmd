package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/pjover/espigol/domain/services"
	"github.com/spf13/cobra"
)

var csvPath string

var importarCmd = &cobra.Command{
    Use:   "importar",
    Short: "Import resources",
}

var socisCmd = &cobra.Command{
    Use:   "socis",
    Short: "Import socis from CSV",
    RunE: func(cmd *cobra.Command, args []string) error {
        if csvPath == "" {
            return errors.New("missing --csv flag")
        }

        importer := services.NewCSVImporter()
        if err := importer.ImportSocisCSV(csvPath, os.Stdout); err != nil {
            return fmt.Errorf("import socis: %w", err)
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(importarCmd)
    importarCmd.AddCommand(socisCmd)

    socisCmd.Flags().StringVar(&csvPath, "csv", "", "path to csv file")
}
