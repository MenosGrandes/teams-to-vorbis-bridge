package impl

import (
	"encoding/csv"
	"fmt"
	"os"

	"vgc/main/export"
)

var _ export.Exporter = (*CSVExporter)(nil)

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

func (e *CSVExporter) Export(rows [][]string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating CSV: %w", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	return w.WriteAll(rows)
}
