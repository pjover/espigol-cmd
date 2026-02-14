package model

import (
	"fmt"
	"time"
)

type Partner struct {
	id               int
	name             string
	surname          string
	vatCode          string
	email            string
	mobile           string
	partnerType      PartnerType
	riaNumber        int
	oliveSection     bool
	livestockSection bool
	addedOn          time.Time
}

func NewPartner(id int, name, surname, vatCode, email, mobile string, partnerType PartnerType, riaNumber int, oliveSection, livestockSection bool, addedOn time.Time) *Partner {
	return &Partner{
		id:               id,
		name:             name,
		surname:          surname,
		vatCode:          vatCode,
		email:            email,
		mobile:           mobile,
		partnerType:      partnerType,
		riaNumber:        riaNumber,
		oliveSection:     oliveSection,
		livestockSection: livestockSection,
		addedOn:          addedOn,
	}
}

func (p *Partner) ID() int {
	return p.id
}

func (p *Partner) Name() string {
	return p.name
}

func (p *Partner) Surname() string {
	return p.surname
}

func (p *Partner) VATCode() string {
	return p.vatCode
}

func (p *Partner) Email() string {
	return p.email
}

func (p *Partner) Mobile() string {
	return p.mobile
}

func (p *Partner) PartnerType() PartnerType {
	return p.partnerType
}

func (p *Partner) RiaNumber() int {
	return p.riaNumber
}

func (p *Partner) OliveSection() bool {
	return p.oliveSection
}

func (p *Partner) LivestockSection() bool {
	return p.livestockSection
}

func (p *Partner) AddedOn() time.Time {
	return p.addedOn
}

func (p *Partner) String() string {
	return fmt.Sprintf("Partner{id=%d, name=%s, surname=%s, vatCode=%s, email=%s, mobile=%s, type=%s, riaNumber=%d, oliveSection=%v, livestockSection=%v, addedOn=%s}",
		p.id, p.name, p.surname, p.vatCode, p.email, p.mobile, p.partnerType, p.riaNumber, p.oliveSection, p.livestockSection, p.addedOn.Format("2006-01-02"))
}
