package spreadsheet

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
	"vgc/main/student"
)

type ExcelFile struct {
	xlsx  *excelize.File
	sheet string
}

func ColToIndex(col string) (int, error) {
	n, err := excelize.ColumnNameToNumber(col)
	if err != nil {
		return 0, fmt.Errorf("invalid column %q: %w", col, err)
	}
	return n - 1, nil
}

func Open(path, sheet string) (*ExcelFile, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	return &ExcelFile{xlsx: f, sheet: sheet}, nil
}

func (f *ExcelFile) Close() error {
	return f.xlsx.Close()
}

func (f *ExcelFile) ReadGrades(firstNameCol, lastNameCol, gradeCol int) ([]student.Student, error) {
	rows, err := f.xlsx.GetRows(f.sheet)
	if err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}
	students := make([]student.Student, 0, len(rows)-1)
	for _, row := range rows[1:] {
		if gradeCol >= len(row) {
			continue
		}
		name := row[firstNameCol] + " " + row[lastNameCol]
		grade, err := strconv.ParseFloat(row[gradeCol], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid grade for %s: %w", name, err)
		}
		students = append(students, student.Student{Name: name, Grade: float32(grade)})
	}
	return students, nil
}

func (f *ExcelFile) Rows() ([][]string, error) {
	return f.xlsx.GetRows(f.sheet)
}

func (f *ExcelFile) SetCell(col string, row int, value interface{}) error {
	cell := col + strconv.Itoa(row)
	return f.xlsx.SetCellValue(f.sheet, cell, value)
}

func (f *ExcelFile) WriteGrades(students []student.Student, nameCol int, gradeCol string, matchFn func(a, b string) bool) error {
	if err := f.SetCell(gradeCol, 1, "FINAL GRADE"); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	rows, err := f.xlsx.GetRows(f.sheet)
	if err != nil {
		return fmt.Errorf("reading rows: %w", err)
	}

	for i, row := range rows[1:] {
		if nameCol >= len(row) {
			continue
		}
		name := row[nameCol]
		for _, s := range students {
			if matchFn(s.Name, name) {
				if err := f.SetCell(gradeCol, i+2, s.Grade); err != nil {
					return fmt.Errorf("writing grade row %d: %w", i+2, err)
				}
				break
			}
		}
	}
	return nil
}
