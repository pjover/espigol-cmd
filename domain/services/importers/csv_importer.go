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

	"github.com/pjover/espigol/domain/interfaces"
	"github.com/pjover/espigol/domain/model"
)

type CSVImporter struct{}

func NewCSVImporter() interfaces.Importer { return &CSVImporter{} }

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
