package iterx

import (
	"iter"
	"time"
)

// ChunkDuration splits sequence by time window and returns chunked sequence.
func ChunkDuration[V any](seq iter.Seq[V], d time.Duration) iter.Seq[[]V] {
	return func(yield func([]V) bool) {
		buffer := make([]V, 0)

		timer := time.NewTimer(d)
		defer timer.Stop()

		done := make(chan struct{})
		values := make(chan V)

		go func() {
			defer close(values)
			for v := range seq {
				values <- v
			}
			close(done)
		}()

		for {
			select {
			case val, ok := <-values:
				if !ok {
					if len(buffer) > 0 {
						yield(buffer)
					}
					return
				}
				buffer = append(buffer, val)
				timer.Reset(d)
			case <-timer.C:
				if len(buffer) > 0 {
					if !yield(buffer) {
						return
					}
					buffer = make([]V, 0)
				}
				timer.Reset(d)
			case <-done:
				if len(buffer) > 0 {
					yield(buffer)
				}
				return
			}
		}
	}
}

func Debounce[V any](seq iter.Seq[V], d time.Duration) iter.Seq[V] {
	return func(yield func(V) bool) {
		buffered := ChunkDuration(seq, d)

		for b := range buffered {
			if len(b) > 0 {
				if !yield(b[len(b)-1]) {
					return
				}
			}
		}
	}
}

// ChunkN splits a sequence into fixed-size batches
func ChunkN[V any](seq iter.Seq[V], n int) iter.Seq[[]V] {
	if n < 1 {
		panic("cannot be less than 1")
	}

	return func(yield func([]V) bool) {
		count := 0
		chunk := make([]V, 0, n)

		emit := func() bool {
			if !yield(chunk) {
				return false
			}

			// reset if seq not break
			count = 0
			chunk = make([]V, 0, n)
			return true
		}

		for e := range seq {
			chunk = append(chunk, e)
			count++

			if count == n {
				if !emit() {
					return
				}
			}
		}

		if len(chunk) > 0 {
			if !emit() {
				return
			}
		}
	}
}
