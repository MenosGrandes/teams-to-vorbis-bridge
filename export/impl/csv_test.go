package impl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVExporter_Export(t *testing.T) {
	rows := [][]string{
		{"Name", "Grade"},
		{"John", "5"},
		{"Jane", "3"},
	}

	path := filepath.Join(t.TempDir(), "out.csv")
	e := NewCSVExporter()
	err := e.Export(rows, path)
	assert.NoError(t, err)

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "John")
	assert.Contains(t, string(data), "Jane")
}
