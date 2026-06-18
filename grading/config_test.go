package grading

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_ValidFile(t *testing.T) {
	yaml := "result_column: H\ngrades:\n  0: 2\n  15: 3\n  28: 5\nformula: \"return getGrade(D)\"\n"
	tmp := filepath.Join(t.TempDir(), "config.yaml")
	os.WriteFile(tmp, []byte(yaml), 0644)

	cfg, err := LoadConfig(tmp)
	assert.NoError(t, err)
	assert.Equal(t, "H", cfg.ResultColumn)
	assert.Equal(t, float64(5), cfg.Grades[28])
}

func TestLoadConfig_NotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent.yaml")
	assert.Error(t, err)
}
