package dbo

import (
	"testing"
	"time"

	"github.com/pjover/espigol/internal/domain/model"
)

func TestConvertExpenseForecastToDbo(t *testing.T) {
	addedOn := time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC)
	plannedDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	partner := model.NewPartner(
		1,
		"John",
		"Doe",
		"VAT123",
		"john@example.com",
		"+34123456789",
		model.Producer,
		42,
		true,
		false,
		addedOn,
	)

	forecast := model.NewExpenseForecast(
		101,
		*partner,
		"Tools",
		"Farming tools purchase",
		1500.50,
		plannedDate,
		model.ExpenseSubtypeA1,
		model.ExpenseScopeCommon,
		[]string{"invoice.pdf"},
		addedOn,
	)

	dboForecast := ConvertExpenseForecastToDbo(forecast)

	if dboForecast.Id != 101 || dboForecast.PartnerEmail != "john@example.com" {
		t.Error("ExpenseForecast conversion to DBO failed")
	}
}

func TestConvertExpenseForecastToModel(t *testing.T) {
	addedOn := time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC)

	partner := model.NewPartner(
		2,
		"Jane",
		"Smith",
		"VAT456",
		"jane@example.com",
		"+34987654321",
		model.Producer,
		99,
		false,
		true,
		addedOn,
	)

	dboForecast := ExpenseForecast{
		Id:           202,
		PartnerEmail: "jane@example.com",
		Concept:      "Seeds",
		AddedOn:      addedOn,
	}

	forecast := ConvertExpenseForecastToModel(dboForecast, partner)

	if forecast.ID() != 202 || forecast.Concept() != "Seeds" {
		t.Error("ExpenseForecast conversion to model failed")
	}
}

func TestExpenseForecastNilAttachments(t *testing.T) {
	addedOn := time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC)
	plannedDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	partner := model.NewPartner(
		3,
		"Bob",
		"Brown",
		"VAT789",
		"bob@example.com",
		"",
		model.Producer,
		0,
		false,
		false,
		addedOn,
	)

	forecast := model.NewExpenseForecast(
		303,
		*partner,
		"Labor",
		"Farm work",
		5000.00,
		plannedDate,
		model.ExpenseSubtypeB1,
		model.ExpenseScopeLivestockSection,
		nil,
		addedOn,
	)

	dboForecast := ConvertExpenseForecastToDbo(forecast)

	if dboForecast.Attachments == nil {
		t.Error("Attachments should be empty slice, not nil")
	}
	if len(dboForecast.Attachments) != 0 {
		t.Error("Attachments should be empty")
	}
}
