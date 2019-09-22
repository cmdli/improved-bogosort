package main

import "fmt"

type InstructionType int

const (
	SET InstructionType = iota
	READ
	WRITE
	COMPARE
	JUMPLESSTHAN
	JUMPEQUAL
	LABEL
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

func jump(label string, program []Instruction) int {
	for i, ins := range program {
		if ins.Type == LABEL && ins.StringArg == label {
			return i
		}
	}
	assert(false, "Invalid label: " + label)
	return -1
}

func run(program []Instruction, mem []int) {
	r0, r1, r2 := 0, 0, 0
	pc := 0
	lessThan, equal := false, false
	for {
		if pc >= len(program) || pc < 0 {
			break
		}
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
			assert(ins.Arg2 > 0, "Incorrect memory address: "+ins.String())
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
			assert(ins.Arg2 > 0, "Incorrect memory address: "+ins.String())
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
				pc = jump(ins.StringArg, program)
			}
		case JUMPEQUAL:
			if equal {
				pc = jump(ins.StringArg, program)
			}
		case LABEL:
			// No-op
		default:
			fmt.Println("Unsupported instruction: ", ins.Type)
		}
		pc += 1
	}
}

func randomIns() Instruction {
	return Instruction{Type: LABEL, StringArg: "NoOp"}
}

func main() {
	prog := []Instruction{}
	prog = append(prog, Instruction{Type: SET, Arg1: R0, Arg2: 2})
	prog = append(prog, Instruction{Type: READ, Arg1: R1, Arg2: 1})
	prog = append(prog, Instruction{Type: WRITE, Arg1: R1, Arg2: 0})
	prog = append(prog, Instruction{Type: WRITE, Arg1: R0, Arg2: 5})
	mem := []int{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	run(prog, mem)
	fmt.Println(mem)
}
