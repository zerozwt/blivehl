package utils

type Future[T any] struct {
	wait chan struct{}
	ret  *T
	err  error
}

func MakeFuture[T any]() *Future[T] {
	return &Future[T]{
		wait: make(chan struct{}),
		ret:  nil,
		err:  nil,
	}
}

func (f *Future[T]) Wait() (*T, error) {
	<-f.wait
	return f.ret, f.err
}

func (f *Future[T]) Done(value *T, err error) {
	f.ret = value
	f.err = err
	close(f.wait)
}
