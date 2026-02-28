package stack_test

import (
	"fmt"

	"github.com/xoctopus/x/container/stack"
)

func ExampleStack() {
	for _, s := range []stack.Stack[any]{
		stack.NewSafeStack[any](),
		stack.NewStack[any](),
	} {
		s.Push(1)
		s.Push(2)
		s.Push(3)
		s.Push(4)

		top, _ := s.Top()
		fmt.Println("top:", top)
		for {
			if v, ok := s.Pop(); ok {
				fmt.Println(v)
				continue
			}
			break
		}
		s.Clear()
		fmt.Println("len:", s.Len())
		_, ok := s.Top()
		fmt.Println("top:", ok)
		_, ok = s.Pop()
		fmt.Println("pop:", ok)
	}

	// Output:
	// top: 4
	// 4
	// 3
	// 2
	// 1
	// len: 0
	// top: false
	// pop: false
	// top: 4
	// 4
	// 3
	// 2
	// 1
	// len: 0
	// top: false
	// pop: false
}
