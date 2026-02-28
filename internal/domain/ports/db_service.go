package ports

import "github.com/pjover/espigol/internal/domain/model"

// DbService defines the persistence layer interface for storing and retrieving entities.
type DbService interface {
	// UpsertPartner inserts or updates a partner in the database.
	UpsertPartner(partner *model.Partner) error

	// FindPartnerByEmail retrieves a partner by its email address.
	// Returns error if partner is not found.
	FindPartnerByEmail(email string) (*model.Partner, error)

	// UpsertExpenseForecast inserts or updates an expense forecast in the database.
	UpsertExpenseForecast(forecast *model.ExpenseForecast) error
}
