package dbo

import "time"

type ExpenseForecast struct {
	Id             int       `bson:"_id"`
	PartnerId      int       `bson:"partner_id"`
	Concept        string    `bson:"concept"`
	Description    string    `bson:"description"`
	GrossAmount    float64   `bson:"gross_amount"`
	PlannedDate    time.Time `bson:"planned_date"`
	ExpenseSubtype string    `bson:"expense_subtype"`
	Scope          string    `bson:"scope"`
	Attachments    []string  `bson:"attachments"`
	AddedOn        time.Time `bson:"added_on"`
}
