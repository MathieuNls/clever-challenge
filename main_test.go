package main

import (
	"testing"
)

func BenchmarkNoConcurrency(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeNoConcurrency()
	}
}

func BenchmarkConcurrencyReadingOnly(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyReadingOnly()
	}
}

func BenchmarkConcurrencyChannelsOneGoroutinePerLine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyChannelsOneGoroutinePerLine()
	}
}

func BenchmarkConcurrencyMutexesOneGoroutinePerLine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyMutexesOneGoroutinePerLine()
	}
}

func BenchmarkConcurrencyMutexesOneGoroutinePerCPU(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyMutexesOneGoroutinePerCPU()
	}
}

func BenchmarkConcurrencyChannelsOneGoroutinePerCPU(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyChannelsOneGoroutinePerCPU()
	}
}
