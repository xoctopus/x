package stringsx_test

import (
	"testing"

	"github.com/xoctopus/x/stringsx"
)

var randoms = map[string]func(int) string{
	"N": stringsx.Random,
	"X": stringsx.RandomN,
}

func TestRandomN(t *testing.T) {
	for _, f := range randoms {
		for range 10 {
			t.Log(f(10))
		}
	}
}

func BenchmarkRandom(b *testing.B) {
	for name, f := range randoms {
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				f(10)
			}
		})
	}
}
