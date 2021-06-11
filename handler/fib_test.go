package handler

import (
	"testing"
)

func BenchmarkFib(b *testing.B) {
	//fmt.Printf("b.N: %v\n",b.N)
	for n := 0; n < b.N; n ++ {
		Fib(20)
	}
}
