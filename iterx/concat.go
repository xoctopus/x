package iterx

import (
	"iter"
	"sync"
)

func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	if len(seqs) == 0 {
		return func(yield func(T) bool) {}
	}

	if len(seqs) == 1 {
		return seqs[0]
	}

	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for x := range seq {
				if !yield(x) {
					return
				}
			}
		}
	}
}

func Merge[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	if len(seqs) == 0 {
		return func(yield func(T) bool) {}
	}

	if len(seqs) == 1 {
		return seqs[0]
	}

	return func(yield func(T) bool) {
		var (
			values = make(chan T)
			done   = make(chan struct{})
			wg     = &sync.WaitGroup{}
		)

		for _, seq := range seqs {
			wg.Go(func() {
				for v := range seq {
					values <- v
				}
			})
		}

		go func() {
			wg.Wait()
			close(done)
		}()

		for {
			select {
			case <-done:
				return
			case v := <-values:
				if !yield(v) {
					return
				}

			}
		}
	}
}
