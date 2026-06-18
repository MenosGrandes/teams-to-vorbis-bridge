package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEngine struct {
	gradeFn func(cols map[string]float64) (float32, error)
}

func (m *mockEngine) Evaluate(cols map[string]float64) (float32, error) {
	return m.gradeFn(cols)
}

func TestGrade_BasicSum(t *testing.T) {
	eng := &mockEngine{gradeFn: func(cols map[string]float64) (float32, error) {
		return float32(cols["D"] + cols["E"]), nil
	}}

	rows := [][]string{
		{"John", "Smith", "j@t.com", "10", "5"},
		{"Jane", "Doe", "d@t.com", "3", "7"},
	}

	grades, err := Grade(eng, rows, 3)
	assert.NoError(t, err)
	assert.Equal(t, []float32{15, 10}, grades)
}

func TestGrade_EmptyCellsAreZero(t *testing.T) {
	eng := &mockEngine{gradeFn: func(cols map[string]float64) (float32, error) {
		return float32(cols["D"]), nil
	}}

	rows := [][]string{
		{"John", "Smith", "j@t.com", ""},
	}

	grades, err := Grade(eng, rows, 3)
	assert.NoError(t, err)
	assert.Equal(t, []float32{0}, grades)
}

func TestGrade_InvalidNumber(t *testing.T) {
	eng := &mockEngine{gradeFn: func(cols map[string]float64) (float32, error) {
		return 0, nil
	}}

	rows := [][]string{
		{"John", "Smith", "j@t.com", "abc"},
	}

	_, err := Grade(eng, rows, 3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid number")
}

func TestGrade_SkipCols(t *testing.T) {
	var receivedCols map[string]float64
	eng := &mockEngine{gradeFn: func(cols map[string]float64) (float32, error) {
		receivedCols = cols
		return 5, nil
	}}

	rows := [][]string{
		{"John", "Smith", "j@t.com", "extra", "10", "20"},
	}

	// skip 4 cols (A,B,C,D) — only E and F should be passed
	grades, err := Grade(eng, rows, 4)
	assert.NoError(t, err)
	assert.Equal(t, []float32{5}, grades)
	_, hasE := receivedCols["E"]
	_, hasF := receivedCols["F"]
	assert.True(t, hasE)
	assert.True(t, hasF)
}

func TestGrade_EmptyRows(t *testing.T) {
	eng := &mockEngine{gradeFn: func(cols map[string]float64) (float32, error) {
		return 5, nil
	}}

	grades, err := Grade(eng, [][]string{}, 3)
	assert.NoError(t, err)
	assert.Empty(t, grades)
}
