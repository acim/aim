package iter

type Iterator[T any] interface {
	HasNext() bool
	Next() *T
	Chan() <-chan *T
}

type Slice[T any] struct {
	inner []T
	idx   int
}

func FromSlice[T any](s []T) *Slice[T] {
	return &Slice[T]{inner: s, idx: -1}
}

func (i *Slice[T]) HasNext() bool {
	if i.idx < len(i.inner)-1 {
		return true
	}

	return false
}

func (i *Slice[T]) Next() *T {
	i.idx++
	return &i.inner[i.idx]
}

func (i *Slice[T]) Chan() <-chan *T {
	ch := make(chan *T)

	go func() {
		for {
			if i.idx < len(i.inner)-1 {
				i.idx++
				ch <- &i.inner[i.idx]
				continue

			}

			close(ch)
			break
		}
	}()

	return ch
}
