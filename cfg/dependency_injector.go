package cfg

import (
	"log"
	"sync"

	"github.com/pjover/espigol/domain/interfaces"
	importers "github.com/pjover/espigol/domain/services/importers"
)

type DependencyInjector struct {
	importer interfaces.Importer
}

// Singleton instance of DiContainer
var (
	instance *DependencyInjector
	once     sync.Once
)

func DI() *DependencyInjector {
	once.Do(func() {
		log.Print("Initializing dependency injection container...")
		var importer = importers.NewCSVImporter()
		instance = &DependencyInjector{
			importer: importer,
		}
	})
	return instance
}

func (di DependencyInjector) Importer() interfaces.Importer {
	return di.importer
}
