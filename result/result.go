package result

type Result[T any] interface {
	Catch(func(error)) Result[T]
	Do(func(T)) Result[T]
}

type success[T any] struct {
	v T
}

func (r success[T]) Catch(func(err error)) Result[T] {
	return r
}

func (r success[T]) Do(f func(v T)) Result[T] {
	f(r.v)
	return r
}

type fail[T any] struct {
	err error
}

func (r fail[T]) Catch(f func(err error)) Result[T] {
	f(r.err)
	return r
}

func (r fail[T]) Do(f func(v T)) Result[T] {
	return r
}

func Wrap[T any](f func() (T, error)) func() Result[T] {
	return func() Result[T] {
		v, err := f()
		if err != nil {
			return fail[T]{err}
		}
		return success[T]{v: v}
	}
}
