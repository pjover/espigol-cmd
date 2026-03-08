package ports

// GenerateService defines operations that generate PDF reports.
type GenerateService interface {
	// ExpenseForecastReport generates the expense forecast PDF for the given year.
	// Returns (hasNegativeRemainder, message, error).
	ExpenseForecastReport(year int) (bool, string, error)
}
