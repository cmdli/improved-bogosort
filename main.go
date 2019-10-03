package main

import "fmt"
import "math/rand"
import "time"
import "strconv"
import "math"
import "sort"
import "os"

type InstructionType int

const (
	SET InstructionType = 0
	READ = 1
	WRITE = 2
	COMPARE = 3
	JUMPLESSTHAN = 4
	JUMPEQUAL = 5
	LABEL = 6
)

type Argument int32

const (
	R0 Argument = -3
	R1 Argument = -2
	R2 Argument = -1
)

type Instruction struct {
	Type InstructionType
	Arg1 Argument
	Arg2 Argument
	StringArg string
}

func (ins Instruction) String() string {
	return fmt.Sprintf("%#v", ins)
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
	} else if a >= 0 && mem != nil {
		return mem[a]
	} else {
		assert(false, "Invalid argument: " + string(a))
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
	lessThan, equal := false, false
	iterations := 0
	for {
		if pc >= len(program) || pc < 0 || iterations > limit{
			break
		}
		iterations++
		ins := program[pc]
		switch ins.Type {
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
		case READ:
			assert(ins.Arg2 >= 0, "Incorrect memory address: "+ins.String())
			if ins.Arg1 == R0 {
				r0 = mem[ins.Arg2]
			} else if ins.Arg1 == R1 {
				r1 = mem[ins.Arg2]
			} else if ins.Arg1 == R2 {
				r2 = mem[ins.Arg2]
			} else {
				assert(false, "Incorrect target: "+ins.String())
			}
		case WRITE:
			assert(ins.Arg2 >= 0, "Incorrect memory address: "+ins.String())
			mem[ins.Arg2] = decodeArgument(ins.Arg1, r0, r1, r2, nil)
		case COMPARE:
			val1 := decodeArgument(ins.Arg1, r0, r1, r2, nil)
			val2 := decodeArgument(ins.Arg2, r0, r1, r2, nil)
			lessThan, equal = false, false
			if val1 < val2 {
				lessThan = true
			} else if val1 == val2 {
				equal = true
			}
		case JUMPLESSTHAN:
			if lessThan {
				pc = jump(ins.StringArg, program, pc)
			}
		case JUMPEQUAL:
			if equal {
				pc = jump(ins.StringArg, program, pc)
			}
		case LABEL:
			// No-op
		default:
			fmt.Println("Unsupported instruction: ", ins.Type)
		}
		pc += 1
	}
}

type Result struct {
	Program []Instruction
	Score float64
}

type ByScore []Result
func (b ByScore) Len() int { return len(b) }
func (b ByScore) Swap(i int, j int) { b[i], b[j] = b[j], b[i] }
func (b ByScore) Less(i, j int) bool { return b[i].Score < b[j].Score }

func testProgram(program []Instruction, originalArray []int) (float64, []int) {
	mem := make([]int, 10000)
	numCounts := make(map[int]int)
	for i,num := range originalArray {
		mem[i] = num
		count, ok := numCounts[num]
		if ok {
			numCounts[num]= count + 1
		} else {
			numCounts[num] = 1
		}
	}
	run(program, mem, 10000)
	score := 0.0
	testCounts := make(map[int]int)
	for i := 0; i < len(originalArray); i++ {
		if mem[i+1] < mem[i] {
			score -= 1.0
		}
		count, ok := testCounts[mem[i]]
		if ok {
			testCounts[mem[i]]= count + 1
		} else {
			testCounts[mem[i]] = 1
		}
	}
	for num,count := range numCounts {
		testCount, ok := testCounts[num]
		if !ok {
			testCount = 0
		}
		score -= math.Abs(float64(count-testCount))
	}
	return score, mem
}

func randomRegister() Argument {
	switch (rand.Intn(3)) {
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

func randomLabel() string {
	return "L" + strconv.Itoa(rand.Intn(10))
}

func randomIns() Instruction {
	Type := InstructionType(rand.Intn(7))
	switch (Type) {
	case SET:
		return Instruction{Type: SET, Arg1: randomRegister(), Arg2: Argument(rand.Intn(1000))}
	case READ:
		return Instruction{Type: READ, Arg1: randomRegister(), Arg2: Argument(rand.Intn(1000))}
	case WRITE:
		return Instruction{Type: WRITE, Arg1: randomRegister(), Arg2: Argument(rand.Intn(1000))}
	case COMPARE:
		return Instruction{Type: COMPARE, Arg1: randomRegister(), Arg2: randomRegister()}
	case JUMPLESSTHAN:
		return Instruction{Type: JUMPLESSTHAN, StringArg: randomLabel()}
	case JUMPEQUAL:
		return Instruction{Type: JUMPEQUAL, StringArg: randomLabel()}
	case LABEL:
		return Instruction{Type: LABEL, StringArg: randomLabel()}
	default:
		assert(false, "Incorrect instruction type: " + string(Type))
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
		program = append(program, Instruction{Type:LABEL, StringArg:"NO-OP"})
	}
	return program
}

func evolve(program []Instruction) []Instruction{
	newProgram := make([]Instruction, len(program))
	copy(newProgram, program)
	if rand.Intn(100) < 100 {
		for i := 0; i < rand.Intn(3); i++ {
			program[rand.Intn(len(program))] = randomIns()
		}
	}
	return newProgram
}

func main() {
	args := os.Args
	numPrograms := 100
	if len(args) >= 2 {
		numPrograms,_ = strconv.Atoi(args[1])
	}
	numIterations := 1000
	if len(args) >= 3 {
		numIterations,_ = strconv.Atoi(args[2])
	}
	learningRate := 0.5
	if len(args) >= 4 {
		learningRate,_ = strconv.ParseFloat(args[3], 64)
	}
	randomize := true
	if len(args) >= 5 {
		randomize = args[4] != "false"
	}
	programLength := 100

	rand.Seed(time.Now().UnixNano())
	programs := [][]Instruction{}

	for i := 0; i < numPrograms; i++ {
		programs = append(programs, randomProgram(programLength))
	}
	programs[0] = nullProgram(programLength)

	originalArray := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		originalArray[i] = rand.Intn(10000)
	}
	results := []Result{}
	for _, prog := range programs {
		score, _ := testProgram(prog, originalArray)
		results = append(results, Result{prog, score})
	}
	sort.Slice(results, func (i, j int) bool {
		return results[i].Score > results[j].Score
	})
	fmt.Println("Best before:", results[0].Score)

	for i := 0; i < numIterations; i++ {
		if randomize {
			for i := 0; i < 1000; i++ {
				originalArray[i] = rand.Intn(10000)
			}
		}
		results := []Result{}
		for _, prog := range programs {
			score, _ := testProgram(prog, originalArray)
			results = append(results, Result{prog, score})
		}
		sort.Slice(results, func (i, j int) bool {
			return results[i].Score > results[j].Score
		})
		programsToKeep := int(float64(len(results))*(1.0-learningRate))
		newPrograms := make([][]Instruction, len(programs))
		for i := 0; i < programsToKeep; i++ {
			newPrograms[i] = results[i].Program
		}
		for i := programsToKeep; i < len(programs); i++ {
			newPrograms[i] = evolve(newPrograms[rand.Intn(programsToKeep)])
		}
		programs = newPrograms
	}

	if randomize {
		for i := 0; i < 1000; i++ {
			originalArray[i] = rand.Intn(10000)
		}
	}
	score, mem := testProgram(programs[0], originalArray)
	fmt.Println("Best: ", score)
	fmt.Println("Mem:",mem)
	fmt.Println("Program:",programs[0])
}
