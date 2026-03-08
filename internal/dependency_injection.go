package internal

import (
	"log"

	"github.com/pjover/espigol/internal/adapters/cfg"
	"github.com/pjover/espigol/internal/adapters/cli"
	csv "github.com/pjover/espigol/internal/adapters/cli/importers/csv"
	"github.com/pjover/espigol/internal/adapters/cli/server"
	httpAdapter "github.com/pjover/espigol/internal/adapters/http"
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

	// 3. DB Service (shared)
	dbService := mongodb.NewDbService(configService)

	// 4. Importers
	importersDI(cmdManager, dbService)

	// 5. Server Commands
	httpServer := httpAdapter.NewServer(configService, dbService)

	serverCmd := server.NewServerCmd(
		server.NewStartCmd(httpServer),
		server.NewStopCmd(),
		server.NewStatusCmd(),
	)
	cmdManager.AddCommand(serverCmd)

	return cmdManager
}

func importersDI(cmdManager ports.CommandManager, dbService ports.DbService) {
	importService := importers.NewCsvImporter(dbService)

	importPartnersCmd := csv.NewImportPartnersCsvCmd(importService)
	importExpenseForecastsCmd := csv.NewImportExpenseForecastsCsvCmd(importService)
	importCmd := csv.NewImportCmd(importPartnersCmd, importExpenseForecastsCmd)

	cmdManager.AddCommand(importCmd)
}
