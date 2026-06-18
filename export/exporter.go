package export

type Exporter interface {
	Export(rows [][]string, path string) error
}
