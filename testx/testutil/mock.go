package testutil

import (
	"fmt"
	"strings"
	"testing"
)

type MockTB struct {
	testing.TB

	output string
}

func (*MockTB) Helper() {}

func (m *MockTB) Fatal(args ...any) {
	out := fmt.Sprint(args...)
	out = strings.TrimPrefix(out, "\n")
	out = strings.ReplaceAll(out, "\u00a0", " ")

	if out[len(out)-1] != '\n' {
		out += "\n"
	}
	m.output = out
}

func (m *MockTB) Fatalf(msg string, args ...any) {
	m.Fatal(fmt.Sprintf(msg, args...))
}

func (m *MockTB) Reset() {
	m.output = ""
}

func (m *MockTB) Output2() string {
	return m.output
}
