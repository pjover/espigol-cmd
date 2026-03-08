package model

import (
	"testing"
	"time"
)

func testForecastWithSubtype(subtype ExpenseSubtype, plannedDate time.Time) *ExpenseForecast {
	partner := NewPartner(1, "Pere", "Jover", "43030928K", "pjover@gmail.com", "+34644421965",
		Producer, 32425, true, false, time.Date(2023, 4, 21, 0, 0, 0, 0, time.UTC))
	return NewExpenseForecast(
		1, *partner, "Test concept", "Test description",
		1000.0, plannedDate, subtype, ExpenseScopeCommon, []string{}, time.Now(),
	)
}

func TestExpenseForecastYear(t *testing.T) {
	tests := []struct {
		name         string
		plannedDate  time.Time
		expectedYear int
	}{
		{"2026 date", time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), 2026},
		{"2025 date", time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), 2025},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef := testForecastWithSubtype(ExpenseSubtypeA1, tt.plannedDate)
			if ef.Year() != tt.expectedYear {
				t.Errorf("expected year %d, got %d", tt.expectedYear, ef.Year())
			}
		})
	}
}

func TestExpenseForecastExpenseCategory(t *testing.T) {
	plannedDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		subtype  ExpenseSubtype
		expected ExpenseCategory
	}{
		{"A1 is current", ExpenseSubtypeA1, ExpenseCategoryCurrent},
		{"A6 is current", ExpenseSubtypeA6, ExpenseCategoryCurrent},
		{"B1 is investment", ExpenseSubtypeB1, ExpenseCategoryInvestment},
		{"B5 is investment", ExpenseSubtypeB5, ExpenseCategoryInvestment},
		{"C1 is current", ExpenseSubtypeC1, ExpenseCategoryCurrent},
		{"C2 is current", ExpenseSubtypeC2, ExpenseCategoryCurrent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef := testForecastWithSubtype(tt.subtype, plannedDate)
			got := ef.ExpenseCategory()
			if got != tt.expected {
				t.Errorf("expected category %s, got %s", tt.expected, got)
			}
		})
	}
}
