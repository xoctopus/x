package list_test

import (
	"fmt"
	"testing"

	"github.com/xoctopus/x/container/list"
	. "github.com/xoctopus/x/testx"
)

func Example() {
	l := list.New[any]()

	e3 := l.PushFront(3)
	e4 := l.PushBack(4)
	l.InsertAfter(5, e4)
	l.InsertBefore(2, e3)

	l2 := list.New[any]()
	l2.PushBack(1)

	l3 := list.New[any]()
	l3.PushFront(6)

	l.PushBackList(l2)
	l.PushFrontList(l3)

	head, tail := l.Front(), l.Back()
	l.MoveToFront(tail)
	l.MoveToBack(head)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}

var NilElement = BeNil[*list.Element[any]]

func TestList(t *testing.T) {
	Expect(t, 1, Equal(1))

	l := list.New[any]()
	Expect(t, l.Len(), Equal(0))
	Expect(t, l.Front(), NilElement())
	Expect(t, l.Back(), NilElement())

	e1 := l.PushFront(1)
	Expect(t, l.Len(), Equal(1))
	Expect(t, e1.Prev(), NilElement())
	Expect(t, e1.Next(), NilElement())

	e2 := l.PushBack(2)
	Expect(t, l.Len(), Equal(2))
	Expect(t, e2.Prev(), Equal(e1))
	Expect(t, e1.Next(), Equal(e2))

	l.MoveToFront(e2)
	Expect(t, l.Front(), Equal(e2))
	l.MoveToBack(e1)
	Expect(t, l.Back(), Equal(e1))

	wild := &list.Element[any]{}
	Expect(t, l.InsertAfter(0, wild), NilElement())
	Expect(t, l.InsertBefore(0, wild), NilElement())

	l.MoveToFront(wild)
	Expect(t, l.Front(), Equal(e2))
	l.MoveToBack(wild)
	Expect(t, l.Back(), Equal(e1))
	l.MoveBefore(wild, e1)
	Expect(t, l.Front(), Equal(e2))
	l.MoveAfter(wild, e2)
	Expect(t, l.Back(), Equal(e1))

	l.MoveBefore(e1, e2)
	Expect(t, l.Front(), Equal(e1))
	Expect(t, l.Back(), Equal(e2))
	l.MoveAfter(e2, e1)
	Expect(t, l.Front(), Equal(e1))
	Expect(t, l.Back(), Equal(e2))

	l.MoveBefore(e1, e1)
	Expect(t, l.Front(), Equal(e1))
	Expect(t, l.Back(), Equal(e2))
	l.MoveBefore(e2, e2)
	Expect(t, l.Front(), Equal(e1))
	Expect(t, l.Back(), Equal(e2))

	l.Remove(e2)
	Expect(t, l.Len(), Equal(1))
	l.Clear()
	Expect(t, l.Len(), Equal(0))
}
