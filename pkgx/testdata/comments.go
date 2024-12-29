// package
package main

import (
	"fmt"
	"time"
)

// Date defines corresponding time.Time
type Date time.Time

// const vars
const (
	// A comment
	A = "A" // A inline comment
	// B comment
	B = "B" // B inline comment
	// C comment
	C = "C" // C inline comment
	// placeholder
	_ = "D"
)

// Test struct
type Test struct {
	// field String
	String string
	// field Int
	Int int
	// field Bool
	Bool bool
	// field Date
	Date Date
}

// Recv recv method
func (Test) Recv() {}

// type group
type (
	// Test2 struct
	Test2 struct {
		// field String
		String string
		// field Int
		Int int
		// field Bool
		Bool bool
	}
)

// var
var test = Test{
	String: "",
	Int:    1 + 1,
	Bool:   true,
}

// var group
var (
	// test2
	test2 = Test{
		String: "",
		Int:    1,
		Bool:   true,
	}
	// test3
	test3 = Test{
		String: "",
		Int:    1,
		Bool:   true,
	}
)

// Print function
//
//go:generate echo
func Print(a string, b string) string {
	return a + b
}

// func fn
func fn() {
	// Call
	res := Print("", "")
	if res != "" {
		// print
		fmt.Println(res)
	}
}
