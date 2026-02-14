package model

import (
	"testing"
	"time"
)

// Helper function to create a test Partner
func createTestPartner(id int) Partner {
	return *NewPartner(
		id, "John", "Doe", "12345678A", "john@example.com",
		"+34600000000", Producer, 1, true, false,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	)
}

func TestNewExpenseForecast(t *testing.T) {
	addedOn := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	plannedDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	attachments := []string{"receipt.pdf"}
	partner := createTestPartner(1)
	
	ef := NewExpenseForecast(
		1, partner, "Fertilizer", "Organic fertilizer for olives",
		150.50, plannedDate, ExpenseSubtypeA6, ExpenseScopeOliveSection,
		attachments, addedOn,
	)
	
	if ef.ID() != 1 {
		t.Errorf("Expected ID 1, got %d", ef.ID())
	}
	if ef.Partner().ID() != 1 {
		t.Errorf("Expected Partner ID 1, got %d", ef.Partner().ID())
	}
	if ef.Concept() != "Fertilizer" {
		t.Errorf("Expected concept 'Fertilizer', got '%s'", ef.Concept())
	}
}

func TestExpenseForecastGetters(t *testing.T) {
	addedOn := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	plannedDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	attachments := []string{"receipt.pdf"}
	partner := createTestPartner(1)
	
	ef := NewExpenseForecast(
		1, partner, "Fertilizer", "Organic fertilizer for olives",
		150.50, plannedDate, ExpenseSubtypeA6, ExpenseScopeOliveSection,
		attachments, addedOn,
	)
	
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Concept", ef.Concept(), "Fertilizer"},
		{"Description", ef.Description(), "Organic fertilizer for olives"},
		{"GrossAmount", ef.GrossAmount(), 150.50},
		{"PlannedDate", ef.PlannedDate(), plannedDate},
		{"ExpenseSubtype", ef.ExpenseSubtype(), ExpenseSubtypeA6},
		{"Scope", ef.Scope(), ExpenseScopeOliveSection},
		{"AddedOn", ef.AddedOn(), addedOn},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.got)
			}
		})
	}
}

func TestExpenseForecastAttachmentsImmutability(t *testing.T) {
	addedOn := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	plannedDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	attachments := []string{"receipt.pdf"}
	partner := createTestPartner(1)
	
	ef := NewExpenseForecast(
		1, partner, "Fertilizer", "Organic fertilizer for olives",
		150.50, plannedDate, ExpenseSubtypeA6, ExpenseScopeOliveSection,
		attachments, addedOn,
	)
	
	// Get attachments and modify the returned slice
	returnedAttachments := ef.Attachments()
	returnedAttachments[0] = "modified.pdf"
	
	// Original should be unchanged
	if ef.Attachments()[0] != "receipt.pdf" {
		t.Error("Attachments() did not return a copy - mutability issue")
	}
}

func TestExpenseForecastString(t *testing.T) {
	addedOn := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	plannedDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	attachments := []string{"receipt.pdf"}
	partner := createTestPartner(1)
	
	ef := NewExpenseForecast(
		1, partner, "Fertilizer", "Organic fertilizer for olives",
		150.50, plannedDate, ExpenseSubtypeA6, ExpenseScopeOliveSection,
		attachments, addedOn,
	)
	
	str := ef.String()
	
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
	
	if len(str) == 0 {
		t.Error("String() returned empty string")
	}
}

func TestExpenseForecastEmptyAttachments(t *testing.T) {
	addedOn := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	plannedDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	partner := createTestPartner(1)
	
	ef := NewExpenseForecast(
		1, partner, "Fertilizer", "Organic fertilizer for olives",
		150.50, plannedDate, ExpenseSubtypeA6, ExpenseScopeOliveSection,
		[]string{}, addedOn,
	)
	
	attachments := ef.Attachments()
	if len(attachments) != 0 {
		t.Errorf("Expected empty attachments, got %d items", len(attachments))
	}
}
