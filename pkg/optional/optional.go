package optional

type Optional[T any] struct {
	value T
	set   bool
}

func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, set: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{}
}

func (o Optional[T]) IsSet() bool { return o.set }

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.set
}
