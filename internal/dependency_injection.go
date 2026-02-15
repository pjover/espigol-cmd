package internal

import (
	"log"

	"github.com/pjover/espigol/internal/adapters/cfg"
	"github.com/pjover/espigol/internal/adapters/cli"
	csv "github.com/pjover/espigol/internal/adapters/cli/importers/csv"
	"github.com/pjover/espigol/internal/domain/ports"
	importers "github.com/pjover/espigol/internal/domain/services/importers"
)

func InjectDependencies() ports.CommandManager {
	log.Print("Initializing dependency injection container...")

	// 1. Config Service
	configService := cfg.NewConfigService()

	// 2. Command Manager (CLI)
	cmdManager := cli.NewCommandManager(configService)

	// 3. Importers
	importersDI(cmdManager)

	return cmdManager
}

func importersDI(cmdManager ports.CommandManager) {
	importService := importers.NewCsvImporter()
	
	importPartnersCsvCmd := csv.NewImportPartnersCsvCmd(importService)
	cmdManager.AddCommand(importPartnersCsvCmd)
	
	importExpenseForecastsCsvCmd := csv.NewImportExpenseForecastsCsvCmd(importService)
	cmdManager.AddCommand(importExpenseForecastsCsvCmd)
}
