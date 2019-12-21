package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
)

var INSTRUCTION_SET = []InstructionType{SET, JUMPLESSTHAN, SWAP, INC, DEC, JUMP, JUMPZERO}

type InstructionType int

/*
Example
SET R0 0
SET R1 1
JUMPLESSTHAN R0 R1 L0
SWAP R0 R1
LABEL L0
*/
const (
	SET          InstructionType = 0
	INC                          = 1
	DEC                          = 2
	JUMPLESSTHAN                 = 3
	JUMPZERO                     = 4
	SWAP                         = 5
	LABEL                        = 6
	READ                         = 7
	JUMP                         = 8
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
	Type       InstructionType
	Arg1       Argument
	Arg2       Argument
	JumpOffset int
	StringArg  string
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
		return "JUMPLESSTHAN " + ins.Arg1.Pretty() + " " + ins.Arg2.Pretty() + " " + strconv.Itoa(ins.JumpOffset)
	case JUMPZERO:
		return "JUMPZERO " + ins.Arg1.Pretty() + " " + strconv.Itoa(ins.JumpOffset)
	case SWAP:
		return "SWAP " + ins.Arg1.Pretty() + " " + ins.Arg2.Pretty()
	case LABEL:
		return "LABEL " + ins.StringArg + ":"
	case JUMP:
		return "JUMP " + strconv.Itoa(ins.JumpOffset)
	default:
		assert(false, "Incorrect instruction type: "+strconv.Itoa(int(ins.Type)))
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

func randomArgument(excluding Argument) Argument {
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

func randomJumpOffset() int {
	offset := rand.Intn(20) - 10
	if offset == -1 {
		offset = 0
	}
	return offset
}

func choice(ins []InstructionType) InstructionType {
	return ins[rand.Intn(len(ins))]
}

func shuffle(arr []Argument) {
	for i := range arr {
		dst := rand.Intn(len(arr)-i) + i
		swap := arr[dst]
		arr[dst] = arr[i]
		arr[i] = swap
	}
}

func randomIns() Instruction {
	registers := []Argument{R0, R1, R2}
	shuffle(registers)
	insType := choice(INSTRUCTION_SET)
	switch insType {
	case READ:
		return Instruction{Type: SET, Arg1: registers[0], Arg2: randomMemLocation()}
	case SET:
		return Instruction{Type: SET, Arg1: registers[0], Arg2: randomMemLocation()}
	case INC:
		return Instruction{Type: INC, Arg1: registers[0]}
	case DEC:
		return Instruction{Type: DEC, Arg1: registers[0]}
	case JUMPLESSTHAN:
		return Instruction{Type: JUMPLESSTHAN, JumpOffset: randomJumpOffset(), Arg1: registers[0], Arg2: registers[1]}
	case JUMPZERO:
		return Instruction{Type: JUMPZERO, JumpOffset: randomJumpOffset(), Arg1: registers[0]}
	case LABEL:
		return Instruction{Type: LABEL, StringArg: randomLabel()}
	case SWAP:
		return Instruction{Type: SWAP, Arg1: registers[0], Arg2: registers[1]}
	case JUMP:
		return Instruction{Type: JUMP, JumpOffset: randomJumpOffset()}
	default:
		assert(false, "Incorrect instruction type: "+string(insType))
	}
	return Instruction{Type: LABEL, StringArg: "NoOp"}
}

func randomProgram(length int) Program {
	prog := []Instruction{}
	for i := 0; i < length; i++ {
		prog = append(prog, randomIns())
	}
	return prog
}

func nullProgram(length int) Program {
	program := []Instruction{}
	for i := 0; i < length; i++ {
		program = append(program, Instruction{Type: LABEL, StringArg: "NO-OP"})
	}
	return program
}

func loadPrograms(filename string) []Program {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	programs := []Program{}
	err = json.Unmarshal(input, &programs)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	return programs
}

func writePrograms(filename string, programs []Program) {
	output, err := json.Marshal(programs)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, output, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func mutate(program Program) Program {
	newProgram := make([]Instruction, len(program))
	copy(newProgram, program)
	for i := 0; i < int(float64(len(program))*MUTATION_RATE); i++ {
		program[rand.Intn(len(program))] = randomIns()
	}
	return newProgram
}
