package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Result struct {
	Program Program
	Score   float64
}

type ByScore []Result

func (b ByScore) Len() int           { return len(b) }
func (b ByScore) Swap(i int, j int)  { b[i], b[j] = b[j], b[i] }
func (b ByScore) Less(i, j int) bool { return b[i].Score < b[j].Score }

// Returns 'a' as interpreted as a value
// Register -> Register value
// Memory location -> Memory value
func getValue(a Argument, r0 int, r1 int, r2 int, mem []int) int {
	if a == R0 {
		return r0
	} else if a == R1 {
		return r1
	} else if a == R2 {
		return r2
	} else if mem != nil && a < ARRAY_SIZE && a >= 0 {
		return mem[a]
	}
	if DEBUG {
		assert(false, "Invalid argument: "+a.Pretty())
	}
	return 0
}

func getMemValue(a Argument, r0 int, r1 int, r2 int, mem []int) int {
	if a == R0 {
		if r0 >= 0 && r0 < len(mem) {
			return mem[r0]
		}
	} else if a == R1 {
		if r1 >= 0 && r1 < len(mem) {
			return mem[r1]
		}
	} else if a == R2 {
		if r2 >= 0 && r2 < len(mem) {
			return mem[r2]
		}
	} else {
		loc := int(a)
		if loc >= 0 && loc < len(mem) {
			return mem[int(a)]
		}
	}
	return 0
}

// Returns 'a' as intepreted as a memory location
func getMemLocation(a Argument, r0 int, r1 int, r2 int) int {
	if a == R0 {
		return r0
	} else if a == R1 {
		return r1
	} else if a == R2 {
		return r2
	} else {
		return int(a)
	}
}

func setRegister(a Argument, r0 *int, r1 *int, r2 *int, val int) {
	if a == R0 {
		*r0 = val
	} else if a == R1 {
		*r1 = val
	} else if a == R2 {
		*r2 = val
	}
	if DEBUG {
		assert(false, "Not a register: "+a.Pretty())
	}
}

func jump(label string, program Program, pc int) int {
	for i, ins := range program {
		if ins.Type == LABEL && ins.StringArg == label {
			return i
		}
	}
	return pc
}

func run(program Program, mem []int, limit int) {
	r0, r1, r2 := 0, 0, 0
	pc := 0
	iterations := 0
	for pc < len(program) && pc >= 0 && iterations < limit {
		ins := program[pc]
		switch ins.Type {
		case READ:
			val := getValue(ins.Arg2, r0, r1, r2, mem)
			setRegister(ins.Arg1, &r0, &r1, &r2, val)
		case SET:
			setRegister(ins.Arg1, &r0, &r1, &r2, int(ins.Arg2))
		case INC:
			if DEBUG {
				assert(ins.Arg1.isRegister(), "Incorrect register argument: "+ins.Pretty())
			}
			val := getValue(ins.Arg1, r0, r1, r2, nil)
			setRegister(ins.Arg1, &r0, &r1, &r2, val+1)
		case DEC:
			if DEBUG {
				assert(ins.Arg1.isRegister(), "Incorrect register argument: "+ins.Pretty())
			}
			val := getValue(ins.Arg1, r0, r1, r2, nil)
			setRegister(ins.Arg1, &r0, &r1, &r2, val-1)
		case JUMPLESSTHAN:
			val1 := getMemValue(ins.Arg1, r0, r1, r2, mem)
			val2 := getMemValue(ins.Arg2, r0, r1, r2, mem)
			if val1 < val2 {
				pc += ins.JumpOffset
			}
		case JUMPZERO:
			val1 := getMemValue(ins.Arg1, r0, r1, r2, mem)
			if val1 == 0 {
				pc += ins.JumpOffset
			}
		case LABEL:
			// No-op
		case SWAP:
			loc1 := getMemLocation(ins.Arg1, r0, r1, r2)
			loc2 := getMemLocation(ins.Arg2, r0, r1, r2)
			if loc1 >= 0 && loc1 < len(mem) && loc2 >= 0 && loc2 < len(mem) {
				swap := mem[loc1]
				mem[loc1] = mem[loc2]
				mem[loc2] = swap
			}
		case JUMP:
			pc += ins.JumpOffset
		default:
			fmt.Println("Unsupported instruction: ", ins.Type)
		}
		pc++
		iterations++
	}
}

func testProgram(program Program, originalArray []int) (Result, []int) {
	mem := make([]int, MEM_SIZE)
	copy(mem, originalArray)
	run(program, mem, NUM_STEPS)
	score := 0.0
	for i := range originalArray {
		if i+1 < len(mem) && mem[i+1] < mem[i] {
			score -= 1.0
		}
	}
	return Result{program, score}, mem
}

func testProgramAsync(program Program, originalArray []int, out chan Result) {
	result, _ := testProgram(program, originalArray)
	out <- result
}

func testPrograms(programs []Program, originalArray []int) []Result {
	resultChannels := []chan Result{}
	for _, program := range programs {
		out := make(chan Result, 1)
		go testProgramAsync(program, originalArray, out)
		resultChannels = append(resultChannels, out)
	}
	results := []Result{}
	for _, out := range resultChannels {
		results = append(results, <-out)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

func best(results []Result) float64 {
	best := 1.0
	for _, result := range results {
		if best > 0.0 || result.Score > best {
			best = result.Score
		}
	}
	return best
}

func average(results []Result) float64 {
	average := 0.0
	for _, result := range results {
		average += result.Score
	}
	return average / float64(len(results))
}

func evolve(programs []Program, rounds int, print bool) []Program {
	array := make([]int, ARRAY_SIZE)
	if print {
		fmt.Printf("Round: ")
	}
	for i := 0; i < rounds; i++ {
		randomize(array, VALUE_SIZE)
		results := testPrograms(programs, array)
		programsToKeep := int(float64(len(results)) * (1.0 - LEARNING_RATE))
		numNewPrograms := len(programs) - programsToKeep
		newPrograms := make([]Program, len(programs))
		for i := 0; i < programsToKeep; i++ {
			newPrograms[i] = results[i].Program
		}
		for i := 0; i < numNewPrograms; i++ {
			newPrograms[programsToKeep+i] = mutate(newPrograms[rand.Intn(programsToKeep)])
		}
		copy(programs, newPrograms)
		if print {
			fmt.Printf("\rRound: %d", i+1)
		}
	}
	if print {
		fmt.Println()
	}
	return programs
}

func measure(program Program, count int) float64 {
	array := make([]int, ARRAY_SIZE)
	sum := int64(0)
	for i := 0; i < count; i++ {
		randomize(array, VALUE_SIZE)
		result, _ := testProgram(program, array)
		sum += int64(result.Score)
	}
	return float64(sum) / float64(count)
}

func measureMulti(programs []Program, count int) float64 {
	sum := 0.0
	for _, program := range programs {
		sum += measure(program, count)
	}
	return sum / float64(len(programs))
}
