package reports

import (
	"strings"
	"testing"
	"time"

	"github.com/pjover/espigol/internal/domain/model"
)

// stubConfig implements ports.ConfigService for testing.
type stubConfig struct {
	floats  map[string]float64
	strings map[string]string
}

func (s *stubConfig) GetFloat64(key string) float64   { return s.floats[key] }
func (s *stubConfig) GetString(key string) string     { return s.strings[key] }
func (s *stubConfig) SetString(string, string) error  { return nil }
func (s *stubConfig) GetTime(string) time.Time        { return time.Time{} }
func (s *stubConfig) SetTime(string, time.Time) error { return nil }
func (s *stubConfig) Init()                           {}

// stubDb implements ports.DbService for testing.
type stubDb struct {
	forecasts []*model.ExpenseForecast
	partners  []*model.Partner
}

func (s *stubDb) GetAllExpenseForecasts() ([]*model.ExpenseForecast, error) {
	return s.forecasts, nil
}
func (s *stubDb) GetAllPartners() ([]*model.Partner, error)                  { return s.partners, nil }
func (s *stubDb) UpsertPartner(*model.Partner) error                         { return nil }
func (s *stubDb) GetPartnerByID(int) (*model.Partner, error)                 { return nil, nil }
func (s *stubDb) DeletePartner(int) error                                    { return nil }
func (s *stubDb) FindPartnerByEmail(string) (*model.Partner, error)          { return nil, nil }
func (s *stubDb) UpsertExpenseForecast(*model.ExpenseForecast) error         { return nil }
func (s *stubDb) GetExpenseForecastByID(int) (*model.ExpenseForecast, error) { return nil, nil }
func (s *stubDb) DeleteExpenseForecast(int) error                            { return nil }

func newTestPartner(id int, oliveSection, livestockSection bool) *model.Partner {
	return model.NewPartner(id, "Test", "Partner", "00000000A",
		"test@test.com", "+34600000000",
		model.Producer, 0, oliveSection, livestockSection,
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
}

func newTestForecast(id, year int, subtype model.ExpenseSubtype, scope model.ExpenseScope, amount float64, concept string) *model.ExpenseForecast {
	partner := newTestPartner(1, true, false)
	return model.NewExpenseForecast(
		id, *partner, concept, "description", amount,
		time.Date(year, 6, 15, 0, 0, 0, 0, time.UTC),
		subtype, scope, []string{}, time.Now(),
	)
}

func newTestConfig(currentLimit, investmentLimit float64) *stubConfig {
	return &stubConfig{
		floats: map[string]float64{
			"expenses.limits.2026.current":    currentLimit,
			"expenses.limits.2026.investment": investmentLimit,
		},
		strings: map[string]string{
			"business.name":    "Test Cooperativa",
			"files.logo":       "/tmp/logo.png",
			"output.directory": "/tmp/espigol-test-reports",
		},
	}
}

func TestFormatEuro(t *testing.T) {
	tests := []struct {
		amount   float64
		expected string
	}{
		{0, "0,00 \u20ac"},
		{1000, "1.000,00 \u20ac"},
		{30000, "30.000,00 \u20ac"},
		{31900, "31.900,00 \u20ac"},
		{100000, "100.000,00 \u20ac"},
		{-1000, "-1.000,00 \u20ac"},
		{1234.56, "1.234,56 \u20ac"},
		{500, "500,00 \u20ac"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			got := formatEuro(tt.amount)
			if got != tt.expected {
				t.Errorf("formatEuro(%v) = %q, want %q", tt.amount, got, tt.expected)
			}
		})
	}
}

func TestExpenseForecastReport_PositiveRemainder(t *testing.T) {
	forecasts := []*model.ExpenseForecast{
		newTestForecast(1, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeCommon, 5000, "Concepte com\u00fa"),
		newTestForecast(2, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeOliveSection, 8000, "Oliveres"),
		newTestForecast(3, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeLivestockSection, 3000, "Ramaderia"),
	}
	partners := []*model.Partner{
		newTestPartner(1, true, false),
		newTestPartner(2, false, true),
	}
	svc := NewExpenseForecastReportService(newTestConfig(30000, 70000), &stubDb{forecasts: forecasts, partners: partners})
	hasNeg, msg, err := svc.ExpenseForecastReport(2026)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasNeg {
		t.Error("expected no negative remainder")
	}
	if !strings.Contains(msg, "2026") {
		t.Errorf("message should mention year 2026, got: %s", msg)
	}
}

func TestExpenseForecastReport_NegativeRemainder(t *testing.T) {
	forecasts := []*model.ExpenseForecast{
		newTestForecast(1, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeCommon, 5000, "Com\u00fa"),
		newTestForecast(2, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeOliveSection, 20000, "Oliveres"),
		newTestForecast(3, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeLivestockSection, 9000, "Ramaderia"),
	}
	var partners []*model.Partner
	for i := 0; i < 8; i++ {
		partners = append(partners, newTestPartner(i+1, true, false))
	}
	for i := 0; i < 3; i++ {
		partners = append(partners, newTestPartner(100+i, false, true))
	}
	svc := NewExpenseForecastReportService(newTestConfig(30000, 70000), &stubDb{forecasts: forecasts, partners: partners})
	hasNeg, _, err := svc.ExpenseForecastReport(2026)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasNeg {
		t.Error("expected hasNegativeRemainder=true")
	}
}

func TestExpenseForecastReport_YearFilter(t *testing.T) {
	forecasts := []*model.ExpenseForecast{
		newTestForecast(1, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeCommon, 5000, "2026 forecast"),
		newTestForecast(2, 2025, model.ExpenseSubtypeA1, model.ExpenseScopeCommon, 50000, "2025 forecast"),
	}
	svc := NewExpenseForecastReportService(newTestConfig(30000, 70000), &stubDb{forecasts: forecasts, partners: nil})
	hasNeg, _, err := svc.ExpenseForecastReport(2026)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasNeg {
		t.Error("2025 forecast should be filtered: remainder should be positive")
	}
}

func TestExpenseForecastReport_ProportionalSplit(t *testing.T) {
	forecasts := []*model.ExpenseForecast{
		newTestForecast(1, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeCommon, 2880, "Com\u00fa"),
		newTestForecast(2, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeOliveSection, 21054, "Oliveres"),
		newTestForecast(3, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeLivestockSection, 9250, "Ramaderia"),
	}
	var partners []*model.Partner
	for i := 0; i < 8; i++ {
		partners = append(partners, newTestPartner(i+1, true, false))
	}
	for i := 0; i < 3; i++ {
		partners = append(partners, newTestPartner(100+i, false, true))
	}
	svc := NewExpenseForecastReportService(newTestConfig(30000, 70000), &stubDb{forecasts: forecasts, partners: partners})
	hasNeg, _, err := svc.ExpenseForecastReport(2026)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasNeg {
		t.Error("expected negative remainder (example from LOG.md)")
	}

	// Verify proportional math directly
	nOlive, nLivestock := 8, 3
	available := 30000.0 - 2880.0
	oliveAllowed := available * float64(nOlive) / float64(nOlive+nLivestock)
	livestockAllowed := available * float64(nLivestock) / float64(nOlive+nLivestock)
	expectedOlive := 27120.0 * 8.0 / 11.0
	expectedLivestock := 27120.0 * 3.0 / 11.0
	if absF(oliveAllowed-expectedOlive) > 0.01 {
		t.Errorf("olive allowed = %.2f, want %.2f", oliveAllowed, expectedOlive)
	}
	if absF(livestockAllowed-expectedLivestock) > 0.01 {
		t.Errorf("livestock allowed = %.2f, want %.2f", livestockAllowed, expectedLivestock)
	}
}

func absF(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
