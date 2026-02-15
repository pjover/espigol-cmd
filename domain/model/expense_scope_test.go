package model

import "testing"

func TestExpenseScopeConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      ExpenseScope
		expected string
	}{
		{"ExpenseScopeCommon", ExpenseScopeCommon, "Comú"},
		{"ExpenseScopeOliveSection", ExpenseScopeOliveSection, "Secció d'oliva"},
		{"ExpenseScopeLivestockSection", ExpenseScopeLivestockSection, "Secció de ramaderia"},
		{"ExpenseScopePartner", ExpenseScopePartner, "Soci"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.got))
			}
		})
	}
}

func TestExpenseScopeString(t *testing.T) {
	tests := []struct {
		scope    ExpenseScope
		expected string
	}{
		{ExpenseScopeCommon, "Comú"},
		{ExpenseScopeOliveSection, "Secció d'oliva"},
		{ExpenseScopeLivestockSection, "Secció de ramaderia"},
		{ExpenseScopePartner, "Soci"},
	}

	for _, tt := range tests {
		t.Run(string(tt.scope), func(t *testing.T) {
			got := tt.scope.String()
			if got != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}
