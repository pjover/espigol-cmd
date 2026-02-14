package model

import (
	"fmt"
	"time"
)

type ExpenseForecast struct {
	id             int
	partner        Partner
	concept        string
	description    string
	grossAmount    float64
	plannedDate    time.Time
	expenseSubtype ExpenseSubtype
	scope          ExpenseScope
	attachments    []string
	addedOn        time.Time
}

func NewExpenseForecast(
	id int,
	partner Partner,
	concept string,
	description string,
	grossAmount float64,
	plannedDate time.Time,
	expenseSubtype ExpenseSubtype,
	scope ExpenseScope,
	attachments []string,
	addedOn time.Time,
) *ExpenseForecast {
	return &ExpenseForecast{
		id:             id,
		partner:        partner,
		concept:        concept,
		description:    description,
		grossAmount:    grossAmount,
		plannedDate:    plannedDate,
		expenseSubtype: expenseSubtype,
		scope:          scope,
		attachments:    attachments,
		addedOn:        addedOn,
	}
}

func (ef *ExpenseForecast) ID() int {
	return ef.id
}

func (ef *ExpenseForecast) Partner() Partner {
	return ef.partner
}

func (ef *ExpenseForecast) Concept() string {
	return ef.concept
}

func (ef *ExpenseForecast) Description() string {
	return ef.description
}

func (ef *ExpenseForecast) GrossAmount() float64 {
	return ef.grossAmount
}

func (ef *ExpenseForecast) PlannedDate() time.Time {
	return ef.plannedDate
}

func (ef *ExpenseForecast) ExpenseSubtype() ExpenseSubtype {
	return ef.expenseSubtype
}

func (ef *ExpenseForecast) Scope() ExpenseScope {
	return ef.scope
}

func (ef *ExpenseForecast) Attachments() []string {
	// Return a copy to maintain immutability
	attachmentsCopy := make([]string, len(ef.attachments))
	copy(attachmentsCopy, ef.attachments)
	return attachmentsCopy
}

func (ef *ExpenseForecast) AddedOn() time.Time {
	return ef.addedOn
}

func (ef *ExpenseForecast) String() string {
	return fmt.Sprintf("ExpenseForecast{id=%d, partnerId=%d, concept=%s, description=%s, grossAmount=%.2f, plannedDate=%s, expenseSubtype=%s, scope=%s, attachments=%v, addedOn=%s}",
		ef.id, ef.partner.ID(), ef.concept, ef.description, ef.grossAmount, ef.plannedDate.Format("2006-01-02"),
		ef.expenseSubtype, ef.scope, ef.attachments, ef.addedOn.Format("2006-01-02"))
}
