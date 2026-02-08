package interfaces

import "io"

// Importer defines methods to import data sources.
type Importer interface {
    // ImportSocisCSV reads the CSV file at path and writes a human-readable
    // representation to the provided writer.
    ImportSocisCSV(path string, w io.Writer) error
}
