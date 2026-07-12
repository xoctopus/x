package iterx_test

import (
	"context"
	"testing"
	"time"

	"github.com/xoctopus/x/iterx"
	"github.com/xoctopus/x/testx"
)

func TestOf(t *testing.T) {
	length := 0
	input := []int{0, 1, 2}
	for v := range iterx.Of(input) {
		testx.Expect(t, v, testx.Equal(length))
		length++
	}
	testx.Expect(t, length, testx.Equal(len(input)))
}

func TestSliceOf(t *testing.T) {
	length := 0
	input := []int{0, 1, 2}
	for i, v := range iterx.SliceOf(input) {
		testx.Expect(t, i, testx.Equal(v))
		length++
		if length == 2 {
			break
		}
	}
	testx.Expect(t, length, testx.Equal(len(input)-1))
}

func TestMapOf(t *testing.T) {
	length := 0
	input := map[int]int{0: 0, 1: 1, 2: 2}
	for k, v := range iterx.MapOf(input) {
		testx.Expect(t, k, testx.Equal(v))
		length++
		if length == 2 {
			break
		}
	}
	testx.Expect(t, length, testx.Equal(len(input)-1))
}

func TestRecv(t *testing.T) {
	c := make(chan int, 3)
	c <- 0
	c <- 1
	c <- 2
	close(c)

	length := 0
	for v := range iterx.Recv(c) {
		testx.Expect(t, v, testx.Equal(length))
		length++
		if length == 2 {
			break
		}
	}
	testx.Expect(t, length, testx.Equal(2))
}

func TestRecvContext(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		c := make(chan int, 3)
		c <- 1
		c <- 2
		c <- 3
		close(c)
		var out []int
		for v := range iterx.RecvContext(context.Background(), c) {
			out = append(out, v)
		}
		testx.Expect(t, out, testx.Equal([]int{1, 2, 3}))
	})

	t.Run("ContextCanceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		c := make(chan int)

		go func() {
			c <- 1
			c <- 2
			cancel() // Triggers context cancellation
			// Optional: try sending more, but they might not be read
			select {
			case c <- 3:
			case <-time.After(10 * time.Millisecond):
			}
			close(c)
		}()

		var out []int
		for v := range iterx.RecvContext(ctx, c) {
			out = append(out, v)
		}

		testx.Expect(t, len(out), testx.BeGte(2))
	})
}
