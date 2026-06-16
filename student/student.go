package student

import (
	"fmt"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type Student struct {
	Name  string
	Grade float32
}
type StudentError struct {
	Name    string
	Message string
}

func (e *StudentError) Error() string {
	return fmt.Sprintf(" %s -  %s", e.Message, e.Name)
}
func (s *Student) FindStudent(students []*Student) (*Student, error) {

	for _, student := range students {
		if student.equalByName(s) {
			return student, nil
		}
	}
	return nil, &StudentError{"NotFound", s.Name}
}
func NewStudent(Name string, Grade float32) *Student {
	p := Student{Name: Name, Grade: Grade}
	return &p
}
func (s *Student) print() {
	fmt.Printf("%s - %f\n", s.Name, s.Grade)
}
func (s *Student) equalByName(other *Student) bool {
	return IsNameMatch(s.Name, other.Name)
}
func printStudents(students []*Student) {

	for _, s := range students {
		s.print()
	}
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = norm.NFD.String(s)

	var b strings.Builder
	for _, r := range s {
		// remove Non UTF characters
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		//IN those strings som non UTF chars are ?
		//So it's a wildcard
		if unicode.IsLetter(r) || r == '?' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func tokenMatch(a, b string) bool {
	ar := []rune(a)
	br := []rune(b)

	if len(ar) != len(br) {
		return false
	}
	//IN those strings som non UTF chars are ?
	//So it's a wildcard
	for i := range ar {
		if ar[i] == '?' || br[i] == '?' {
			continue
		}
		if ar[i] != br[i] {
			return false
		}
	}
	return true
}

// tokenize Name into normalized words
func tokenize(Name string) []string {
	parts := strings.Fields(Name)

	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := normalize(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func IsNameMatch(a, b string) bool {
	A := tokenize(a)
	B := tokenize(b)

	if len(A) == 0 || len(B) == 0 {
		return false
	}

	matchedB := make([]bool, len(B))
	matchCount := 0

	// greedy matching
	for _, ta := range A {
		for j, tb := range B {
			if matchedB[j] {
				continue
			}
			if tokenMatch(ta, tb) {
				matchedB[j] = true
				matchCount++
				break
			}
		}
	}

	coverageA := float64(matchCount) / float64(len(A))
	coverageB := float64(matchCount) / float64(len(B))

	// accept if either side is mostly covered
	return coverageA >= 0.9 || coverageB >= 0.9
}
func SetupGrades(students_from, students_into []*Student) {

	const workers = 8

	jobs := make(chan *Student)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for st := range jobs {
			s, err := st.FindStudent(students_from)
			if err == nil {
				st.Grade = s.Grade
			}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	for _, st := range students_into {
		jobs <- st
	}

	close(jobs)
	wg.Wait()

}
