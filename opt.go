package logx

type optval[T any] struct {
	value T
	isset bool
}

func (o *optval[T]) has() bool { return o.isset }
func (o *optval[T]) set(v T)   { o.value = v; o.isset = true }
func (o *optval[T]) get() T    { return o.value }
