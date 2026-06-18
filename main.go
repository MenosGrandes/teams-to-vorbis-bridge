package main

import (
	"flag"
	"fmt"
	"os"

	"vgc/main/engine/impl"
	exportimpl "vgc/main/export/impl"
	"vgc/main/grading"
	"vgc/main/matching"
	"vgc/main/pipeline"
	"vgc/main/spreadsheet"
)

type config struct {
	from        string
	into        string
	output      string
	gradeConfig string
	sheet       string
}

func parseFlags() (*config, error) {
	cfg := &config{}
	flag.StringVar(&cfg.into, "into", "", "Teams export file path")
	flag.StringVar(&cfg.from, "from", "", "Vorbis grades file path")
	flag.StringVar(&cfg.output, "output", "output.csv", "output CSV file path")
	flag.StringVar(&cfg.gradeConfig, "grade-config", "", "grading YAML config path")
	flag.StringVar(&cfg.sheet, "sheet", "a", "sheet name in xlsx files")
	flag.Parse()

	if cfg.into == "" || cfg.from == "" || cfg.gradeConfig == "" {
		return nil, fmt.Errorf("usage: app --into <file> --from <file> --grade-config <file> [--output <file>] [--sheet <name>]")
	}

	for _, p := range []string{cfg.into, cfg.from, cfg.gradeConfig} {
		if _, err := os.Stat(p); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func run() error {
	cfg, err := parseFlags()
	if err != nil {
		return err
	}

	gradeCfg, err := grading.LoadConfig(cfg.gradeConfig)
	if err != nil {
		return err
	}

	skipCols, err := spreadsheet.ColToIndex(gradeCfg.SkipCols)
	if err != nil {
		fmt.Println("DUPA1")
		return err
	}
	fromFirstCol, err := spreadsheet.ColToIndex(gradeCfg.FromFirstNameCol)
	if err != nil {
				fmt.Println("DUPA2")

		return err
	}
	fromLastCol, err := spreadsheet.ColToIndex(gradeCfg.FromLastNameCol)
	if err != nil {
				fmt.Println("DUPA3")

		return err
	}
	fromGradeCol, err := spreadsheet.ColToIndex(gradeCfg.FromGradeCol)
	if err != nil {
				fmt.Println("DUPA4")

		return err
	}
	intoNameCol, err := spreadsheet.ColToIndex(gradeCfg.IntoNameCol)
	if err != nil {
				fmt.Println("DUPA15")

		return err
	}

	eng, err := impl.NewLuaEngine(gradeCfg.Formula, gradeCfg.Grades)
	if err != nil {
		return err
	}
	defer eng.Close()

	fromFile, err := spreadsheet.Open(cfg.from, cfg.sheet)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	rows, err := fromFile.Rows()
	if err != nil {
		return err
	}
	grades, err := pipeline.Grade(eng, rows[1:], skipCols)
	if err != nil {
		return err
	}
	for i, grade := range grades {
		if err := fromFile.SetCell(gradeCfg.ResultColumn, i+2, grade); err != nil {
			return err
		}
	}

	studentsFrom, err := fromFile.ReadGrades(fromFirstCol, fromLastCol, fromGradeCol)
	if err != nil {
		return err
	}

	intoFile, err := spreadsheet.Open(cfg.into, cfg.sheet)
	if err != nil {
		return err
	}
	defer intoFile.Close()

	if err := intoFile.WriteGrades(studentsFrom, intoNameCol, gradeCfg.IntoGradeCol, matching.IsMatch); err != nil {
		return err
	}

	rows, err = intoFile.Rows()
	if err != nil {
		return err
	}

	return exportimpl.NewCSVExporter().Export(rows, cfg.output)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}