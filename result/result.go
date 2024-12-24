package result

type Result[T any] interface {
	Catch(func(error)) Result[T]
	Do(func(T)) Result[T]
	IsSuccess() bool
	Value() T
	Error() error
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

func (success[T]) IsSuccess() bool {
	return true
}

func (r success[T]) Value() T {
	return r.v
}

func (success[T]) Error() error {
	return nil
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

func (fail[T]) IsSuccess() bool {
	return false
}

func (r fail[T]) Value() T {
	panic("Get on fail[T]")
}

func (r fail[T]) Error() error {
	return r.err
}

func Wrap[T any](f func() (T, error)) func() Result[T] {
	return func() Result[T] {
		v, err := f()
		return newResult(v, err)
	}
}

func Wrap1[T, U any](f func(T) (U, error)) func(T) Result[U] {
	return func(in T) Result[U] {
		v, err := f(in)
		return newResult(v, err)
	}
}

func newResult[T any](v T, err error) Result[T] {
	if err != nil {
		return fail[T]{err}
	}
	return success[T]{v: v}
}

func FlatMap[T, U any](f func(T) (U, error)) func(Result[T]) Result[U] {
	return func(r Result[T]) Result[U] {
		if r.IsSuccess() {
			return Wrap1(f)(r.Value())
		}
		return fail[U]{r.Error()}
	}
}
