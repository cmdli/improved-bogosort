package main

import "testing"

func BenchmarkRun(b *testing.B) {
	programs := loadPrograms("benchmark.progs")
	array := make([]int, ARRAY_SIZE)
	randomize(array, VALUE_SIZE)
	for n := 0; n < b.N; n++ {
		testProgram(programs[0], array)
	}
}
