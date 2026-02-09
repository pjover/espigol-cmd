package interfaces

// Importer defines methods to import entities and resources.
type Importer interface {
	
	// Reads the partners from the sourceAddress and stores the Partners.
	ImportPartners(sourceAddress string) error
}
