package ports

import "github.com/pjover/espigol/internal/domain/model"

// DbService defines the persistence layer interface for storing and retrieving entities.
type DbService interface {
	// UpsertPartner inserts or updates a partner in the database.
	UpsertPartner(partner *model.Partner) error

	// GetPartnerByID retrieves a partner by its numeric ID.
	GetPartnerByID(id int) (*model.Partner, error)

	// GetAllPartners retrieves all partners.
	GetAllPartners() ([]*model.Partner, error)

	// DeletePartner deletes a partner by its numeric ID.
	DeletePartner(id int) error

	// FindPartnerByEmail retrieves a partner by its email address.
	// Returns error if partner is not found.
	FindPartnerByEmail(email string) (*model.Partner, error)

	// UpsertExpenseForecast inserts or updates an expense forecast in the database.
	UpsertExpenseForecast(forecast *model.ExpenseForecast) error

	// GetExpenseForecastByID retrieves an expense forecast by its numeric ID.
	GetExpenseForecastByID(id int) (*model.ExpenseForecast, error)

	// GetAllExpenseForecasts retrieves all expense forecasts.
	GetAllExpenseForecasts() ([]*model.ExpenseForecast, error)

	// DeleteExpenseForecast deletes an expense forecast by its numeric ID.
	DeleteExpenseForecast(id int) error
}
