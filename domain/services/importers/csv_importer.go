package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pjover/espigol/domain/interfaces"
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

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read row: %w", err)
		}

		for i, h := range header {
			val := ""
			if i < len(rec) {
				val = rec[i]
			}
			fmt.Printf("%s: %s\n", h, val)
		}
		fmt.Println("---")
	}

	return nil
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
