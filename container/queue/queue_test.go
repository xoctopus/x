package queue_test

import (
	"fmt"

	"github.com/xoctopus/x/container/queue"
)

func ExampleQueue() {
	for _, q := range []queue.Queue[any]{
		queue.NewSafeQueue[any](),
		queue.NewQueue[any](),
	} {
		q.Push(1)
		q.Push(2)
		q.Push(3)
		q.Push(4)

		head, _ := q.Head()
		tail, _ := q.Tail()
		fmt.Println("head:", head)
		fmt.Println("tail:", tail)
		for {
			if v, ok := q.Pop(); ok {
				fmt.Println(v)
				continue
			}
			break
		}
		q.Clear()
		q.Len()
		_, ok := q.Head()
		fmt.Println("head:", ok)
		_, ok = q.Tail()
		fmt.Println("tail:", ok)
		_, ok = q.Pop()
		fmt.Println("pop: ", ok)
	}

	// Output:
	// head: 1
	// tail: 4
	// 1
	// 2
	// 3
	// 4
	// head: false
	// tail: false
	// pop:  false
	// head: 1
	// tail: 4
	// 1
	// 2
	// 3
	// 4
	// head: false
	// tail: false
	// pop:  false
}
