package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"vgc/main/constants"
	"vgc/main/student"

	"github.com/xuri/excelize/v2"
)

type Axis struct {
	row int
	col string
}

// The student names are always in 1st column
func SaveStudentGradesToCSV(students []*student.Student, f *excelize.File, gradeColumName, output_file_path string) {
	header_cell := string(gradeColumName + strconv.Itoa(1))

	sheet_err := f.SetSheetRow(constants.SHEET_NAME, header_cell, &[]any{"FINAL GRADE"})
	if sheet_err != nil {
		log.Fatal(sheet_err)

	}
	rows, err := f.GetRows(constants.SHEET_NAME)
	if err != nil {
		fmt.Println(err)
	}
	for i, row := range rows[1:] {
		tmp_student := student.NewStudent(row[1], 0)
		found, student_not_found := tmp_student.FindStudent(students)
		if student_not_found != nil {
			continue
		}

		cell := string(gradeColumName + strconv.Itoa(i+2))
		sheet_err := f.SetSheetRow(constants.SHEET_NAME, cell, &[]any{found.Grade})
		if sheet_err != nil {
			log.Fatal(sheet_err)

		}
	}

	err = generateCSV(f, Axis{1, "A"}, Axis{200, "J"}, output_file_path)

	if err != nil {
		log.Fatal(err)
	}

}

func generateCSV(f *excelize.File, start, end Axis, output_file_path string) error {
	var data [][]string

	for i := start.row; i <= end.row; i++ {
		row := []string{}
		for j := []rune(start.col)[0]; j <= []rune(end.col)[0]; j++ {
			value, err := f.GetCellValue(constants.SHEET_NAME, fmt.Sprintf("%s%d", string(j), i), excelize.Options{})
			if err != nil {
				return err
			}
			row = append(row, value)
		}
		data = append(data, row)
	}

	file, err := os.Create(output_file_path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(file)
	return writer.WriteAll(data)
}
func PrintSpreadsheet(f *excelize.File) {
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(constants.SHEET_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
