package game

type (
	Attribute[T int | float64] struct {
		Base       T `yaml:"base"`
		Delta      T `yaml:"delta"`
		TotalValue T `yaml:"total_value"`
	}
)

func NewAttribute[T int | float64](base T) *Attribute[T] {
	return &Attribute[T]{
		Base: base,
	}
}

func (a *Attribute[T]) Recalculate() {
	a.TotalValue = a.Base + a.Delta
}
