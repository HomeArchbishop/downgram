package binding

type tBindVar[T any] struct {
	data T
}

func (v *tBindVar[T]) Set(value T) {
	v.data = value
	window.Invalidate()
}

func (v *tBindVar[T]) Get() T {
	return v.data
}

func BindVar[T any](hdata T) tBindVar[T] {
	return tBindVar[T]{data: hdata}
}
