package dbo

import (
	"testing"
	"time"

	"github.com/pjover/espigol/internal/domain/model"
)

func TestConvertPartnerToDbo(t *testing.T) {
	addedOn := time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC)
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

	dboPartner := ConvertPartnerToDbo(partner)

	if dboPartner.Id != 1 || dboPartner.Name != "John" || dboPartner.Email != "john@example.com" {
		t.Error("Partner conversion to DBO failed")
	}
}

func TestConvertPartnerToModel(t *testing.T) {
	addedOn := time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC)
	dboPartner := Partner{
		Id:      5,
		Name:    "Jane",
		Surname: "Smith",
		Email:   "jane@example.com",
		AddedOn: addedOn,
	}

	partner := ConvertPartnerToModel(dboPartner)

	if partner.ID() != 5 || partner.Name() != "Jane" || partner.Email() != "jane@example.com" {
		t.Error("Partner conversion to model failed")
	}
}

func TestPartnerRoundTrip(t *testing.T) {
	addedOn := time.Date(2026, 1, 15, 14, 30, 0, 0, time.UTC)
	original := model.NewPartner(
		10,
		"Alice",
		"Johnson",
		"VAT999",
		"alice@example.com",
		"+34555666777",
		model.Producer,
		77,
		true,
		true,
		addedOn,
	)

	dboDbo := ConvertPartnerToDbo(original)
	restored := ConvertPartnerToModel(dboDbo)

	if restored.ID() != original.ID() || restored.Name() != original.Name() {
		t.Error("Round-trip conversion lost data")
	}
}
