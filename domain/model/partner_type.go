package model

type PartnerType string

const (
	Producer     PartnerType = "Productor"
	Sponsor      PartnerType = "Patrocinador"
	Collaborator PartnerType = "Col·laborador"
)

func (pt PartnerType) String() string {
	return string(pt)
}

