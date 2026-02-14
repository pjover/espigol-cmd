package services

import (
	"os"
	"path/filepath"
	"testing"
)

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

	importer := NewCSVImporter()
	err = importer.ImportPartners(csvPath)
	if err != nil {
		t.Errorf("ImportPartners failed: %v", err)
	}
}

func TestImportPartnersInvalidPath(t *testing.T) {
	importer := NewCSVImporter()
	err := importer.ImportPartners("/nonexistent/path/partners.csv")

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

	importer := NewCSVImporter()
	err = importer.ImportPartners(csvPath)

	if err == nil {
		t.Error("Expected error for invalid data, got nil")
	}
}
