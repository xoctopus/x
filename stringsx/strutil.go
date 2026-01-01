package stringsx

import (
	randx "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func Random(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func RandomN(n int) string {
	b := make([]byte, n)
	for i := 0; i < len(b); {
		if num, err := randx.Int(randx.Reader, big.NewInt(int64(len(letters)))); err == nil {
			b[i] = letters[num.Int64()]
			i++
		}
	}
	return string(b)
}
