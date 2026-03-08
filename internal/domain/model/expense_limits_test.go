package model

import (
	"testing"
	"time"
)

// stubConfigService is a minimal ports.ConfigService for testing LimitsForYear.
type stubConfigService struct {
	values map[string]float64
}

func (s *stubConfigService) GetFloat64(key string) float64 {
	return s.values[key]
}

func (s *stubConfigService) GetString(_ string) string           { return "" }
func (s *stubConfigService) SetString(_ string, _ string) error  { return nil }
func (s *stubConfigService) GetTime(_ string) time.Time          { return time.Time{} }
func (s *stubConfigService) SetTime(_ string, _ time.Time) error { return nil }
func (s *stubConfigService) Init()                               {}

func TestLimitsForYear_found(t *testing.T) {
	cfg := &stubConfigService{
		values: map[string]float64{
			"expenses.limits.2026.current":    30000.0,
			"expenses.limits.2026.investment": 70000.0,
		},
	}

	limits, ok := LimitsForYear(2026, cfg)
	if !ok {
		t.Fatal("expected limits to be found for year 2026")
	}
	if limits.CurrentExpense != 30000.0 {
		t.Errorf("expected CurrentExpense 30000, got %f", limits.CurrentExpense)
	}
	if limits.InvestmentExpense != 70000.0 {
		t.Errorf("expected InvestmentExpense 70000, got %f", limits.InvestmentExpense)
	}
	if limits.Total != 100000.0 {
		t.Errorf("expected Total 100000, got %f", limits.Total)
	}
}

func TestLimitsForYear_notFound(t *testing.T) {
	cfg := &stubConfigService{values: map[string]float64{}}

	_, ok := LimitsForYear(2025, cfg)
	if ok {
		t.Error("expected no limits for year 2025, but got some")
	}
}
