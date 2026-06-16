package student

import (
	"math/rand"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/suite"
)

type NameSpec struct {
	Tokens            int
	UnicodeRatio      float64
	QuestionMarkRatio float64
}

type Generator struct {
	r *rand.Rand
}

func New(seed int64) *Generator {
	return &Generator{
		r: rand.New(rand.NewSource(seed)),
	}
}

var syllables = []string{
	"al", "an", "ar", "ber", "dan", "el", "mar", "mir",
	"tan", "tor", "vin", "zor", "kan", "len", "rik",
}

var unicodeLetters = []rune{
	'ğ', 'ü', 'ö', 'ş', 'ç', 'ı', 'İ',
	'à', 'á', 'â', 'ä', 'ã', 'å', 'æ',
	'è', 'é', 'ê', 'ë',
	'ì', 'í', 'î', 'ï',
	'ñ', 'ò', 'ó', 'ô', 'ö', 'õ', 'ø',
	'ù', 'ú', 'û', 'ü',
	'č', 'š', 'ž', 'ł', 'ń', 'ą', 'ę',
}

func (g *Generator) randomToken() string {
	a := syllables[g.r.Intn(len(syllables))]
	b := syllables[g.r.Intn(len(syllables))]
	return strings.Title(a + b)
}

func (g *Generator) mutate(token string, spec NameSpec) string {
	r := []rune(token)

	for i := range r {
		p := g.r.Float64()

		if p < spec.QuestionMarkRatio {
			r[i] = '?'
			continue
		}

		if p < spec.QuestionMarkRatio+spec.UnicodeRatio && unicode.IsLetter(r[i]) {
			r[i] = unicodeLetters[g.r.Intn(len(unicodeLetters))]
		}
	}

	return string(r)
}

func (g *Generator) Generate(spec NameSpec) string {
	tokens := make([]string, spec.Tokens)

	for i := 0; i < spec.Tokens; i++ {
		tokens[i] = g.mutate(g.randomToken(), spec)
	}

	return strings.Join(tokens, " ")
}

type StudentSuite struct {
	suite.Suite
	gen  *Generator
	spec *NameSpec
}

func (s *StudentSuite) SetupTest() {
	// deterministic seed for reproducible tests
	s.gen = New(1)

}

func (s *StudentSuite) TestGenerateBasic() {
	spec := NameSpec{
		Tokens:            3,
		UnicodeRatio:      0.2,
		QuestionMarkRatio: 0.1,
	}

	name := s.gen.Generate(spec)

	s.NotEmpty(name)
	s.Contains(name, " ")
}

func (s *StudentSuite) TestManyGenerations() {
	spec := NameSpec{
		Tokens:            3,
		UnicodeRatio:      0.3,
		QuestionMarkRatio: 0.2,
	}

	for i := 0; i < 100; i++ {
		name := s.gen.Generate(spec)

		s.NotEmpty(name)
	}
}

func (s *StudentSuite) TestEqualByNameSame() {
	spec := NameSpec{
		Tokens:            3,
		UnicodeRatio:      0.2,
		QuestionMarkRatio: 0.1,
	}
	name := s.gen.Generate(spec)

	student1 := NewStudent(name, 1)
	s.True(student1.equalByName(student1))

}

func (s *StudentSuite) Test_EqualByName_NonUtfSameNamePermute_ShouldReturnTrue() {

	student1 := NewStudent("Abdülbakioğlu Batuhan", 1)
	student2 := NewStudent(" Batuhan Abdülbakioğlu", 1)

	s.True(student1.equalByName(student2))
	s.True(student2.equalByName(student1))

}
func (s *StudentSuite) Test_EqualByName_NonUtfSameNamePermuteWildcard_ShouldReturnTrue() {

	student1 := NewStudent("Abdülbakio?lu Batuhan", 1)
	student2 := NewStudent(" Batuhan Abdülbakioğlu", 1)

	s.True(student1.equalByName(student2))
	s.True(student2.equalByName(student1))

}
func (s *StudentSuite) Test_EqualByName_NonUtfDifferentNamePermuteWildcard_ShouldReturnFalse() {
	spec := NameSpec{
		Tokens:            3,
		UnicodeRatio:      0.2,
		QuestionMarkRatio: 0.1,
	}
	student1 := NewStudent(s.gen.Generate(spec), 1)
	student2 := NewStudent(s.gen.Generate(spec), 1)

	s.False(student1.equalByName(student2))
	s.False(student2.equalByName(student1))

}

func (s *StudentSuite) Test_FindStudent_NonUtfWholeDifferentNamePermuteWildcard_ShouldReturnTrue() {
	spec := NameSpec{
		Tokens:            3,
		UnicodeRatio:      0.2,
		QuestionMarkRatio: 0.1,
	}

	foundStudentName := s.gen.Generate(spec)
	students := make([]*Student, 0)
	students = append(students, NewStudent(s.gen.Generate(spec), 1), NewStudent(foundStudentName, 1))
	student := NewStudent(foundStudentName, 1)
	{
		found_student, error := student.FindStudent(students)
		s.NoError(error)
		s.NotNil(found_student)
		s.Equal(found_student.Name, foundStudentName)

	}

}
func (s *StudentSuite) Test_FindStudent_NonUtfWholeDifferentNamePermuteWildcard_ShouldReturnFalse() {
	spec := NameSpec{
		Tokens:            3,
		UnicodeRatio:      0.2,
		QuestionMarkRatio: 0.1,
	}
	students := make([]*Student, 0)
	students = append(students, NewStudent(s.gen.Generate(spec), 1), NewStudent(s.gen.Generate(spec), 1))
	student := NewStudent(s.gen.Generate(spec), 1)
	{
		found_student, error := student.FindStudent(students)
		s.Error(error)
		s.Nil(found_student)

	}

}
func TestStudentSuite(t *testing.T) {
	suite.Run(t, new(StudentSuite))
}
