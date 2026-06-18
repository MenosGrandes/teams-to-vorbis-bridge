package impl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testGrades = map[int]float64{0: 2, 15: 3, 20: 3.5, 22: 4, 26: 4.5, 28: 5}

func TestNewLuaEngine_ValidFormula(t *testing.T) {
	eng, err := NewLuaEngine("return getGrade(D + E)", testGrades)
	assert.NoError(t, err)
	assert.NotNil(t, eng)
}

func TestNewLuaEngine_InvalidFormula(t *testing.T) {
	_, err := NewLuaEngine("if D > 0 then", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid formula")
}

func TestEvaluate_SimpleSum(t *testing.T) {
	eng, _ := NewLuaEngine("return getGrade(D + E)", testGrades)
	grade, err := eng.Evaluate(map[string]float64{"D": 10, "E": 15})
	assert.NoError(t, err)
	assert.Equal(t, float32(4), grade) // 25 >= 22
}

func TestEvaluate_LowScore(t *testing.T) {
	eng, _ := NewLuaEngine("return getGrade(D + E)", testGrades)
	grade, err := eng.Evaluate(map[string]float64{"D": 3, "E": 5})
	assert.NoError(t, err)
	assert.Equal(t, float32(2), grade) // 8 >= 0
}

func TestEvaluate_MaxScore(t *testing.T) {
	eng, _ := NewLuaEngine("return getGrade(D + E)", testGrades)
	grade, err := eng.Evaluate(map[string]float64{"D": 15, "E": 15})
	assert.NoError(t, err)
	assert.Equal(t, float32(5), grade) // 30 >= 28
}

func TestEvaluate_Conditional(t *testing.T) {
	grades := map[int]float64{0: 2, 15: 3, 28: 5}
	eng, _ := NewLuaEngine("if D > 0 then return getGrade(D) end\nreturn getGrade(E + F)", grades)

	grade, err := eng.Evaluate(map[string]float64{"D": 28, "E": 0, "F": 0})
	assert.NoError(t, err)
	assert.Equal(t, float32(5), grade)

	grade, err = eng.Evaluate(map[string]float64{"D": 0, "E": 10, "F": 8})
	assert.NoError(t, err)
	assert.Equal(t, float32(3), grade) // 18 >= 15
}

func TestEvaluate_ZeroScore(t *testing.T) {
	eng, _ := NewLuaEngine("return getGrade(D + E)", testGrades)
	grade, err := eng.Evaluate(map[string]float64{"D": 0, "E": 0})
	assert.NoError(t, err)
	assert.Equal(t, float32(2), grade) // 0 >= 0
}
