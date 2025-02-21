package main

import (
	"testing"
)

func BenchmarkSequentialProcess(b *testing.B) {
	b.ReportMetric(0, "ns/op")
	input := make([]int, 1_000_000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SequentialProcess(input)
	}
	b.ReportMetric(float64(b.Elapsed().Seconds())/float64(b.N), "sec/op")
}

func BenchmarkParallelProcess(b *testing.B) {
	b.ReportMetric(0, "ns/op")
	input := make([]int, 1_000_000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelProcess(input)
	}
	b.ReportMetric(float64(b.Elapsed().Seconds())/float64(b.N), "sec/op")
}
