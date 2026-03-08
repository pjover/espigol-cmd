package dbo

import (
	"strings"

	"github.com/pjover/espigol/internal/domain/model"
)

// expenseSubtypeCode extracts the short code (e.g. "a2") from a full ExpenseSubtype string.
func expenseSubtypeCode(subtype model.ExpenseSubtype) string {
	s := string(subtype)
	if len(s) > 2 && s[0] == '[' {
		if end := strings.Index(s, "]"); end > 1 {
			return s[1:end]
		}
	}
	return s
}

// expenseSubtypeFromCode maps a short code (e.g. "a2") back to the full ExpenseSubtype constant.
func expenseSubtypeFromCode(code string) model.ExpenseSubtype {
	map_ := map[string]model.ExpenseSubtype{
		"a1": model.ExpenseSubtypeA1,
		"a2": model.ExpenseSubtypeA2,
		"a3": model.ExpenseSubtypeA3,
		"a4": model.ExpenseSubtypeA4,
		"a5": model.ExpenseSubtypeA5,
		"a6": model.ExpenseSubtypeA6,
		"b1": model.ExpenseSubtypeB1,
		"b2": model.ExpenseSubtypeB2,
		"b3": model.ExpenseSubtypeB3,
		"b4": model.ExpenseSubtypeB4,
		"b5": model.ExpenseSubtypeB5,
		"c1": model.ExpenseSubtypeC1,
		"c2": model.ExpenseSubtypeC2,
	}
	if st, ok := map_[strings.ToLower(strings.TrimSpace(code))]; ok {
		return st
	}
	return model.ExpenseSubtype(code)
}

func ConvertExpenseForecastToModel(dbo ExpenseForecast, partner *model.Partner) *model.ExpenseForecast {
	return model.NewExpenseForecast(
		dbo.Id,
		*partner,
		dbo.Concept,
		dbo.Description,
		dbo.GrossAmount,
		dbo.PlannedDate,
		expenseSubtypeFromCode(dbo.ExpenseSubtype),
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
		PartnerId:      forecast.Partner().ID(),
		Concept:        forecast.Concept(),
		Description:    forecast.Description(),
		GrossAmount:    forecast.GrossAmount(),
		PlannedDate:    forecast.PlannedDate(),
		ExpenseSubtype: expenseSubtypeCode(forecast.ExpenseSubtype()),
		Scope:          forecast.Scope().String(),
		Attachments:    attachments,
		AddedOn:        forecast.AddedOn(),
	}
}
