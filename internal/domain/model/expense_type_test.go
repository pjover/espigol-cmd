package model

import "testing"

func TestExpenseTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      ExpenseType
		expected string
	}{
		{"ExpenseTypeA", ExpenseTypeA, "[a] Despeses corrents de caràcter fungible o temporals"},
		{"ExpenseTypeB", ExpenseTypeB, "[b] Despeses de caràcter permanent o d'inversió"},
		{"ExpenseTypeC", ExpenseTypeC, "[c] Despeses indirectes o d'estructura de l'entitat beneficiària"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.got))
			}
		})
	}
}

func TestExpenseTypeString(t *testing.T) {
	tests := []struct {
		expenseType ExpenseType
		expected    string
	}{
		{ExpenseTypeA, "[a] Despeses corrents de caràcter fungible o temporals"},
		{ExpenseTypeB, "[b] Despeses de caràcter permanent o d'inversió"},
		{ExpenseTypeC, "[c] Despeses indirectes o d'estructura de l'entitat beneficiària"},
	}

	for _, tt := range tests {
		t.Run(string(tt.expenseType), func(t *testing.T) {
			got := tt.expenseType.String()
			if got != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}

func TestExpenseTypeCategory(t *testing.T) {
	tests := []struct {
		expenseType ExpenseType
		expected    ExpenseCategory
	}{
		{ExpenseTypeA, ExpenseCategoryCurrent},
		{ExpenseTypeB, ExpenseCategoryInvestment},
		{ExpenseTypeC, ExpenseCategoryCurrent},
	}

	for _, tt := range tests {
		t.Run(string(tt.expenseType), func(t *testing.T) {
			got := tt.expenseType.Category()
			if got != tt.expected {
				t.Errorf("Expected category '%s', got '%s'", tt.expected, got)
			}
		})
	}
}

func TestExpenseCategoryString(t *testing.T) {
	tests := []struct {
		category ExpenseCategory
		expected string
	}{
		{ExpenseCategoryCurrent, "Despesa corrent"},
		{ExpenseCategoryInvestment, "Despesa d'inversió"},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			got := tt.category.String()
			if got != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}
