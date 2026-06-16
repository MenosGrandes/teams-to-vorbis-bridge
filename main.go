package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"vgc/main/constants"
	utils "vgc/main/excel"
	"vgc/main/student"

	"github.com/xuri/excelize/v2"
)

//

// Names are always in second, split them
func readStudentsInto(f *excelize.File, stundet_name_row int) ([]*student.Student, error) {
	studets := make([]*student.Student, 0)

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(constants.SHEET_NAME)
	if err != nil {
		fmt.Println(err)
		return studets, nil
	}

	for _, row := range rows[1:] {
		studets = append(studets, student.NewStudent(row[stundet_name_row], 0))
	}
	return studets, nil
}
func readStudentsFrom(f *excelize.File, row_number_first_name, row_number_second_name, row_number_grade int) ([]*student.Student, error) {
	studets := make([]*student.Student, 0)

	rows, err := f.GetRows(constants.SHEET_NAME)
	if err != nil {
		fmt.Println(err)
		return studets, nil
	}

	for _, row := range rows[1:] {
		name_array := make([]string, 2)
		name_array = append(name_array, row[row_number_first_name], row[row_number_second_name])
		name := strings.Join(name_array, " ")
		grade, error_grade := strconv.ParseFloat(row[row_number_grade], 32)
		if error_grade != nil {
			fmt.Println(err)
			return studets, &student.StudentError{Message: "Grade is not int in excel", Name: name}
		}
		studets = append(studets, student.NewStudent(name, float32(grade)))
	}
	return studets, nil
}
func validateFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", path)
		}
		return fmt.Errorf("failed to check file %s: %w", path, err)
	}
	return nil
}

func main() {

	into_file_path := flag.String("into", "", "input file path")
	from_file_path := flag.String("from", "", "from file path")
	output_file_path := flag.String("output", "", "output file path")

	flag.Parse()

	if *into_file_path == "" || *from_file_path == "" {
		fmt.Println("usage: app --into <file> --from <file>")
		return
	}
	if err := validateFile(*into_file_path); err != nil {
		fmt.Println("error:", err)
		return
	}
	if err := validateFile(*from_file_path); err != nil {
		fmt.Println("error:", err)
		return
	}
	//
	into_file, err := excelize.OpenFile(*into_file_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := into_file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	students_into, eror := readStudentsInto(into_file, 1)

	if eror != nil {
		fmt.Println(eror)
		return
	}

	from_file, err := excelize.OpenFile(*from_file_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := from_file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	students_from, eror__fropm := readStudentsFrom(from_file, 0, 1, 7)

	if eror__fropm != nil {
		fmt.Println(eror__fropm)
		return
	}
	student.SetupGrades(students_from, students_into)
	//students_into write to file. Loop over into.xml, find all student and put the grade.
	//Always in E column
	utils.SaveStudentGradesToCSV(students_from, into_file, "E", *output_file_path)

}
