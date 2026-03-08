package model

import (
	"fmt"
)

// ConfigReader is a minimal interface for reading float64 configuration values.
// It avoids an import cycle between model and ports.
type ConfigReader interface {
	GetFloat64(key string) float64
}

// ExpenseLimits holds the maximum amounts allowed for a grant year.
type ExpenseLimits struct {
	CurrentExpense    float64
	InvestmentExpense float64
	Total             float64
}

// LimitsForYear reads the grant limits for the given year from the config reader.
// Returns the limits and true if limits are defined, or a zero-value and false otherwise.
func LimitsForYear(year int, config ConfigReader) (ExpenseLimits, bool) {
	prefix := fmt.Sprintf("expenses.limits.%d", year)
	current := config.GetFloat64(prefix + ".current")
	investment := config.GetFloat64(prefix + ".investment")
	if current == 0 && investment == 0 {
		return ExpenseLimits{}, false
	}
	return ExpenseLimits{
		CurrentExpense:    current,
		InvestmentExpense: investment,
		Total:             current + investment,
	}, true
}
