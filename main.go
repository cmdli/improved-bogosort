package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

const LEARNING_RATE = 0.3
const MUTATION_RATE = 0.1
const PROGRAM_LENGTH = 100
const ARRAY_SIZE = 100
const VALUE_SIZE = 10000
const NUM_STEPS = 1000
const MEM_SIZE = 100

var INSTRUCTION_SET = []InstructionType{JUMPLESSTHAN, JUMPZERO, SWAP, LABEL}

type InstructionType int

const (
	SET          InstructionType = 0
	INC                          = 1
	DEC                          = 2
	JUMPLESSTHAN                 = 3
	JUMPZERO                     = 4
	SWAP                         = 5
	LABEL                        = 6
	READ                         = 7
)

type Argument int32

const (
	R0 Argument = -3
	R1 Argument = -2
	R2 Argument = -1
)

func (arg Argument) isRegister() bool {
	return arg == R0 || arg == R1 || arg == R2
}

func (arg Argument) Pretty() string {
	if arg == R0 {
		return "R0"
	} else if arg == R1 {
		return "R1"
	} else if arg == R2 {
		return "R2"
	}
	return strconv.Itoa(int(arg))
}

type Instruction struct {
	Type      InstructionType
	Arg1      Argument
	Arg2      Argument
	StringArg string
}

func (ins Instruction) String() string {
	return fmt.Sprintf("%#v", ins)
}

func (ins Instruction) Pretty() string {
	switch ins.Type {
	case SET:
		return "SET " + ins.Arg1.Pretty() + " " + ins.Arg2.Pretty()
	case INC:
		return "INC " + ins.Arg1.Pretty()
	case DEC:
		return "DEC " + ins.Arg1.Pretty()
	case JUMPLESSTHAN:
		return "JLT " + ins.Arg1.Pretty() + " " + ins.Arg2.Pretty() + " " + ins.StringArg
	case JUMPZERO:
		return "JZ " + ins.Arg1.Pretty() + " " + ins.StringArg
	case SWAP:
		return "SWAP " + ins.Arg1.Pretty() + " " + ins.Arg2.Pretty()
	case LABEL:
		return "LABEL " + ins.StringArg + ":"
	default:
		assert(false, "Incorrect instruction type: "+string(ins.Type))
	}
	return "INVALID"
}

type Program []Instruction

func (program Program) Pretty() string {
	out := ""
	for _, ins := range program {
		out += ins.Pretty() + "\n"
	}
	return out
}

func assert(b bool, err string) {
	if !b {
		fmt.Println(err)
		panic(err)
	}
}

func decodeArgument(a Argument, r0 int, r1 int, r2 int, mem []int) int {
	if a == R0 {
		return r0
	} else if a == R1 {
		return r1
	} else if a == R2 {
		return r2
	} else if a < ARRAY_SIZE && a >= 0 {
		return mem[a]
	} else {
		assert(false, "Invalid argument: "+a.Pretty())
		return 0
	}
}

func jump(label string, program []Instruction, pc int) int {
	for i, ins := range program {
		if ins.Type == LABEL && ins.StringArg == label {
			return i
		}
	}
	return pc
}

func run(program []Instruction, mem []int, limit int) {
	r0, r1, r2 := 0, 0, 0
	pc := 0
	iterations := 0
	for {
		if pc >= len(program) || pc < 0 || iterations > limit {
			break
		}
		iterations++
		ins := program[pc]
		switch ins.Type {
		case READ:
			val := decodeArgument(ins.Arg2, r0, r1, r2, nil)
			if ins.Arg1 == R0 {
				r0 = val
			} else if ins.Arg1 == R1 {
				r1 = val
			} else if ins.Arg1 == R2 {
				r2 = val
			} else {
				assert(false, "Incorrect target: "+ins.String())
			}
		case SET:
			if ins.Arg1 == R0 {
				r0 = int(ins.Arg2)
			} else if ins.Arg1 == R1 {
				r1 = int(ins.Arg2)
			} else if ins.Arg1 == R2 {
				r2 = int(ins.Arg2)
			} else {
				assert(false, "Incorrect target: "+ins.String())
			}
		case INC:
			assert(ins.Arg1.isRegister(), "Incorrect register argument: "+ins.Pretty())
			if ins.Arg1 == R0 {
				r0++
			} else if ins.Arg1 == R1 {
				r1++
			} else if ins.Arg1 == R2 {
				r2++
			}
		case DEC:
			assert(ins.Arg1.isRegister(), "Incorrect register argument: "+ins.Pretty())
			if ins.Arg1 == R0 {
				r0--
			} else if ins.Arg1 == R1 {
				r1--
			} else if ins.Arg1 == R2 {
				r2--
			}
		case JUMPLESSTHAN:
			val1 := decodeArgument(ins.Arg1, r0, r1, r2, mem)
			val2 := decodeArgument(ins.Arg2, r0, r1, r2, mem)
			if val1 < val2 {
				pc = jump(ins.StringArg, program, pc)
			}
		case JUMPZERO:
			val1 := decodeArgument(ins.Arg1, r0, r1, r2, mem)
			if val1 == 0 {
				pc = jump(ins.StringArg, program, pc)
			}
		case LABEL:
			// No-op
		case SWAP:
			assert(ins.Arg1.isRegister(), "Needs a register: "+ins.Pretty())
			assert(ins.Arg2.isRegister(), "Needs a register: "+ins.Pretty())
			val1 := decodeArgument(ins.Arg1, r0, r1, r2, nil)
			val2 := decodeArgument(ins.Arg2, r0, r1, r2, nil)
			if val1 < 0 || val1 >= len(mem) || val2 < 0 || val2 >= len(mem) {
				return
			}
			swap := mem[val1]
			mem[val1] = mem[val2]
			mem[val2] = swap
		default:
			fmt.Println("Unsupported instruction: ", ins.Type)
		}
		pc++
	}
}

type Result struct {
	Program []Instruction
	Score   float64
}

type ByScore []Result

func (b ByScore) Len() int           { return len(b) }
func (b ByScore) Swap(i int, j int)  { b[i], b[j] = b[j], b[i] }
func (b ByScore) Less(i, j int) bool { return b[i].Score < b[j].Score }

func testProgram(program []Instruction, originalArray []int) (float64, []int) {
	mem := make([]int, MEM_SIZE)
	copy(mem, originalArray)
	run(program, mem, NUM_STEPS)
	score := 0.0
	for i := range originalArray {
		if i+1 < len(mem) && mem[i+1] < mem[i] {
			score -= 1.0
		}
	}
	return score, mem
}

func testPrograms(programs [][]Instruction, originalArray []int) []Result {
	results := []Result{}
	for _, prog := range programs {
		score, _ := testProgram(prog, originalArray)
		results = append(results, Result{prog, score})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

func randomMemLocation() Argument {
	return Argument(rand.Intn(ARRAY_SIZE))
}

func randomRegister() Argument {
	switch rand.Intn(3) {
	case 0:
		return R0
	case 1:
		return R1
	case 2:
		return R2
	}
	assert(false, "Invalid code in randomRegister")
	return R0
}

func randomArgument() Argument {
	switch rand.Intn(4) {
	case 0:
		return R0
	case 1:
		return R1
	case 2:
		return R2
	case 3:
		return randomMemLocation()
	}
	assert(false, "Invalid code in randomRegister")
	return R0
}

func randomLabel() string {
	return "L" + strconv.Itoa(rand.Intn(10))
}

func choice(ins []InstructionType) InstructionType {
	return ins[rand.Intn(len(ins))]
}

func randomIns() Instruction {
	insType := choice(INSTRUCTION_SET)
	switch insType {
	case READ:
		return Instruction{Type: SET, Arg1: randomMemLocation(), Arg2: randomRegister()}
	case SET:
		return Instruction{Type: SET, Arg1: randomRegister(), Arg2: Argument(rand.Intn(ARRAY_SIZE))}
	case INC:
		return Instruction{Type: INC, Arg1: randomRegister()}
	case DEC:
		return Instruction{Type: DEC, Arg1: randomRegister()}
	case JUMPLESSTHAN:
		return Instruction{Type: JUMPLESSTHAN, StringArg: randomLabel(), Arg1: randomArgument(), Arg2: randomRegister()}
	case JUMPZERO:
		return Instruction{Type: JUMPZERO, StringArg: randomLabel(), Arg1: randomArgument()}
	case LABEL:
		return Instruction{Type: LABEL, StringArg: randomLabel()}
	case SWAP:
		return Instruction{Type: SWAP, Arg1: randomRegister(), Arg2: randomRegister()}
	default:
		assert(false, "Incorrect instruction type: "+string(insType))
	}
	return Instruction{Type: LABEL, StringArg: "NoOp"}
}

func randomProgram(length int) []Instruction {
	prog := []Instruction{}
	for i := 0; i < length; i++ {
		prog = append(prog, randomIns())
	}
	return prog
}

func nullProgram(length int) []Instruction {
	program := []Instruction{}
	for i := 0; i < length; i++ {
		program = append(program, Instruction{Type: LABEL, StringArg: "NO-OP"})
	}
	return program
}

func mutate(program []Instruction) []Instruction {
	newProgram := make([]Instruction, len(program))
	copy(newProgram, program)
	for i := 0; i < int(float64(len(program))*MUTATION_RATE); i++ {
		program[rand.Intn(len(program))] = randomIns()
	}
	return newProgram
}

func evolve(programs [][]Instruction, rounds int) [][]Instruction {
	array := make([]int, ARRAY_SIZE)
	fmt.Printf("Round: ")
	for i := 0; i < rounds; i++ {
		randomize(array, VALUE_SIZE)
		results := testPrograms(programs, array)
		programsToKeep := int(float64(len(results)) * (1.0 - LEARNING_RATE))
		numNewPrograms := len(programs) - programsToKeep
		newPrograms := make([][]Instruction, len(programs))
		for i := 0; i < programsToKeep; i++ {
			newPrograms[i] = results[i].Program
		}
		fmt.Println("To keep", programsToKeep)
		fmt.Println("Num new programs", numNewPrograms)
		for i := 0; i < numNewPrograms/2; i++ {
			newPrograms[programsToKeep+i] = mutate(newPrograms[rand.Intn(programsToKeep)])
		}
		for i := numNewPrograms / 2; i < numNewPrograms; i++ {
			newPrograms[programsToKeep+i] = randomProgram(len(programs[0]))
		}
		programs = newPrograms
		fmt.Printf("\rRound: %d", i+1)
	}
	fmt.Println()
	return programs
}

func randomize(array []int, numberRange int) {
	for i := 0; i < len(array); i++ {
		array[i] = rand.Intn(numberRange)
	}
}

func loadPrograms(filename string) [][]Instruction {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	programs := [][]Instruction{}
	err = json.Unmarshal(input, &programs)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	return programs
}

func writePrograms(filename string, programs [][]Instruction) {
	output, err := json.Marshal(programs)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, output, 0644)
	if err != nil {
		log.Fatal(err)
	}
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

func main() {
	rand.Seed(time.Now().UnixNano())
	args := os.Args
	if args[1] == "generate" {

		programs := [][]Instruction{}
		numPrograms, _ := strconv.Atoi(args[3])
		for i := 0; i < numPrograms; i++ {
			programs = append(programs, randomProgram(PROGRAM_LENGTH))
		}
		programs[0] = nullProgram(PROGRAM_LENGTH)
		writePrograms(args[2], programs)
	} else if args[1] == "test" {
		programs := loadPrograms(args[2])
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
	} else if args[1] == "evolve" {
		programs := loadPrograms(args[2])
		numIterations := 100
		if len(args) >= 3 {
			numIterations, _ = strconv.Atoi(args[3])
		}
		originalArray := make([]int, ARRAY_SIZE)
		randomize(originalArray, VALUE_SIZE)
		results := testPrograms(programs, originalArray)
		fmt.Println("Average before:", average(results), "Best before:", best(results))
		evolve(programs, numIterations)
		results = testPrograms(programs, originalArray)
		fmt.Println("Average after:", average(results), "Best after:", best(results))
		writePrograms(args[2], programs)
	} else if args[1] == "print" {
		programs := loadPrograms(args[2])
		index, _ := strconv.Atoi(args[3])
		fmt.Println(Program(programs[index]).Pretty())
	}
}
