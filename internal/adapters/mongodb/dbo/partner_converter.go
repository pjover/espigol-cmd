package dbo

import "github.com/pjover/espigol/internal/domain/model"

func ConvertPartnerToModel(dbo Partner) *model.Partner {
	return model.NewPartner(
		dbo.Id,
		dbo.Name,
		dbo.Surname,
		dbo.VatCode,
		dbo.Email,
		dbo.Mobile,
		model.PartnerType(dbo.PartnerType),
		dbo.RiaNumber,
		dbo.OliveSection,
		dbo.LivestockSection,
		dbo.AddedOn,
	)
}

func ConvertPartnerToDbo(partner *model.Partner) Partner {
	return Partner{
		Id:               partner.ID(),
		Name:             partner.Name(),
		Surname:          partner.Surname(),
		VatCode:          partner.VATCode(),
		Email:            partner.Email(),
		Mobile:           partner.Mobile(),
		PartnerType:      partner.PartnerType().String(),
		RiaNumber:        partner.RiaNumber(),
		OliveSection:     partner.OliveSection(),
		LivestockSection: partner.LivestockSection(),
		AddedOn:          partner.AddedOn(),
	}
}
