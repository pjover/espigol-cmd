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

func TestBuildPartnersSubReport(t *testing.T) {
	forecasts := []*model.ExpenseForecast{
		newTestForecast(10, 2026, model.ExpenseSubtypeA1, model.ExpenseScopePartner, 1000, "Soci A1 a"),
		newTestForecast(11, 2026, model.ExpenseSubtypeA1, model.ExpenseScopePartner, 2000, "Soci A1 b"),
		newTestForecast(12, 2026, model.ExpenseSubtypeA6, model.ExpenseScopePartner, 500, "Soci A6"),
		newTestForecast(13, 2026, model.ExpenseSubtypeB1, model.ExpenseScopePartner, 8000, "Soci B1"),
		// Common scope — should be excluded
		newTestForecast(14, 2026, model.ExpenseSubtypeA1, model.ExpenseScopeCommon, 9999, "Comú"),
	}

	svc := NewExpenseForecastReportService(newTestConfig(30000, 70000), &stubDb{forecasts: forecasts})

	// Test current category with no excess (remainder=10000, total=3500)
	allocations := []partnerAllocation{
		{partnerID: 1, partnerName: "Test Partner", requested: 3500, allocated: 3500},
	}
	subs := svc.buildPartnersSubReport(model.ExpenseCategoryCurrent, forecasts, 10000, allocations, 6500)
	if len(subs) != 1 {
		t.Fatalf("expected 1 sub-report, got %d", len(subs))
	}
	table, ok := subs[0].(CustomTableSubReport)
	if !ok {
		t.Fatal("expected CustomTableSubReport")
	}
	if table.Title != "Despesa corrent (socis)" {
		t.Errorf("unexpected title: %s", table.Title)
	}
	// 2 subtype rows + 1 total row + 1 remanent final row = 4 rows
	if len(table.Rows) != 4 {
		t.Errorf("expected 4 rows, got %d", len(table.Rows))
	}
	// A1 total = 3000
	if table.Rows[0].Cells[1] != formatEuro(3000) {
		t.Errorf("A1 total = %s, want %s", table.Rows[0].Cells[1], formatEuro(3000))
	}
	// A6 total = 500
	if table.Rows[1].Cells[1] != formatEuro(500) {
		t.Errorf("A6 total = %s, want %s", table.Rows[1].Cells[1], formatEuro(500))
	}
	// Grand total = 3500
	if table.Rows[2].Cells[1] != formatEuro(3500) {
		t.Errorf("grand total = %s, want %s", table.Rows[2].Cells[1], formatEuro(3500))
	}
	if !table.Rows[2].Bold {
		t.Error("total row should be bold")
	}

	// Test investment category: should include B1 only
	allocsInv := []partnerAllocation{
		{partnerID: 1, partnerName: "Test Partner", requested: 8000, allocated: 8000},
	}
	subs2 := svc.buildPartnersSubReport(model.ExpenseCategoryInvestment, forecasts, 50000, allocsInv, 42000)
	if len(subs2) != 1 {
		t.Fatalf("expected 1 sub-report for investment, got %d", len(subs2))
	}
	table2 := subs2[0].(CustomTableSubReport)
	// 1 subtype row + 1 total row + 1 remanent final row = 3 rows
	if len(table2.Rows) != 3 {
		t.Errorf("expected 3 rows for investment, got %d", len(table2.Rows))
	}
	if table2.Rows[0].Cells[1] != formatEuro(8000) {
		t.Errorf("B1 total = %s, want %s", table2.Rows[0].Cells[1], formatEuro(8000))
	}
}

func TestDistributeRemainder_NoExcess(t *testing.T) {
	totals := map[int]float64{1: 1000, 2: 1000, 3: 1000}
	names := map[int]string{1: "A", 2: "B", 3: "C"}
	allocs, finalRem := distributeRemainder(5000, totals, names)
	if len(allocs) != 3 {
		t.Fatalf("expected 3 allocations, got %d", len(allocs))
	}
	for _, a := range allocs {
		if a.allocated != a.requested {
			t.Errorf("partner %s: allocated=%.2f, want %.2f", a.partnerName, a.allocated, a.requested)
		}
	}
	if absF(finalRem-2000) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want 2000", finalRem)
	}
}

func TestDistributeRemainder_UniformExcess(t *testing.T) {
	totals := map[int]float64{1: 4000, 2: 4000, 3: 4000}
	names := map[int]string{1: "A", 2: "B", 3: "C"}
	allocs, finalRem := distributeRemainder(9000, totals, names)
	for _, a := range allocs {
		if absF(a.allocated-3000) > 0.01 {
			t.Errorf("partner %s: allocated=%.2f, want 3000", a.partnerName, a.allocated)
		}
	}
	if absF(finalRem) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want ~0", finalRem)
	}
}

func TestDistributeRemainder_NonUniformExcess(t *testing.T) {
	// A=6000, B=2000, C=5000, remainder=9000
	// Mean = 3000. B(2000) <= mean → fixed. Budget left = 7000 for A,C.
	// Mean = 3500. Both A,C > 3500 → capped at 3500.
	// Total = 2000 + 3500 + 3500 = 9000
	totals := map[int]float64{1: 6000, 2: 2000, 3: 5000}
	names := map[int]string{1: "A", 2: "B", 3: "C"}
	allocs, finalRem := distributeRemainder(9000, totals, names)

	allocMap := map[string]float64{}
	for _, a := range allocs {
		allocMap[a.partnerName] = a.allocated
	}
	if absF(allocMap["B"]-2000) > 0.01 {
		t.Errorf("B allocated=%.2f, want 2000", allocMap["B"])
	}
	if absF(allocMap["A"]-3500) > 0.01 {
		t.Errorf("A allocated=%.2f, want 3500", allocMap["A"])
	}
	if absF(allocMap["C"]-3500) > 0.01 {
		t.Errorf("C allocated=%.2f, want 3500", allocMap["C"])
	}
	if absF(finalRem) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want ~0", finalRem)
	}
}

func TestDistributeRemainder_MultipleRounds(t *testing.T) {
	// A=10000, B=1000, C=500, D=8000, remainder=9000
	// Round 1: mean=2250. B(1000),C(500) fixed. Budget left=7500 for A,D.
	// Round 2: mean=3750. Both A,D > 3750 → capped.
	// Total = 1000 + 500 + 3750 + 3750 = 9000
	totals := map[int]float64{1: 10000, 2: 1000, 3: 500, 4: 8000}
	names := map[int]string{1: "A", 2: "B", 3: "C", 4: "D"}
	allocs, finalRem := distributeRemainder(9000, totals, names)

	var total float64
	allocMap := map[string]float64{}
	for _, a := range allocs {
		allocMap[a.partnerName] = a.allocated
		total += a.allocated
	}
	if absF(allocMap["B"]-1000) > 0.01 {
		t.Errorf("B=%.2f, want 1000", allocMap["B"])
	}
	if absF(allocMap["C"]-500) > 0.01 {
		t.Errorf("C=%.2f, want 500", allocMap["C"])
	}
	if absF(allocMap["A"]-3750) > 0.01 {
		t.Errorf("A=%.2f, want 3750", allocMap["A"])
	}
	if absF(allocMap["D"]-3750) > 0.01 {
		t.Errorf("D=%.2f, want 3750", allocMap["D"])
	}
	if absF(total-9000) > 0.01 {
		t.Errorf("total=%.2f, want 9000", total)
	}
	if absF(finalRem) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want ~0", finalRem)
	}
}

func TestDistributeRemainder_SinglePartner(t *testing.T) {
	// Partner requests more than remainder
	totals := map[int]float64{1: 5000}
	names := map[int]string{1: "A"}
	allocs, finalRem := distributeRemainder(3000, totals, names)
	if len(allocs) != 1 {
		t.Fatalf("expected 1 allocation, got %d", len(allocs))
	}
	if absF(allocs[0].allocated-3000) > 0.01 {
		t.Errorf("allocated=%.2f, want 3000", allocs[0].allocated)
	}
	if absF(finalRem) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want ~0", finalRem)
	}

	// Partner requests less than remainder
	allocs2, finalRem2 := distributeRemainder(8000, totals, names)
	if absF(allocs2[0].allocated-5000) > 0.01 {
		t.Errorf("allocated=%.2f, want 5000", allocs2[0].allocated)
	}
	if absF(finalRem2-3000) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want 3000", finalRem2)
	}
}

func TestDistributeRemainder_ZeroRemainder(t *testing.T) {
	totals := map[int]float64{1: 1000, 2: 2000}
	names := map[int]string{1: "A", 2: "B"}
	allocs, finalRem := distributeRemainder(0, totals, names)
	for _, a := range allocs {
		if absF(a.allocated) > 0.01 {
			t.Errorf("partner %s: allocated=%.2f, want 0", a.partnerName, a.allocated)
		}
	}
	if absF(finalRem) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want 0", finalRem)
	}
}

func TestDistributeRemainder_Empty(t *testing.T) {
	allocs, finalRem := distributeRemainder(5000, map[int]float64{}, map[int]string{})
	if len(allocs) != 0 {
		t.Errorf("expected 0 allocations, got %d", len(allocs))
	}
	if absF(finalRem-5000) > 0.01 {
		t.Errorf("finalRemainder=%.2f, want 5000", finalRem)
	}
}

func TestBuildPartnersSubReport_WithExcess(t *testing.T) {
	forecasts := []*model.ExpenseForecast{
		newTestForecast(10, 2026, model.ExpenseSubtypeA1, model.ExpenseScopePartner, 5000, "Soci A1"),
		newTestForecast(11, 2026, model.ExpenseSubtypeA6, model.ExpenseScopePartner, 3000, "Soci A6"),
	}
	svc := NewExpenseForecastReportService(newTestConfig(30000, 70000), &stubDb{forecasts: forecasts})

	allocations := []partnerAllocation{
		{partnerID: 1, partnerName: "Test Partner", requested: 8000, allocated: 5000},
	}
	// remainder=5000 but total=8000 → excess, should produce 2 sub-reports
	subs := svc.buildPartnersSubReport(model.ExpenseCategoryCurrent, forecasts, 5000, allocations, 0)
	if len(subs) != 2 {
		t.Fatalf("expected 2 sub-reports (table + adjustment), got %d", len(subs))
	}
	// Second sub-report should be the adjustment table
	adjTable, ok := subs[1].(CustomTableSubReport)
	if !ok {
		t.Fatal("expected CustomTableSubReport for adjustment")
	}
	if !strings.Contains(adjTable.Title, "Ajust") {
		t.Errorf("adjustment title should contain 'Ajust', got: %s", adjTable.Title)
	}
	// 1 partner row + 1 total row = 2 rows
	if len(adjTable.Rows) != 2 {
		t.Errorf("expected 2 rows in adjustment table, got %d", len(adjTable.Rows))
	}
}
