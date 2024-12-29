package main

import (
	"context"
	_ "embed"
)

var (
	//go:embed embed
	data []byte
)

// Foo struct
type Foo struct {
	i int
}

// Bar interface
type Bar interface {
	Do(ctx context.Context) error
}

// main
func main() {
	_ = 1
}
