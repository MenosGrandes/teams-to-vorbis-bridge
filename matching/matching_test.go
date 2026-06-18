package matching

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"vgc/main/tests/testutil"
)

type MatchingSuite struct {
	suite.Suite
	gen *testutil.Generator
}

func (s *MatchingSuite) SetupTest() {
	s.gen = testutil.NewGenerator(1)
}

func (s *MatchingSuite) TestGenerateBasic() {
	spec := testutil.NameSpec{Tokens: 3, UnicodeRatio: 0.2, QuestionMarkRatio: 0.1}
	name := s.gen.Generate(spec)
	s.NotEmpty(name)
	s.Contains(name, " ")
}

func (s *MatchingSuite) TestManyGenerations() {
	spec := testutil.NameSpec{Tokens: 3, UnicodeRatio: 0.3, QuestionMarkRatio: 0.2}
	for i := 0; i < 100; i++ {
		s.NotEmpty(s.gen.Generate(spec))
	}
}

func (s *MatchingSuite) TestSameNameMatches() {
	spec := testutil.NameSpec{Tokens: 3, UnicodeRatio: 0.2, QuestionMarkRatio: 0.1}
	name := s.gen.Generate(spec)
	s.True(IsMatch(name, name))
}

func (s *MatchingSuite) Test_NonUtfSameNamePermute_ShouldReturnTrue() {
	s.True(IsMatch("Abd\u00fclbakio\u011flu Batuhan", "Batuhan Abd\u00fclbakio\u011flu"))
	s.True(IsMatch("Batuhan Abd\u00fclbakio\u011flu", "Abd\u00fclbakio\u011flu Batuhan"))
}

func (s *MatchingSuite) Test_NonUtfSameNamePermuteWildcard_ShouldReturnTrue() {
	s.True(IsMatch("Abd\u00fclbakio?lu Batuhan", "Batuhan Abd\u00fclbakio\u011flu"))
	s.True(IsMatch("Batuhan Abd\u00fclbakio\u011flu", "Abd\u00fclbakio?lu Batuhan"))
}

func (s *MatchingSuite) Test_DifferentNames_ShouldReturnFalse() {
	spec := testutil.NameSpec{Tokens: 3, UnicodeRatio: 0.2, QuestionMarkRatio: 0.1}
	nameA := s.gen.Generate(spec)
	nameB := s.gen.Generate(spec)
	s.False(IsMatch(nameA, nameB))
	s.False(IsMatch(nameB, nameA))
}

func (s *MatchingSuite) Test_FindInList_ShouldReturnTrue() {
	spec := testutil.NameSpec{Tokens: 3, UnicodeRatio: 0.2, QuestionMarkRatio: 0.1}
	target := s.gen.Generate(spec)
	other := s.gen.Generate(spec)
	names := []string{other, target}
	found := false
	for _, n := range names {
		if IsMatch(n, target) {
			found = true
			break
		}
	}
	s.True(found)
}

func (s *MatchingSuite) Test_FindInList_ShouldReturnFalse() {
	spec := testutil.NameSpec{Tokens: 3, UnicodeRatio: 0.2, QuestionMarkRatio: 0.1}
	target := s.gen.Generate(spec)
	names := []string{s.gen.Generate(spec), s.gen.Generate(spec)}
	found := false
	for _, n := range names {
		if IsMatch(n, target) {
			found = true
			break
		}
	}
	s.False(found)
}

func TestMatchingSuite(t *testing.T) {
	suite.Run(t, new(MatchingSuite))
}
