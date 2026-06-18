package grading

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ResultColumn     string          `yaml:"result_column"`
	Grades           map[int]float64 `yaml:"grades"`
	Formula          string          `yaml:"formula"`
	FromFirstNameCol string          `yaml:"from_first_name_col"`
	FromLastNameCol  string          `yaml:"from_last_name_col"`
	FromGradeCol     string          `yaml:"from_grade_col"`
	IntoNameCol      string          `yaml:"into_name_col"`
	IntoGradeCol     string          `yaml:"into_grade_col"`
	SkipCols         string          `yaml:"skip_cols"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}
