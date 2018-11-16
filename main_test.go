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

func BenchmarkConcurrency(b *testing.B) {
	for n := 0; n < b.N; n++ {
		compute()
	}
}

func BenchmarkConcurrencyMutexes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyMutexes()
	}
}

func BenchmarkConcurrencyMutexesWithWorkers(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyMutexesWithWorkers()
	}
}

func BenchmarkConcurrencyChannelsWithWorkers(b *testing.B) {
	for n := 0; n < b.N; n++ {
		computeConcurrencyChannelsWithWorkers()
	}
}
