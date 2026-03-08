package dbo

import "time"

type Partner struct {
	Id               int       `bson:"_id"`
	Name             string    `bson:"name"`
	Surname          string    `bson:"surname"`
	VatCode          string    `bson:"vat_code"`
	Email            string    `bson:"email"`
	Mobile           string    `bson:"mobile"`
	PartnerType      string    `bson:"partner_type"`
	RiaNumber        int       `bson:"ria_number"`
	OliveSection     bool      `bson:"olive_section"`
	LivestockSection bool      `bson:"livestock_section"`
	AddedOn          time.Time `bson:"added_on"`
}
