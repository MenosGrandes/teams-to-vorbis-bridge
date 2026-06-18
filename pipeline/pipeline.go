package pipeline

import (
	"fmt"
	"strconv"

	"vgc/main/engine"
)

func colLetter(index int) string {
	return string(rune('A' + index))
}

func Grade(eng engine.Engine, rows [][]string, skipCols int) ([]float32, error) {
	grades := make([]float32, len(rows))
	for i, row := range rows {
		cols := make(map[string]float64)
		for colIdx := skipCols; colIdx < len(row); colIdx++ {
			val := 0.0
			if row[colIdx] != "" {
				v, err := strconv.ParseFloat(row[colIdx], 64)
				if err != nil {
					return nil, fmt.Errorf("row %d col %s: invalid number %q: %w", i+2, colLetter(colIdx), row[colIdx], err)
				}
				val = v
			}
			cols[colLetter(colIdx)] = val
		}
		grade, err := eng.Evaluate(cols)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+2, err)
		}
		grades[i] = grade
	}
	return grades, nil
}
