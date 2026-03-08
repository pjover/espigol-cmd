package ports

// ImportService defines methods to import entities and resources.
type ImportService interface {
	// Reads the partners from the sourceAddress and stores the Partners.
	ImportPartners(sourceAddress string) (msg string, err error)

	// Reads expense forecasts from the provided CSV file and outputs them.
	ImportExpenseForecasts(sourceAddress string) (msg string, err error)
}
