package binding

type tBindSlice[T any] struct {
	data []T
}

func (s *tBindSlice[T]) Append(value T) {
	s.data = append(s.data, value)
	window.Invalidate()
}

func (s *tBindSlice[T]) AppendSlice(value []T) {
	s.data = append(s.data, value...)
	window.Invalidate()
}

func (s *tBindSlice[T]) Overwrite(hdate *[]T) {
	s.data = *hdate
	window.Invalidate()
}

func (s *tBindSlice[T]) Len() int {
	return len(s.data)
}

func (s *tBindSlice[T]) Val() *[]T {
	return &s.data
}

func (s *tBindSlice[T]) Get(i int) T {
	return s.data[i]
}

func (s *tBindSlice[T]) Set(i int, value T) {
	s.data[i] = value
	window.Invalidate()
}

func BindSlice[T any](hdata *[]T) tBindSlice[T] {
	if hdata == nil {
		return tBindSlice[T]{}
	}
	return tBindSlice[T]{data: *hdata}
}
