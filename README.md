# teams-to-vorbis-bridge

CLI tool that computes student grades from a Vorbis spreadsheet using a configurable Lua formula, matches students to an MS Teams export by fuzzy name matching, and outputs a CSV with final grades.

## Usage

```bash
go build -o vgc .
./vgc --from vorbis.xlsx --into teams.xlsx --grade-config config.yaml [--output grades.csv] [--sheet a]
```

### Prerequisites format of XLSX
There should be no strings in XLSX files, as all cells are formatted as int.


### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--from` | yes | | Vorbis grades xlsx file |
| `--into` | yes | | MS Teams export xlsx file |
| `--grade-config` | yes | | Grading config YAML |
| `--output` | no | `output.csv` | Output CSV path |
| `--sheet` | no | `a` | Sheet name in xlsx files |

## Configuration

`config.yaml` defines column mappings, grade thresholds, and the grading formula:

```yaml
result_column: H
from_first_name_col: A
from_last_name_col: B
from_grade_col: H
into_name_col: B
into_grade_col: E
skip_cols: D

grades:
  0: 2
  15: 3
  20: 3.5
  22: 4
  26: 4.5
  28: 5

formula: |
  if D > 0 then
    return getGrade(D)
  end
  return getGrade(E + F)
```

### Config fields

| Field | Description |
|-------|-------------|
| `result_column` | Column to write computed grade in the "from" file |
| `from_first_name_col` | First name column in Vorbis file |
| `from_last_name_col` | Last name column in Vorbis file |
| `from_grade_col` | Column to read final grade from (after formula is applied) |
| `into_name_col` | Student name column in Teams file |
| `into_grade_col` | Column to write grade in Teams file |
| `skip_cols` | First column exposed to the formula (columns before this are skipped) |
| `grades` | Score threshold → grade mapping (highest threshold ≤ score wins) |
| `formula` | Lua code evaluated per student row |

### Formula

The formula is Lua. Column letters (D, E, F, ...) from `skip_cols` onward are available as numeric variables. Empty cells are 0.

Built-in function:
- `getGrade(score)` — looks up the grade for a given score using the `grades` threshold table

Standard Lua libraries (`math`, `string`, etc.) are available.

## Name matching

Students are matched between the two files by fuzzy name comparison:
- Unicode normalization (diacritics stripped for comparison)
- Token order independent ("Smith John" matches "John Smith")
- `?` wildcard support (matches any single character)
- 90% token coverage threshold

## Project structure

```
├── main.go              # CLI wiring
├── student/             # Student domain type
├── matching/            # Fuzzy name matching
├── grading/             # YAML config loading
├── engine/              # Engine interface
│   └── impl/            # Lua implementation
├── export/              # Exporter interface
│   └── impl/            # CSV implementation
├── pipeline/            # Grading logic (score → grade)
├── spreadsheet/         # Excel file I/O
└── testutil/            # Test name generator
```

## Testing

```bash
go test ./...
```

Generate test data with diacritics and wildcards:

```bash
go run testdata/gen.go
```

## Requirements

- Go 1.24+
