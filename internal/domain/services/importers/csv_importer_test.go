package services

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pjover/espigol/internal/domain/model"
)

// MockDbService implements ports.DbService for testing.
type MockDbService struct {
	UpsertedPartners         []*model.Partner
	UpsertedExpenseForecasts []*model.ExpenseForecast
	UpsertPartnerErr         error
	UpsertExpenseForecastErr error
}

func (m *MockDbService) UpsertPartner(partner *model.Partner) error {
	if m.UpsertPartnerErr != nil {
		return m.UpsertPartnerErr
	}
	m.UpsertedPartners = append(m.UpsertedPartners, partner)
	return nil
}

func (m *MockDbService) FindPartnerByEmail(email string) (*model.Partner, error) {
	p := model.NewPartner(1, "Test", "Partner", "00000000A", email, "+34600000000",
		model.Producer, 0, false, false, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	return p, nil
}

func (m *MockDbService) UpsertExpenseForecast(forecast *model.ExpenseForecast) error {
	if m.UpsertExpenseForecastErr != nil {
		return m.UpsertExpenseForecastErr
	}
	m.UpsertedExpenseForecasts = append(m.UpsertedExpenseForecasts, forecast)
	return nil
}

func (m *MockDbService) GetPartnerByID(id int) (*model.Partner, error) {
	return nil, nil
}

func (m *MockDbService) GetAllPartners() ([]*model.Partner, error) {
	return nil, nil
}

func (m *MockDbService) DeletePartner(id int) error {
	return nil
}

func (m *MockDbService) GetExpenseForecastByID(id int) (*model.ExpenseForecast, error) {
	return nil, nil
}

func (m *MockDbService) GetAllExpenseForecasts() ([]*model.ExpenseForecast, error) {
	return nil, nil
}

func (m *MockDbService) DeleteExpenseForecast(id int) error {
	return nil
}

func TestImportPartners(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test_partners.csv")

	csvContent := `id,name,surname,vatCode,email,mobile,partnerType,riaNumber,oliveSection,livestockSection,addedOn
1,Partner,One,12345678A,partner1@example.com,+34600000001,Productor,1001,true,false,01/02/2020
2,Partner,Two,87654321B,partner2@example.com,+34600000002,Patrocinador,1002,false,true,15/03/2021
3,Partner,Three,11111111C,partner3@example.com,+34600000003,Col·laborador,1003,true,true,20/06/2022
`

	err := os.WriteFile(csvPath, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	mock := &MockDbService{}
	importer := NewCsvImporter(mock)
	_, err = importer.ImportPartners(csvPath)
	if err != nil {
		t.Errorf("ImportPartners failed: %v", err)
	}

	if len(mock.UpsertedPartners) != 3 {
		t.Errorf("Expected 3 partners upserted, got %d", len(mock.UpsertedPartners))
	}
}

func TestImportPartnersInvalidPath(t *testing.T) {
	mock := &MockDbService{}
	importer := NewCsvImporter(mock)
	_, err := importer.ImportPartners("/nonexistent/path/partners.csv")

	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestImportPartnersInvalidData(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "invalid_partners.csv")

	csvContent := `id,name,surname,vatCode,email,mobile,partnerType,riaNumber,oliveSection,livestockSection,addedOn
invalid,Partner,One,12345678A,partner1@example.com,+34600000001,Productor,1001,true,false,01/02/2020
`

	err := os.WriteFile(csvPath, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	mock := &MockDbService{}
	importer := NewCsvImporter(mock)
	_, err = importer.ImportPartners(csvPath)

	if err == nil {
		t.Error("Expected error for invalid data, got nil")
	}
}

func TestImportExpenseForecastsCommonSections(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "expense_forecasts_common.csv")

	csvContent := `id,Timestamp,Email address,Àmbit,Concepte,Descripció,Brut,Data,Pressuposts,Tipus de despesa
1,28/01/2026 08:18:00,anon1@example.com,Comú,Comunicacio,Projecte anual,11280,01/03/2026,,[a2] Activitats d'informació i promoció de productes agraris
2,21/01/2026 21:31:00,anon2@example.com,Secció d'oliva,Formacio,Curs tecnic,1200,01/03/2026,,[a3] Activitats d'informació i promoció de productes agraris
`

	err := os.WriteFile(csvPath, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	mock := &MockDbService{}
	importer := NewCsvImporter(mock)
	_, err = importer.ImportExpenseForecasts(csvPath)
	if err != nil {
		t.Errorf("ImportExpenseForecasts failed: %v", err)
	}

	if len(mock.UpsertedExpenseForecasts) != 2 {
		t.Errorf("Expected 2 expense forecasts upserted, got %d", len(mock.UpsertedExpenseForecasts))
	}
}

func TestImportExpenseForecastsPartners(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "expense_forecasts_partners.csv")

	csvContent := `id,Timestamp,Email address,Concepte,Descripció,Brut,Data,Pressuposts,Tipus de despesa
1,13/01/2026 19:51:00,anon3@example.com,Menjar animals,Compra pinso,4000,01/03/2026,,[a6] Despeses de fertilitzants, productes d'alimentació animal i ormejos
`

	err := os.WriteFile(csvPath, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	mock := &MockDbService{}
	importer := NewCsvImporter(mock)
	_, err = importer.ImportExpenseForecasts(csvPath)
	if err != nil {
		t.Errorf("ImportExpenseForecasts failed: %v", err)
	}

	if len(mock.UpsertedExpenseForecasts) != 1 {
		t.Errorf("Expected 1 expense forecast upserted, got %d", len(mock.UpsertedExpenseForecasts))
	}
}
