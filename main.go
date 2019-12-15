package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const LEARNING_RATE = 0.3
const MUTATION_RATE = 0.1
const PROGRAM_LENGTH = 1000
const ARRAY_SIZE = 100
const VALUE_SIZE = 10000
const NUM_STEPS = 1000
const MEM_SIZE = 100

func assert(b bool, err string) {
	if !b {
		fmt.Println(err)
		panic(err)
	}
}

func randomize(array []int, numberRange int) {
	for i := 0; i < len(array); i++ {
		array[i] = rand.Intn(numberRange)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	args := os.Args
	command := flag.String("cmd", "", "Command to run")
	numPrograms := flag.Int("num_programs", -1, "Number of programs to generate")
	programFile := flag.String("programs", "", "Programs to load")
	numIterations := flag.Int("iterations", 100, "Number of iterations to run")
	index := flag.Int("index", -1, "Index into the program list")
	flag.Parse()
	if *command == "generate" {
		assert(*numPrograms >= 0, "Need number of programs, see -h for help")
		programs := []Program{}
		for i := 0; i < *numPrograms; i++ {
			programs = append(programs, randomProgram(PROGRAM_LENGTH))
		}
		programs[0] = nullProgram(PROGRAM_LENGTH)
		writePrograms(args[2], programs)
	} else if *command == "test" {
		assert(*programFile != "", "Need program file, see -h for help")
		programs := loadPrograms(*programFile)
		array := make([]int, ARRAY_SIZE)
		sum := 0.0
		count := 0
		for i := 0; i < 100; i++ {
			randomize(array, VALUE_SIZE)
			results := testPrograms(programs, array)
			for _, result := range results {
				sum += result.Score
				count++
			}
		}
		println("Average score:", sum/float64(count))
	} else if *command == "evolve" {
		assert(*programFile != "", "Need program file, see -h for help")
		programs := loadPrograms(*programFile)
		originalArray := make([]int, ARRAY_SIZE)
		randomize(originalArray, VALUE_SIZE)
		results := testPrograms(programs, originalArray)
		fmt.Println("Average before:", average(results), "Best before:", best(results))
		evolve(programs, *numIterations)
		results = testPrograms(programs, originalArray)
		fmt.Println("Average after:", average(results), "Best after:", best(results))
		writePrograms(args[2], programs)
	} else if *command == "print" {
		assert(*programFile != "", "Need program file, see -h for help")
		assert(*index >= 0, "Need index, see -h for help")
		programs := loadPrograms(*programFile)
		fmt.Println(Program(programs[*index]).Pretty())
	}
}
