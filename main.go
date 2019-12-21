package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const DEBUG = false
const LEARNING_RATE = 0.1
const MUTATION_RATE = 0.1
const PROGRAM_LENGTH = 100
const ARRAY_SIZE = 100
const VALUE_SIZE = 10000
const NUM_STEPS = 10000
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

func sortProgram() Program {
	return []Instruction{
		Instruction{Type: SET, Arg1: R0, Arg2: 0},
		Instruction{Type: SET, Arg1: R1, Arg2: 1},
		Instruction{Type: SET, Arg1: R2, Arg2: 100},
		Instruction{Type: JUMPLESSTHAN, Arg1: R0, Arg2: R1, JumpOffset: 1},
		Instruction{Type: SWAP, Arg1: R0, Arg2: R1},
		Instruction{Type: INC, Arg1: R0},
		Instruction{Type: INC, Arg1: R1},
		Instruction{Type: DEC, Arg1: R2},
		Instruction{Type: JUMPZERO, Arg1: R2, JumpOffset: -9},
		Instruction{Type: JUMP, JumpOffset: -7},
	}
}

func swap() Program {
	return []Instruction{
		Instruction{Type: SET, Arg1: R0, Arg2: 0},
		Instruction{Type: SET, Arg1: R1, Arg2: 1},
		Instruction{Type: SWAP, Arg1: R0, Arg2: R1},
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	command := flag.String("cmd", "", "Command to run")
	numPrograms := flag.Int("num_programs", -1, "Number of programs to generate")
	programFile := flag.String("programs", "", "Programs to load")
	numIterations := flag.Int("iterations", 100, "Number of iterations to run")
	index := flag.Int("index", -1, "Index into the program list")
	null := flag.Bool("null", false, "Generate null programs")
	length := flag.Int("length", PROGRAM_LENGTH, "Program length")
	flag.Parse()
	if *command == "generate" {
		assert(*numPrograms >= 0, "Need number of programs, see -h for help")
		programs := []Program{}
		for i := 0; i < *numPrograms; i++ {
			if *null {
				programs = append(programs, nullProgram(*length))
			} else {
				programs = append(programs, randomProgram(*length))
			}
		}
		programs[0] = nullProgram(PROGRAM_LENGTH)
		writePrograms(*programFile, programs)
	} else if *command == "test" {
		assert(*programFile != "", "Need program file, see -h for help")
		programs := loadPrograms(*programFile)
		if *index >= 0 {
			fmt.Println("Score:", measure(programs[*index], 1000))
		} else {
			fmt.Println("Average score:", measureMulti(programs, 100))
		}
	} else if *command == "evolve" {
		assert(*programFile != "", "Need program file, see -h for help")
		programs := loadPrograms(*programFile)
		fmt.Println("Before:", measureMulti(programs, 10))
		evolve(programs, *numIterations, false)
		fmt.Println("After:", measureMulti(programs, 10))
		writePrograms(*programFile, programs)
	} else if *command == "print" {
		assert(*programFile != "", "Need program file, see -h for help")
		assert(*index >= 0, "Need index, see -h for help")
		programs := loadPrograms(*programFile)
		fmt.Println(Program(programs[*index]).Pretty())
	} else if *command == "test_sort" {
		program := sortProgram()
		array := make([]int, ARRAY_SIZE)
		randomize(array, VALUE_SIZE)
		fmt.Println("Before:", array)
		result, after := testProgram(program, array)
		fmt.Println("After:", after)
		fmt.Println("Score:", result.Score)
	} else if *command == "random_prog" {
		fmt.Println(randomProgram(PROGRAM_LENGTH).Pretty())
	} else {
		fmt.Println("Unrecognized command ", *command)
	}
}
