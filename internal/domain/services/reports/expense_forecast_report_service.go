package reports

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/pjover/espigol/internal/domain/model"
	"github.com/pjover/espigol/internal/domain/ports"
)

// ExpenseForecastReportService generates the expense forecast PDF report.
type ExpenseForecastReportService struct {
	config ports.ConfigService
	db     ports.DbService
}

// NewExpenseForecastReportService creates a new ExpenseForecastReportService.
func NewExpenseForecastReportService(config ports.ConfigService, db ports.DbService) *ExpenseForecastReportService {
	return &ExpenseForecastReportService{config: config, db: db}
}

// ExpenseForecastReport generates the PDF report for the given year.
// Returns (hasNegativeRemainder, message, error).
// Implements ports.GenerateService.
func (s *ExpenseForecastReportService) ExpenseForecastReport(year int) (bool, string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Generant l'informe de previsions de despeses per a l'any %d...\n", year))

	forecasts, err := s.db.GetAllExpenseForecasts()
	if err != nil {
		return false, "", fmt.Errorf("error en carregar les previsions: %w", err)
	}

	partners, err := s.db.GetAllPartners()
	if err != nil {
		return false, "", fmt.Errorf("error en carregar els socis: %w", err)
	}

	// Filter forecasts for the requested year
	var yearForecasts []*model.ExpenseForecast
	for _, f := range forecasts {
		if f.Year() == year {
			yearForecasts = append(yearForecasts, f)
		}
	}

	hasNegativeRemainder := false
	var subReports []SubReport

	for _, cat := range []model.ExpenseCategory{model.ExpenseCategoryCurrent, model.ExpenseCategoryInvestment} {
		commonSub, commonTotal := s.buildCommonSubReport(cat, yearForecasts)
		subReports = append(subReports, commonSub)

		sectionsSub, oliveTotal, livestockTotal, sectionsTotal := s.buildSectionsSubReport(cat, yearForecasts)
		subReports = append(subReports, sectionsSub)

		limits, hasLimits := model.LimitsForYear(year, s.config)
		var limit float64
		if hasLimits {
			if cat == model.ExpenseCategoryCurrent {
				limit = limits.CurrentExpense
			} else {
				limit = limits.InvestmentExpense
			}
		}

		remainderSub, remainder := s.buildRemainderSubReport(cat, year, limit, commonTotal, sectionsTotal)
		subReports = append(subReports, remainderSub)

		if remainder < 0 {
			hasNegativeRemainder = true
			availableForSections := limit - commonTotal
			warnSub := s.buildWarningSubReport(cat, year, availableForSections, oliveTotal, livestockTotal, partners)
			subReports = append(subReports, warnSub)
		}

		// Compute partner totals and allocations for this category
		partnerTotals, partnerNames := s.computePartnerTotals(cat, yearForecasts)
		allocations, finalRemainder := distributeRemainder(remainder, partnerTotals, partnerNames)

		partnersSub := s.buildPartnersSubReport(cat, yearForecasts, remainder, allocations, finalRemainder)
		for _, sub := range partnersSub {
			subReports = append(subReports, sub)
		}

		// Add per-partner detail sections
		partnerDetails := s.buildPartnerDetailSubReports(cat, year, yearForecasts, allocations)
		for _, sub := range partnerDetails {
			subReports = append(subReports, sub)
		}

		// Page break between categories (not after the last one)
		if cat != model.ExpenseCategoryInvestment {
			subReports = append(subReports, NewPageBreak())
		}
	}

	def := ReportDefinition{
		PageOrientation: consts.Portrait,
		Title:           fmt.Sprintf("Previsions de despeses %d", year),
		Footer:          time.Now().Format("02/01/2006"),
		SubReports:      subReports,
	}

	outputDir := expandOutputDir(s.config.GetString("output.directory"))
	filePath := path.Join(outputDir, fmt.Sprintf("Despeses %d.pdf", year))

	renderer := NewReportPdf(s.config)
	if err := renderer.SaveToFile(def, filePath); err != nil {
		return false, "", fmt.Errorf("error en guardar el PDF: %w", err)
	}

	buf.WriteString(fmt.Sprintf("Informe generat a '%s'", filePath))
	return hasNegativeRemainder, buf.String(), nil
}

// buildCommonSubReport builds the common-scope table for a category.
func (s *ExpenseForecastReportService) buildCommonSubReport(
	cat model.ExpenseCategory,
	forecasts []*model.ExpenseForecast,
) (SubReport, float64) {
	title := fmt.Sprintf("%s (comú)", cat.String())
	var rows []RowDef
	var total float64

	var filtered []*model.ExpenseForecast
	for _, f := range forecasts {
		if f.ExpenseCategory() == cat && f.Scope() == model.ExpenseScopeCommon {
			filtered = append(filtered, f)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Concept() < filtered[j].Concept()
	})

	for i, f := range filtered {
		scopeLabel := ""
		if i == 0 {
			scopeLabel = f.Scope().String()
			rows = append(rows, RowDef{
				Cells: []string{scopeLabel, "", ""},
				Bold:  true,
			})
		}
		rows = append(rows, RowDef{
			Cells: []string{"", f.Concept(), formatEuro(f.GrossAmount())},
		})
		total += f.GrossAmount()
	}

	// Total row
	rows = append(rows, RowDef{
		Cells: []string{"Total comú", "", formatEuro(total)},
		Bold:  true,
	})

	return CustomTableSubReport{
		Title:   title,
		Widths:  []uint{3, 7, 2},
		Headers: []string{"Àmbit", "Concepte", "Brut"},
		Rows:    rows,
	}, total
}

// buildSectionsSubReport builds the sections-scope table for a category.
// Returns the sub-report, olive total, livestock total, and grand sections total.
func (s *ExpenseForecastReportService) buildSectionsSubReport(
	cat model.ExpenseCategory,
	forecasts []*model.ExpenseForecast,
) (SubReport, float64, float64, float64) {
	title := fmt.Sprintf("%s (seccions)", cat.String())
	var rows []RowDef

	sectionScopes := []model.ExpenseScope{model.ExpenseScopeOliveSection, model.ExpenseScopeLivestockSection}
	sectionTotals := map[model.ExpenseScope]float64{}
	var grandTotal float64

	for _, scope := range sectionScopes {
		var filtered []*model.ExpenseForecast
		for _, f := range forecasts {
			if f.ExpenseCategory() == cat && f.Scope() == scope {
				filtered = append(filtered, f)
			}
		}
		if len(filtered) == 0 {
			continue
		}
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Concept() < filtered[j].Concept()
		})

		var sectionTotal float64
		for i, f := range filtered {
			scopeLabel := ""
			if i == 0 {
				scopeLabel = f.Scope().String()
				rows = append(rows, RowDef{
					Cells: []string{scopeLabel, "", ""},
					Bold:  true,
				})
			}
			rows = append(rows, RowDef{
				Cells: []string{"", f.Concept(), formatEuro(f.GrossAmount())},
			})
			sectionTotal += f.GrossAmount()
		}
		rows = append(rows, RowDef{
			Cells: []string{fmt.Sprintf("Total (%s)", scope.String()), "", formatEuro(sectionTotal)},
			Bold:  true,
		})
		sectionTotals[scope] = sectionTotal
		grandTotal += sectionTotal
	}

	// Grand total row
	rows = append(rows, RowDef{
		Cells: []string{"Total seccions", "", formatEuro(grandTotal)},
		Bold:  true,
	})

	return CustomTableSubReport{
		Title:   title,
		Widths:  []uint{3, 7, 2},
		Headers: []string{"Àmbit", "Concepte", "Brut (Suma)"},
		Rows:    rows,
	}, sectionTotals[model.ExpenseScopeOliveSection], sectionTotals[model.ExpenseScopeLivestockSection], grandTotal
}

// buildRemainderSubReport builds the remainder summary table.
func (s *ExpenseForecastReportService) buildRemainderSubReport(
	cat model.ExpenseCategory,
	year int,
	limit float64,
	commonTotal float64,
	sectionsTotal float64,
) (SubReport, float64) {
	title := fmt.Sprintf("Remanent de %s", strings.ToLower(cat.String()))
	availableForSections := limit - commonTotal
	categoryTotal := commonTotal + sectionsTotal
	remainder := limit - categoryTotal

	rows := []RowDef{
		{Cells: []string{"", ""}, Bold: false},
		{Cells: []string{fmt.Sprintf("Disponible any %d", year), formatEuro(limit)}},
		{Cells: []string{"Total comú", formatEuro(commonTotal)}},
		{Cells: []string{"Disponible per seccions", formatEuro(availableForSections)}},
		{Cells: []string{"Total seccions", formatEuro(sectionsTotal)}},
		{Cells: []string{fmt.Sprintf("Total %s", strings.ToLower(cat.String())), formatEuro(categoryTotal)}},
		{Cells: []string{fmt.Sprintf("Remanent de %s", strings.ToLower(cat.String())), formatEuro(remainder)}, Bold: true},
	}

	return CustomTableSubReport{
		Title:   title,
		Widths:  []uint{8, 4},
		Headers: []string{"", ""},
		Rows:    rows,
	}, remainder
}

// buildWarningSubReport builds the proportional-adjustment warning table when remainder < 0.
func (s *ExpenseForecastReportService) buildWarningSubReport(
	cat model.ExpenseCategory,
	year int,
	availableForSections float64,
	oliveRequested float64,
	livestockRequested float64,
	partners []*model.Partner,
) SubReport {
	_ = year

	nOlive := 0
	nLivestock := 0
	for _, p := range partners {
		if p.PartnerType() == model.Producer {
			if p.OliveSection() {
				nOlive++
			}
			if p.LivestockSection() {
				nLivestock++
			}
		}
	}
	nTotal := nOlive + nLivestock

	var oliveAllowed, livestockAllowed float64
	if nTotal > 0 {
		oliveAllowed = availableForSections * float64(nOlive) / float64(nTotal)
		livestockAllowed = availableForSections * float64(nLivestock) / float64(nTotal)
	}

	title := fmt.Sprintf("⚠ AVÍS: Ajust necessari per %s", cat.String())
	rows := []RowDef{
		{Cells: []string{"", "", "", ""}, Bold: false},
		{
			Cells: []string{
				model.ExpenseScopeOliveSection.String(),
				fmt.Sprintf("%d", nOlive),
				formatEuro(oliveAllowed),
				formatEuro(oliveRequested - oliveAllowed),
			},
		},
		{
			Cells: []string{
				model.ExpenseScopeLivestockSection.String(),
				fmt.Sprintf("%d", nLivestock),
				formatEuro(livestockAllowed),
				formatEuro(livestockRequested - livestockAllowed),
			},
		},
	}

	return CustomTableSubReport{
		Title:   title,
		Widths:  []uint{4, 2, 3, 3},
		Headers: []string{"Secció", "Socis productors", "Disponible", "Ajust"},
		Rows:    rows,
	}
}

// computePartnerTotals aggregates partner-scope forecasts by partner ID for a category.
func (s *ExpenseForecastReportService) computePartnerTotals(
	cat model.ExpenseCategory,
	forecasts []*model.ExpenseForecast,
) (map[int]float64, map[int]string) {
	totals := map[int]float64{}
	names := map[int]string{}
	for _, f := range forecasts {
		if f.ExpenseCategory() == cat && f.Scope() == model.ExpenseScopePartner {
			pID := f.Partner().ID()
			totals[pID] += f.GrossAmount()
			if _, ok := names[pID]; !ok {
				names[pID] = f.Partner().Name() + " " + f.Partner().Surname()
			}
		}
	}
	return totals, names
}

// partnerAllocation holds the redistribution result for a single partner.
type partnerAllocation struct {
	partnerID   int
	partnerName string
	requested   float64
	allocated   float64
}

// distributeRemainder distributes the available remainder among partners iteratively.
// Partners requesting less than or equal to the fair share keep their amount;
// those above are capped, and surplus is redistributed until convergence.
// Returns allocations sorted by partner name and the final unallocated remainder.
func distributeRemainder(
	remainder float64,
	partnerTotals map[int]float64,
	partnerNames map[int]string,
) ([]partnerAllocation, float64) {
	if len(partnerTotals) == 0 {
		return nil, remainder
	}

	// Initialize allocations with requested amounts
	type entry struct {
		id    int
		name  string
		req   float64
		alloc float64
		fixed bool
	}
	entries := make([]entry, 0, len(partnerTotals))
	var totalRequested float64
	for id, amount := range partnerTotals {
		entries = append(entries, entry{id: id, name: partnerNames[id], req: amount, alloc: amount})
		totalRequested += amount
	}

	// Case 1: no excess
	if totalRequested <= remainder {
		result := make([]partnerAllocation, len(entries))
		for i, e := range entries {
			result[i] = partnerAllocation{partnerID: e.id, partnerName: e.name, requested: e.req, allocated: e.req}
		}
		sort.Slice(result, func(i, j int) bool { return result[i].partnerName < result[j].partnerName })
		return result, remainder - totalRequested
	}

	// Case 2: iterative redistribution
	budgetLeft := remainder
	for iter := 0; iter < 100; iter++ {
		// Count unfixed partners and their total
		nUnfixed := 0
		var unfixedTotal float64
		for i := range entries {
			if !entries[i].fixed {
				nUnfixed++
				unfixedTotal += entries[i].alloc
			}
		}
		if nUnfixed == 0 {
			break
		}

		mean := budgetLeft / float64(nUnfixed)

		// Fix partners at or below the mean
		newlyFixed := false
		for i := range entries {
			if !entries[i].fixed && entries[i].alloc <= mean {
				entries[i].fixed = true
				budgetLeft -= entries[i].alloc
				newlyFixed = true
			}
		}

		if !newlyFixed {
			// All unfixed are above the mean → cap them equally
			for i := range entries {
				if !entries[i].fixed {
					entries[i].alloc = mean
					entries[i].fixed = true
				}
			}
			break
		}

		// Check convergence
		var totalAllocated float64
		for _, e := range entries {
			totalAllocated += e.alloc
		}
		if math.Abs(remainder-totalAllocated) < 0.01 {
			break
		}
	}

	// Build result
	var totalAllocated float64
	result := make([]partnerAllocation, len(entries))
	for i, e := range entries {
		result[i] = partnerAllocation{partnerID: e.id, partnerName: e.name, requested: e.req, allocated: e.alloc}
		totalAllocated += e.alloc
	}
	sort.Slice(result, func(i, j int) bool { return result[i].partnerName < result[j].partnerName })
	return result, remainder - totalAllocated
}

// buildPartnersSubReport builds the partner-scope expenses table for a category,
// grouped by ExpenseSubtype, plus remainder or adjustment info.
func (s *ExpenseForecastReportService) buildPartnersSubReport(
	cat model.ExpenseCategory,
	forecasts []*model.ExpenseForecast,
	remainder float64,
	allocations []partnerAllocation,
	finalRemainder float64,
) []SubReport {
	title := fmt.Sprintf("%s (socis)", cat.String())

	var filtered []*model.ExpenseForecast
	for _, f := range forecasts {
		if f.ExpenseCategory() == cat && f.Scope() == model.ExpenseScopePartner {
			filtered = append(filtered, f)
		}
	}

	// Group by ExpenseSubtype
	type subtypeGroup struct {
		subtype model.ExpenseSubtype
		total   float64
	}
	seen := map[model.ExpenseSubtype]int{}
	var groups []subtypeGroup

	sort.Slice(filtered, func(i, j int) bool {
		return string(filtered[i].ExpenseSubtype()) < string(filtered[j].ExpenseSubtype())
	})

	for _, f := range filtered {
		st := f.ExpenseSubtype()
		if idx, ok := seen[st]; ok {
			groups[idx].total += f.GrossAmount()
		} else {
			seen[st] = len(groups)
			groups = append(groups, subtypeGroup{subtype: st, total: f.GrossAmount()})
		}
	}

	var rows []RowDef
	var grandTotal float64
	for _, g := range groups {
		rows = append(rows, RowDef{
			Cells: []string{g.subtype.String(), formatEuro(g.total)},
		})
		grandTotal += g.total
	}
	rows = append(rows, RowDef{
		Cells: []string{"Total socis", formatEuro(grandTotal)},
		Bold:  true,
	})

	// Add remainder info
	if grandTotal <= remainder {
		rows = append(rows, RowDef{
			Cells: []string{"Remanent final", formatEuro(finalRemainder)},
			Bold:  true,
		})
	}

	result := []SubReport{
		CustomTableSubReport{
			Title:   title,
			Widths:  []uint{8, 4},
			Headers: []string{"Subtipus de despesa", "Brut (Suma)"},
			Rows:    rows,
		},
	}

	// If excess: add adjustment table per partner
	if grandTotal > remainder {
		result = append(result, s.buildPartnerAdjustmentSubReport(cat, allocations))
	}

	return result
}

// buildPartnerAdjustmentSubReport builds the per-partner adjustment table.
func (s *ExpenseForecastReportService) buildPartnerAdjustmentSubReport(
	cat model.ExpenseCategory,
	allocations []partnerAllocation,
) SubReport {
	title := fmt.Sprintf("Ajust de despeses per soci (%s)", cat.String())
	var rows []RowDef
	var totalRequested, totalAllocated float64

	redColor := &color.Color{Red: 200, Green: 0, Blue: 0}
	for _, a := range allocations {
		row := RowDef{
			Cells: []string{a.partnerName, formatEuro(a.requested), formatEuro(a.allocated)},
		}
		if a.allocated < a.requested {
			row.Color = redColor
		}
		rows = append(rows, row)
		totalRequested += a.requested
		totalAllocated += a.allocated
	}
	rows = append(rows, RowDef{
		Cells: []string{"Total", formatEuro(totalRequested), formatEuro(totalAllocated)},
		Bold:  true,
	})

	return CustomTableSubReport{
		Title:   title,
		Widths:  []uint{5, 4, 3},
		Headers: []string{"Soci", "Sol·licitat", "Assignat"},
		Rows:    rows,
	}
}

// partnerSectionName returns the section name(s) for a partner.
func partnerSectionName(p *model.Partner) string {
	var sections []string
	if p.OliveSection() {
		sections = append(sections, model.ExpenseScopeOliveSection.String())
	}
	if p.LivestockSection() {
		sections = append(sections, model.ExpenseScopeLivestockSection.String())
	}
	if len(sections) == 0 {
		return ""
	}
	return strings.Join(sections, ", ")
}

// forecastCode generates the CP code: "CP" + last 2 year digits + ID with 3 digits.
func forecastCode(year, id int) string {
	return fmt.Sprintf("CP%02d%03d", year%100, id)
}

// buildPartnerDetailSubReports builds per-partner detail tables for a category.
// Each partner with forecasts gets a section with their individual expense lines.
func (s *ExpenseForecastReportService) buildPartnerDetailSubReports(
	cat model.ExpenseCategory,
	year int,
	forecasts []*model.ExpenseForecast,
	allocations []partnerAllocation,
) []SubReport {
	// Build allocation lookup
	allocMap := map[int]*partnerAllocation{}
	for i := range allocations {
		allocMap[allocations[i].partnerID] = &allocations[i]
	}

	// Group forecasts by partner
	type partnerForecasts struct {
		partner   *model.Partner
		forecasts []*model.ExpenseForecast
	}
	partnerMap := map[int]*partnerForecasts{}
	var partnerOrder []int

	for _, f := range forecasts {
		if f.ExpenseCategory() != cat || f.Scope() != model.ExpenseScopePartner {
			continue
		}
		pID := f.Partner().ID()
		if _, ok := partnerMap[pID]; !ok {
			partnerMap[pID] = &partnerForecasts{partner: f.Partner()}
			partnerOrder = append(partnerOrder, pID)
		}
		partnerMap[pID].forecasts = append(partnerMap[pID].forecasts, f)
	}

	if len(partnerOrder) == 0 {
		return nil
	}

	// Sort partners by name
	sort.Slice(partnerOrder, func(i, j int) bool {
		pi := partnerMap[partnerOrder[i]].partner
		pj := partnerMap[partnerOrder[j]].partner
		ni := pi.Name() + " " + pi.Surname()
		nj := pj.Name() + " " + pj.Surname()
		return ni < nj
	})

	var subReports []SubReport
	redColor := &color.Color{Red: 200, Green: 0, Blue: 0}

	for _, pID := range partnerOrder {
		pf := partnerMap[pID]
		p := pf.partner
		alloc := allocMap[pID]
		isCapped := alloc != nil && alloc.allocated < alloc.requested

		sectionName := partnerSectionName(p)
		title := p.Name() + " " + p.Surname()
		if sectionName != "" {
			title += ", " + sectionName
		}

		// Sort forecasts by concept
		sort.Slice(pf.forecasts, func(i, j int) bool {
			return pf.forecasts[i].Concept() < pf.forecasts[j].Concept()
		})

		var rows []RowDef
		var total float64
		for _, f := range pf.forecasts {
			row := RowDef{
				Cells: []string{
					forecastCode(year, f.ID()),
					f.Concept(),
					formatEuro(f.GrossAmount()),
					f.ExpenseSubtype().Type().String(),
				},
			}
			if isCapped {
				row.Strikethrough = map[int]bool{2: true}
				row.Color = redColor
			}
			rows = append(rows, row)
			total += f.GrossAmount()
		}

		// Total row
		rows = append(rows, RowDef{
			Cells: []string{"", "Total", formatEuro(total), ""},
			Bold:  true,
		})

		// If capped, add the max authorized amount
		if isCapped {
			rows = append(rows, RowDef{
				Cells: []string{"", "Import màxim autoritzat", formatEuro(alloc.allocated), ""},
				Bold:  true,
				Color: redColor,
			})
		}

		subReports = append(subReports, CustomTableSubReport{
			Title:   title,
			Widths:  []uint{1, 4, 2, 5},
			Headers: []string{"CP", "Concepte", "Brut", "Tipus de despesa"},
			Rows:    rows,
		})
	}

	return subReports
}

// formatEuro formats a float as European currency: "31.900,00 €".
func formatEuro(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	decPart := parts[1]

	negative := strings.HasPrefix(intPart, "-")
	if negative {
		intPart = intPart[1:]
	}

	n := len(intPart)
	var result []byte
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, intPart[i])
	}

	formatted := string(result) + "," + decPart + " €"
	if negative {
		return "-" + formatted
	}
	return formatted
}

// expandOutputDir expands a leading ~ in the output directory.
func expandOutputDir(dir string) string {
	if !strings.HasPrefix(dir, "~/") {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return dir
	}
	return path.Join(home, dir[2:])
}
