package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const sheet_name = "a"

func permute(nums []string) [][]string {
	var result [][]string
	backtrack(&result, nums, []string{})
	return result
}

// backtrack is the recursive helper function for generating permutations.
func backtrack(result *[][]string, remaining []string, currentPermutation []string) {
	if len(remaining) == 0 {
		temp := make([]string, len(currentPermutation))
		copy(temp, currentPermutation)
		*result = append(*result, temp)
		return
	}

	for i := 0; i < len(remaining); i++ {
		element := remaining[i]

		newRemaining := make([]string, 0, len(remaining)-1)
		for j := 0; j < len(remaining); j++ {
			if i != j {
				newRemaining = append(newRemaining, remaining[j])
			}
		}

		newCurrentPermutation := append(currentPermutation, element)
		backtrack(result, newRemaining, newCurrentPermutation)
	}
}

// permute generates all permutations of the given tokens.
func permuteTokens(tokens []string) [][]string {
	var result [][]string
	backtrack(&result, tokens, []string{})
	return result
}

// preprocessAndPermute takes a raw string input  and returns all unique permutations of its constituent words.
func preprocessAndPermute(rawInput string) [][]string {

	tokens := strings.Fields(rawInput)

	if len(tokens) == 0 {
		return nil // Handle empty input
	}
	sort.Strings(tokens)
	permutations := permuteTokens(tokens)
	var finalResults [][]string
	for _, p := range permutations {
		joinedString := strings.Join(p, " ")
		finalResults = append(finalResults, []string{joinedString})
	}

	return finalResults
}

func printSpreadsheet(f *excelize.File) {
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheet_name)
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

type Student struct {
	name  string
	grade float32
}
type StudentError struct {
	name    string
	message string
}

func (e *StudentError) Error() string {
	return fmt.Sprintf(" %s -  %s", e.message, e.name)
}
func (s *Student) findStudent(students []*Student) (*Student, error) {

	for _, student := range students {
		if student.equalByName(s) {
			return student, nil
		}
	}
	return nil, &StudentError{"NotFound", s.name}
}
func newStudent(name string, grade float32) *Student {
	p := Student{name: name, grade: grade}
	return &p
}
func (s *Student) print() {
	fmt.Printf("%#v\n", s)
}

func removeSpacesFromSlice(input []string) []string {
	output := make([]string, len(input))

	for i, s := range input {
		cleanedString := strings.ReplaceAll(s, " ", "")
		output[i] = cleanedString
	}
	return output
}

func (s *Student) equalByName(other *Student) bool {
	//So split Name into array, make all of them lowercase and compare allcombinations?
	s_name_perm := preprocessAndPermute(s.name)
	s_name_other_perm := preprocessAndPermute(other.name)
	for _, s_perm := range s_name_perm {
		s_perm_name := strings.Join(s_perm, " ")

		for _, s_other_perm := range s_name_other_perm {
			s_other_perm_name := strings.Join(s_other_perm, " ")
			if s_perm_name == s_other_perm_name {
				return true
			}
		}
	}

	return false
}
func printStudents(students []*Student) {

	for _, s := range students {
		s.print()
	}
}

// Names are always in second, split them
func readStudentsInto(f *excelize.File, stundet_name_row int) ([]*Student, error) {
	studets := make([]*Student, 0)

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheet_name)
	if err != nil {
		fmt.Println(err)
		return studets, nil
	}

	for _, row := range rows[1:] {
		studets = append(studets, newStudent(row[stundet_name_row], 0))
	}
	return studets, nil
}
func readStudentsFrom(f *excelize.File, row_number_first_name, row_number_second_name, row_number_grade int) ([]*Student, error) {
	studets := make([]*Student, 0)

	rows, err := f.GetRows(sheet_name)
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
			return studets, &StudentError{"Grade is not int in excel", name}
		}
		studets = append(studets, newStudent(name, float32(grade)))
	}
	return studets, nil
}
func main() {
	into_file, err := excelize.OpenFile("into.xlsx")
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

	from_file, err := excelize.OpenFile("from.xlsx")
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
	setupGrades(students_from, students_into)

	//students_into write to file. Loop over into.xml, find all student and put the grade.
	//Always in E column
	saveStudentGradesToCSV(students_from, into_file, "E")

}
func setupGrades(students_from, students_into []*Student) {
	for _, student := range students_into {
		s, err := student.findStudent(students_from)
		if err == nil {
			student.grade = s.grade
		}
	}

}

type Axis struct {
	row int
	col string
}

// The student names are always in 1st column
func saveStudentGradesToCSV(students []*Student, f *excelize.File, gradeColumName string) {
	header_cell := string(gradeColumName + strconv.Itoa(1))

	sheet_err := f.SetSheetRow(sheet_name, header_cell, &[]any{"FINAL GRADE"})
	if sheet_err != nil {
		log.Fatal(sheet_err)

	}
	rows, err := f.GetRows(sheet_name)
	if err != nil {
		fmt.Println(err)
	}
	for i, row := range rows[1:] {
		tmp_student := newStudent(row[1], 0)
		found, student_not_found := tmp_student.findStudent(students)
		if student_not_found != nil {
			continue
		}

		cell := string(gradeColumName + strconv.Itoa(i+2))
		sheet_err := f.SetSheetRow(sheet_name, cell, &[]any{found.grade})
		if sheet_err != nil {
			log.Fatal(sheet_err)

		}
	}

	err = generateCSV(f, Axis{1, "A"}, Axis{200, "J"})

	if err != nil {
		log.Fatal(err)
	}

}

func generateCSV(f *excelize.File, start, end Axis) error {
	var data [][]string

	for i := start.row; i <= end.row; i++ {
		row := []string{}
		for j := []rune(start.col)[0]; j <= []rune(end.col)[0]; j++ {
			value, err := f.GetCellValue(sheet_name, fmt.Sprintf("%s%d", string(j), i), excelize.Options{})
			if err != nil {
				return err
			}
			row = append(row, value)
		}
		data = append(data, row)
	}

	file, err := os.Create("final_grades.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(file)
	return writer.WriteAll(data)
}
