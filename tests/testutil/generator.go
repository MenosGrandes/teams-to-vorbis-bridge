package testutil

import (
	"math/rand"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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

type NameSpec struct {
	Tokens            int
	UnicodeRatio      float64
	QuestionMarkRatio float64
}

type Generator struct {
	R *rand.Rand
}

func NewGenerator(seed int64) *Generator {
	return &Generator{R: rand.New(rand.NewSource(seed))}
}

func (g *Generator) RandomToken() string {
	a := syllables[g.R.Intn(len(syllables))]
	b := syllables[g.R.Intn(len(syllables))]
	return cases.Title(language.Und).String(a + b)
}

func (g *Generator) Mutate(token string, spec NameSpec) string {
	r := []rune(token)
	for i := range r {
		p := g.R.Float64()
		if p < spec.QuestionMarkRatio {
			r[i] = '?'
			continue
		}
		if p < spec.QuestionMarkRatio+spec.UnicodeRatio && unicode.IsLetter(r[i]) {
			r[i] = unicodeLetters[g.R.Intn(len(unicodeLetters))]
		}
	}
	return string(r)
}

func (g *Generator) Generate(spec NameSpec) string {
	tokens := make([]string, spec.Tokens)
	for i := 0; i < spec.Tokens; i++ {
		tokens[i] = g.Mutate(g.RandomToken(), spec)
	}
	return strings.Join(tokens, " ")
}

func (g *Generator) GenerateName(unicodeRatio float64) (string, string) {
	spec := NameSpec{Tokens: 1, UnicodeRatio: unicodeRatio}
	first := g.Mutate(g.RandomToken(), spec)
	last := g.Mutate(g.RandomToken(), spec)
	return first, last
}

func (g *Generator) AddWildcard(name string) string {
	r := []rune(name)
	if len(r) == 0 {
		return name
	}
	idx := g.R.Intn(len(r))
	for !unicode.IsLetter(r[idx]) {
		idx = g.R.Intn(len(r))
	}
	r[idx] = '?'
	return string(r)
}
