package internal

import (
	"log"

	"github.com/pjover/espigol/internal/adapters/cfg"
	"github.com/pjover/espigol/internal/adapters/cli"
	csv "github.com/pjover/espigol/internal/adapters/cli/importers/csv"
	"github.com/pjover/espigol/internal/adapters/mongodb"
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
	importersDI(cmdManager, configService)

	return cmdManager
}

func importersDI(cmdManager ports.CommandManager, configService ports.ConfigService) {
	dbService := mongodb.NewDbService(configService)
	importService := importers.NewCsvImporter(dbService)

	importPartnersCmd := csv.NewImportPartnersCsvCmd(importService)
	importExpenseForecastsCmd := csv.NewImportExpenseForecastsCsvCmd(importService)
	importCmd := csv.NewImportCmd(importPartnersCmd, importExpenseForecastsCmd)

	cmdManager.AddCommand(importCmd)
}
