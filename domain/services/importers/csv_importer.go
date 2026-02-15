package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pjover/espigol/domain/model"
	"github.com/pjover/espigol/domain/ports"
)

type CSVImporter struct{}

func NewCSVImporter() ports.Importer { return &CSVImporter{} }

// Reads the CSV file at path and stores the Partners.
func (c *CSVImporter) ImportPartners(path string) error {
	p, err := expandPath(path)
	if err != nil {
		return err
	}

	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	columnIndexes := make(map[string]int)
	for i, h := range header {
		columnIndexes[h] = i
	}

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read row: %w", err)
		}

		partner, err := c.parsePartner(rec, columnIndexes)
		if err != nil {
			return fmt.Errorf("parse partner: %w", err)
		}

		fmt.Println(partner)
	}

	return nil
}

// Reads the CSV file at the provided path and outputs ExpenseForecasts.
func (c *CSVImporter) ImportExpenseForecasts(path string) error {
	return c.importExpenseForecastsFromFile(path)
}

func (c *CSVImporter) importExpenseForecastsFromFile(path string) error {
	p, err := expandPath(path)
	if err != nil {
		return err
	}

	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	columnIndexes := make(map[string]int)
	for i, h := range header {
		columnIndexes[strings.TrimSpace(h)] = i
	}

	var scopeOverride *model.ExpenseScope
	if _, ok := columnIndexes["Àmbit"]; !ok {
		scope := model.ExpenseScopePartner
		scopeOverride = &scope
	}

	rowID := 1
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read row: %w", err)
		}

		forecast, err := c.parseExpenseForecast(rec, columnIndexes, scopeOverride, rowID)
		if err != nil {
			return fmt.Errorf("parse expense forecast: %w", err)
		}

		fmt.Println(forecast)
		rowID++
	}

	return nil
}

func (c *CSVImporter) parsePartner(rec []string, columnIndexes map[string]int) (*model.Partner, error) {
	getField := func(name string) string {
		if idx, ok := columnIndexes[name]; ok && idx < len(rec) {
			return strings.TrimSpace(rec[idx])
		}
		return ""
	}

	id, err := strconv.Atoi(getField("id"))
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	name := getField("name")
	surname := getField("surname")
	vatCode := getField("vatCode")
	email := getField("email")
	mobile := getField("mobile")
	partnerTypeStr := getField("partnerType")

	riaNumber, err := strconv.Atoi(getField("riaNumber"))
	if err != nil {
		return nil, fmt.Errorf("invalid riaNumber: %w", err)
	}

	oliveSection := parseBoolean(getField("oliveSection"))
	livestockSection := parseBoolean(getField("livestockSection"))

	addedOnStr := getField("addedOn")
	addedOn, err := time.Parse("02/01/2006", addedOnStr)
	if err != nil {
		return nil, fmt.Errorf("invalid addedOn date: %w", err)
	}

	return model.NewPartner(id, name, surname, vatCode, email, mobile, model.PartnerType(partnerTypeStr), riaNumber, oliveSection, livestockSection, addedOn), nil
}

func (c *CSVImporter) parseExpenseForecast(rec []string, columnIndexes map[string]int, scopeOverride *model.ExpenseScope, rowID int) (*model.ExpenseForecast, error) {
	getField := func(name string) string {
		if idx, ok := columnIndexes[name]; ok && idx < len(rec) {
			return strings.TrimSpace(rec[idx])
		}
		return ""
	}

	shiftAfter := -1
	if idx, ok := columnIndexes["Email address"]; ok {
		shiftAfter = idx
	}
	getFieldShifted := func(name string) string {
		idx, ok := columnIndexes[name]
		if !ok {
			return ""
		}
		if shiftAfter >= 0 && idx > shiftAfter {
			idx++
		}
		if idx < len(rec) {
			return strings.TrimSpace(rec[idx])
		}
		return ""
	}

	useShifted := false

	addedOn, err := time.Parse("02/01/2006 15:04:05", getField("Timestamp"))
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	plannedDate, err := time.Parse("02/01/2006", getField("Data"))
	if err != nil {
		if len(rec) == len(columnIndexes)+1 {
			plannedDate, err = time.Parse("02/01/2006", getFieldShifted("Data"))
			if err == nil {
				useShifted = true
			}
		}
		if err != nil {
			return nil, fmt.Errorf("invalid planned date: %w", err)
		}
	}

	grossAmountStr := strings.ReplaceAll(getField("Brut"), ",", ".")
	grossAmount, err := strconv.ParseFloat(grossAmountStr, 64)
	if err != nil {
		if len(rec) == len(columnIndexes)+1 {
			grossAmountStr = strings.ReplaceAll(getFieldShifted("Brut"), ",", ".")
			grossAmount, err = strconv.ParseFloat(grossAmountStr, 64)
			if err == nil {
				useShifted = true
			}
		}
		if err != nil {
			return nil, fmt.Errorf("invalid gross amount: %w", err)
		}
	}

	getValue := getField
	if useShifted {
		getValue = getFieldShifted
	}

	var scope model.ExpenseScope
	if scopeOverride != nil {
		scope = *scopeOverride
	} else {
		parsedScope, err := parseExpenseScope(getValue("Àmbit"))
		if err != nil {
			return nil, err
		}
		scope = parsedScope
	}

	expenseSubtype, err := parseExpenseSubtype(getValue("Tipus de despesa"))
	if err != nil {
		return nil, err
	}

	partner := partnerFromEmail(getValue("Email address"), rowID, addedOn)

	return model.NewExpenseForecast(
		rowID,
		partner,
		getValue("Concepte"),
		getValue("Descripció"),
		grossAmount,
		plannedDate,
		expenseSubtype,
		scope,
		nil,
		addedOn,
	), nil
}

func partnerFromEmail(email string, id int, addedOn time.Time) model.Partner {
	return *model.NewPartner(
		id,
		"Unknown",
		"",
		"",
		email,
		"",
		model.Producer,
		0,
		false,
		false,
		addedOn,
	)
}

func parseExpenseScope(value string) (model.ExpenseScope, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))

	if strings.Contains(normalized, "oliva") {
		return model.ExpenseScopeOliveSection, nil
	}
	if strings.Contains(normalized, "ramaderia") {
		return model.ExpenseScopeLivestockSection, nil
	}
	if strings.Contains(normalized, "com") {
		return model.ExpenseScopeCommon, nil
	}

	return "", fmt.Errorf("unknown scope: %s", value)
}

func parseExpenseSubtype(value string) (model.ExpenseSubtype, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("missing expense subtype")
	}

	if strings.HasPrefix(trimmed, "[") {
		if end := strings.Index(trimmed, "]"); end > 1 {
			code := trimmed[1:end]
			switch strings.ToLower(code) {
			case "a1":
				return model.ExpenseSubtypeA1, nil
			case "a2":
				return model.ExpenseSubtypeA2, nil
			case "a3":
				return model.ExpenseSubtypeA3, nil
			case "a4":
				return model.ExpenseSubtypeA4, nil
			case "a5":
				return model.ExpenseSubtypeA5, nil
			case "a6":
				return model.ExpenseSubtypeA6, nil
			case "b1":
				return model.ExpenseSubtypeB1, nil
			case "b2":
				return model.ExpenseSubtypeB2, nil
			case "b3":
				return model.ExpenseSubtypeB3, nil
			case "b4":
				return model.ExpenseSubtypeB4, nil
			case "b5":
				return model.ExpenseSubtypeB5, nil
			case "c1":
				return model.ExpenseSubtypeC1, nil
			case "c2":
				return model.ExpenseSubtypeC2, nil
			}
		}
	}

	switch trimmed {
	case string(model.ExpenseSubtypeA1):
		return model.ExpenseSubtypeA1, nil
	case string(model.ExpenseSubtypeA2):
		return model.ExpenseSubtypeA2, nil
	case string(model.ExpenseSubtypeA3):
		return model.ExpenseSubtypeA3, nil
	case string(model.ExpenseSubtypeA4):
		return model.ExpenseSubtypeA4, nil
	case string(model.ExpenseSubtypeA5):
		return model.ExpenseSubtypeA5, nil
	case string(model.ExpenseSubtypeA6):
		return model.ExpenseSubtypeA6, nil
	case string(model.ExpenseSubtypeB1):
		return model.ExpenseSubtypeB1, nil
	case string(model.ExpenseSubtypeB2):
		return model.ExpenseSubtypeB2, nil
	case string(model.ExpenseSubtypeB3):
		return model.ExpenseSubtypeB3, nil
	case string(model.ExpenseSubtypeB4):
		return model.ExpenseSubtypeB4, nil
	case string(model.ExpenseSubtypeB5):
		return model.ExpenseSubtypeB5, nil
	case string(model.ExpenseSubtypeC1):
		return model.ExpenseSubtypeC1, nil
	case string(model.ExpenseSubtypeC2):
		return model.ExpenseSubtypeC2, nil
	default:
		return "", fmt.Errorf("unknown expense subtype: %s", value)
	}
}

func parseBoolean(s string) bool {
	return strings.ToLower(s) == "true"
}

func expandPath(p string) (string, error) {
	p = os.ExpandEnv(p)
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if p == "~" {
			return home, nil
		}
		if strings.HasPrefix(p, "~/") {
			return filepath.Join(home, p[2:]), nil
		}
	}
	return p, nil
}
