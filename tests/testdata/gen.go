//go:build ignore

package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"vgc/main/tests/testutil"
)

type fromRecord struct {
	FirstName string
	LastName  string
	Email     string
	Test1     int
	Test2     int
	Resit     int
}

type intoRecord struct {
	ID         string
	Name       string
	Email      string
	Submission string
}

func main() {
	const numStudents = 15
	g := testutil.NewGenerator(42)

	fromStudents := make([]fromRecord, numStudents)
	for i := range fromStudents {
		first, last := g.GenerateName(0.3)
		s := fromRecord{FirstName: first, LastName: last, Email: fmt.Sprintf("student%d@uni.pl", i+1)}
		if g.R.Float64() < 0.25 {
			s.Resit = 28 + g.R.Intn(3)
		} else {
			s.Test1 = g.R.Intn(16)
			s.Test2 = g.R.Intn(16)
		}
		fromStudents[i] = s
	}

	from := excelize.NewFile()
	from.SetSheetName("Sheet1", "a")
	headers := []interface{}{"First Name", "Last Name", "Email", "Resit", "Second Test", "First Test"}
	for j, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(j+1, 1)
		from.SetCellValue("a", cell, h)
	}
	for i, s := range fromStudents {
		row := []interface{}{s.FirstName, s.LastName, s.Email, s.Resit, s.Test2, s.Test1}
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			from.SetCellValue("a", cell, val)
		}
	}
	from.SaveAs("tests/testdata/from.xlsx")
	from.Close()

	intoStudents := make([]intoRecord, numStudents)
	for i, s := range fromStudents {
		var name string
		if g.R.Float64() < 0.5 {
			name = s.LastName + " " + s.FirstName
		} else {
			name = s.FirstName + " " + s.LastName
		}
		if g.R.Float64() < 0.3 {
			name = g.AddWildcard(name)
		}
		intoStudents[i] = intoRecord{
			ID:         fmt.Sprintf("%d", i+1),
			Name:       name,
			Email:      s.Email,
			Submission: "submitted",
		}
	}

	into := excelize.NewFile()
	into.SetSheetName("Sheet1", "a")
	intoHeaders := []interface{}{"ID", "Student Name", "Email", "Submission"}
	for j, h := range intoHeaders {
		cell, _ := excelize.CoordinatesToCellName(j+1, 1)
		into.SetCellValue("a", cell, h)
	}
	for i, s := range intoStudents {
		row := []interface{}{s.ID, s.Name, s.Email, s.Submission}
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			into.SetCellValue("a", cell, val)
		}
	}
	into.SaveAs("tests/testdata/into.xlsx")
	into.Close()

	fmt.Printf("Generated %d students with diacritics, reversed names, and wildcards\n", numStudents)
}
