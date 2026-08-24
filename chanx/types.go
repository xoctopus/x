package chanx

// Deprecated: use github.com/xoctopus/concx/pkg/chanx. Cancelable is archived here.
type Cancelable interface {
	CancelCause(err error)
}

// Deprecated: use github.com/xoctopus/concx/pkg/chanx.
type ValueNotifier[T any] interface {
	Send(x T)
}

// Deprecated: use github.com/xoctopus/concx/pkg/chanx.
type ValueObserver[T any] interface {
	Value() <-chan T
}

// Deprecated: use github.com/xoctopus/concx/pkg/chanx.
type Observable[T any] interface {
	Observe() Observer[T]
}

// Deprecated: use github.com/xoctopus/concx/pkg/chanx.
type Observer[T any] interface {
	ValueObserver[T]
	Cancelable

	Done() <-chan struct{}
	Err() error
}

// Deprecated: use github.com/xoctopus/concx/pkg/chanx.
type NotifiableObserver[T any] interface {
	Observer[T]
	ValueNotifier[T]
}

// Deprecated: use github.com/xoctopus/concx/pkg/chanx.
type Subscriber[T any] interface {
	ValueNotifier[T]
	Cancelable

	Done() <-chan struct{}
	Err() error
}

// Deprecated: use github.com/xoctopus/concx/pkg/chanx.
type ObserverFunc[T any] func() Observer[T]

func (f ObserverFunc[T]) Observe() Observer[T] {
	return f()
}
