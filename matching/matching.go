package matching

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func IsMatch(a, b string) bool {
	tokensA := tokenize(a)
	tokensB := tokenize(b)

	if len(tokensA) == 0 || len(tokensB) == 0 {
		return false
	}

	matchedB := make([]bool, len(tokensB))
	matchCount := 0

	for _, ta := range tokensA {
		for j, tb := range tokensB {
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

	coverageA := float64(matchCount) / float64(len(tokensA))
	coverageB := float64(matchCount) / float64(len(tokensB))

	return coverageA >= 0.9 || coverageB >= 0.9
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = norm.NFD.String(s)

	var b strings.Builder
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
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

func tokenize(name string) []string {
	parts := strings.Fields(name)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := normalize(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
