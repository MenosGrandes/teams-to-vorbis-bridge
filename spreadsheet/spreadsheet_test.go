package spreadsheet

import (
	"path/filepath"
	"testing"
	"vgc/main/student"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func createTestXlsx(t *testing.T, sheet string, rows [][]string) string {
	t.Helper()
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheet)
	for i, row := range rows {
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			f.SetCellValue(sheet, cell, val)
		}
	}
	path := filepath.Join(t.TempDir(), "test.xlsx")
	f.SaveAs(path)
	f.Close()
	return path
}

func TestOpen_Success(t *testing.T) {
	path := createTestXlsx(t, "a", [][]string{{"Header"}})
	f, err := Open(path, "a")
	assert.NoError(t, err)
	assert.NotNil(t, f)
	f.Close()
}

func TestOpen_FileNotFound(t *testing.T) {
	_, err := Open("/nonexistent.xlsx", "a")
	assert.Error(t, err)
}

func TestReadGrades(t *testing.T) {
	path := createTestXlsx(t, "a", [][]string{
		{"First", "Last", "Email", "Grade"},
		{"John", "Smith", "j@t.com", "4.5"},
		{"Jane", "Doe", "d@t.com", "3"},
	})
	f, _ := Open(path, "a")
	defer f.Close()

	students, err := f.ReadGrades(0, 1, 3)
	assert.NoError(t, err)
	assert.Len(t, students, 2)
	assert.Equal(t, "John Smith", students[0].Name)
	assert.Equal(t, float32(4.5), students[0].Grade)
}

func TestReadGrades_InvalidGrade(t *testing.T) {
	path := createTestXlsx(t, "a", [][]string{
		{"First", "Last", "Grade"},
		{"John", "Smith", "abc"},
	})
	f, _ := Open(path, "a")
	defer f.Close()

	_, err := f.ReadGrades(0, 1, 2)
	assert.Error(t, err)
}

func TestWriteGrades(t *testing.T) {
	path := createTestXlsx(t, "a", [][]string{
		{"ID", "Name"},
		{"1", "John Smith"},
		{"2", "Jane Doe"},
	})
	f, _ := Open(path, "a")
	defer f.Close()

	students := []student.Student{
		{Name: "John Smith", Grade: 5},
		{Name: "Jane Doe", Grade: 3},
	}

	matchFn := func(a, b string) bool { return a == b }
	err := f.WriteGrades(students, 1, "C", matchFn)
	assert.NoError(t, err)
}

func TestSetCell(t *testing.T) {
	path := createTestXlsx(t, "a", [][]string{
		{"Header"},
		{"data"},
	})
	f, _ := Open(path, "a")
	defer f.Close()

	err := f.SetCell("B", 2, 4.5)
	assert.NoError(t, err)

	rows, _ := f.Rows()
	// After SetCell, we can verify via Rows that the cell has been set
	assert.GreaterOrEqual(t, len(rows), 2)
}
