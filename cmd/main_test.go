package main

import (
	"testing"
)

func BenchmarkSequentialProcess(b *testing.B) {
	input := make([]int, 1_000_000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SequentialProcess(input)
	}
}

func BenchmarkParallelProcess(b *testing.B) {
	input := make([]int, 1_000_000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelProcess(input)
	}
}
