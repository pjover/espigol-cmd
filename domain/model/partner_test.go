package model

import (
	"testing"
	"time"
)

func TestNewPartner(t *testing.T) {
	addedOn := time.Date(2023, 4, 21, 0, 0, 0, 0, time.UTC)

	partner := NewPartner(1, "John", "Doe", "12345678A", "john@example.com", "+34600000000", Producer, 13937, true, true, addedOn)

	if partner.ID() != 1 {
		t.Errorf("Expected ID 1, got %d", partner.ID())
	}
	if partner.Name() != "John" {
		t.Errorf("Expected name 'John', got '%s'", partner.Name())
	}
	if partner.Surname() != "Doe" {
		t.Errorf("Expected surname 'Doe', got '%s'", partner.Surname())
	}
}

func TestPartnerGetters(t *testing.T) {
	addedOn := time.Date(2023, 4, 21, 0, 0, 0, 0, time.UTC)

	partner := NewPartner(1, "John", "Doe", "12345678A", "john@example.com", "+34600000000", Producer, 13937, true, true, addedOn)

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"VATCode", partner.VATCode(), "12345678A"},
		{"Email", partner.Email(), "john@example.com"},
		{"Mobile", partner.Mobile(), "+34600000000"},
		{"PartnerType", partner.PartnerType(), Producer},
		{"RiaNumber", partner.RiaNumber(), 13937},
		{"OliveSection", partner.OliveSection(), true},
		{"LivestockSection", partner.LivestockSection(), true},
		{"AddedOn", partner.AddedOn(), addedOn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.got)
			}
		})
	}
}

func TestPartnerScreenName(t *testing.T) {
	addedOn := time.Date(2023, 4, 21, 0, 0, 0, 0, time.UTC)
	partner := NewPartner(1, "John", "Doe", "12345678A", "john@example.com", "+34600000000", Producer, 13937, true, true, addedOn)

	screenName := partner.ScreenName()
	expected := "1 - John"

	if screenName != expected {
		t.Errorf("Expected screen name '%s', got '%s'", expected, screenName)
	}
}

func TestPartnerString(t *testing.T) {
	addedOn := time.Date(2023, 4, 21, 0, 0, 0, 0, time.UTC)
	partner := NewPartner(1, "John", "Doe", "12345678A", "john@example.com", "+34600000000", Producer, 13937, true, true, addedOn)

	str := partner.String()

	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}
