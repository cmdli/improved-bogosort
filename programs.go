package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
)

var INSTRUCTION_SET = []InstructionType{JUMPLESSTHAN, SWAP, LABEL}

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
		return Instruction{Type: JUMPLESSTHAN, StringArg: randomLabel(), Arg1: randomMemLocation(), Arg2: randomMemLocation()}
	case JUMPZERO:
		return Instruction{Type: JUMPZERO, StringArg: randomLabel(), Arg1: randomMemLocation()}
	case LABEL:
		return Instruction{Type: LABEL, StringArg: randomLabel()}
	case SWAP:
		return Instruction{Type: SWAP, Arg1: randomMemLocation(), Arg2: randomMemLocation()}
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

func mutate(program []Instruction) []Instruction {
	newProgram := make([]Instruction, len(program))
	copy(newProgram, program)
	for i := 0; i < int(float64(len(program))*MUTATION_RATE); i++ {
		program[rand.Intn(len(program))] = randomIns()
	}
	return newProgram
}
