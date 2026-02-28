package dbo

import "github.com/pjover/espigol/internal/domain/model"

func ConvertExpenseForecastToModel(dbo ExpenseForecast, partner *model.Partner) *model.ExpenseForecast {
	return model.NewExpenseForecast(
		dbo.Id,
		*partner,
		dbo.Concept,
		dbo.Description,
		dbo.GrossAmount,
		dbo.PlannedDate,
		model.ExpenseSubtype(dbo.ExpenseSubtype),
		model.ExpenseScope(dbo.Scope),
		dbo.Attachments,
		dbo.AddedOn,
	)
}

func ConvertExpenseForecastToDbo(forecast *model.ExpenseForecast) ExpenseForecast {
	attachments := forecast.Attachments()
	if attachments == nil {
		attachments = []string{}
	}

	return ExpenseForecast{
		Id:             forecast.ID(),
		PartnerEmail:   forecast.Partner().Email(),
		Concept:        forecast.Concept(),
		Description:    forecast.Description(),
		GrossAmount:    forecast.GrossAmount(),
		PlannedDate:    forecast.PlannedDate(),
		ExpenseSubtype: forecast.ExpenseSubtype().String(),
		Scope:          forecast.Scope().String(),
		Attachments:    attachments,
		AddedOn:        forecast.AddedOn(),
	}
}
