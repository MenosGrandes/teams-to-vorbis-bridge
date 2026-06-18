package engine

type Engine interface {
	Evaluate(cols map[string]float64) (float32, error)
}
