package processor

import (
	"fmt"
	"sort"
)

type Opcode uint8

// Instruction represents a single Cpu instruction. This can be generated from the static data.
type Instruction struct {
	Opcode         Opcode
	AddressingFunc AddressingFunc
	Operation      Operation

	// The baseline number of cycles that this instruction normally requires (over and above.
	// the single cycle the CPU took to read the operation code first).
	Cycles uint

	// If true, this operation incurs an additional cycle if the addressing mode indicates a
	// page boundary has been crossed.
	PageBoundaryPenalty bool
}

type Instructions []Instruction

// Converts the Instruction instance into a canonical string form.
func (i Instruction) String() string {
	return fmt.Sprintf(
		"Opcode: $%02X, Cycles: %v, PBP: %v",
		i.Opcode, i.Cycles, i.PageBoundaryPenalty)
}

// Execute takes a starting State and returns the changed State after
// executing the instruction operation. also returned are the number
// of cycles taken to execute the operation.
func (i Instruction) Execute(state State, memory Memory) (State, uint, error) {
	if memory == nil {
		return State{}, 0, MemoryMustBeProvided
	}

	if i.AddressingFunc == nil {
		return State{}, 0, NoAddressingModeFunction
	}

	if i.Operation == nil {
		return State{}, 0, NoOperationFunction
	}

	addressingState, err := i.AddressingFunc(state, memory)
	if err != nil {
		return State{}, 0, err
	}

	cycles := i.Cycles
	if i.PageBoundaryPenalty && addressingState.PageBoundaryCrossed {
		cycles++
	}
	state.PC += addressingState.ProgramCounterChange

	state, err = i.Operation(state, addressingState)
	if err != nil {
		return State{}, 0, err
	}

	return state, cycles, nil
}

// InstructionSet is a straight forward map of opcodes to Instruction instances.
type InstructionSet struct {
	instructions map[Opcode]Instruction
}

// validate returns an error if the instruction set is empty.
func (is InstructionSet) validate() error {
	if len(is.instructions) <= 0 {
		return InstructionSetEmpty
	}
	return nil
}

// Get returns the Instruction represented by the opcode from InstructionSet.
// An error is returned if the opcode does not exist in the Instruction Set.
func (is InstructionSet) Get(opcode Opcode) (Instruction, error) {
	if err := is.validate(); err != nil {
		return Instruction{Opcode: opcode}, err
	}

	if val, ok := is.instructions[opcode]; ok {
		return val, nil
	}

	return Instruction{Opcode: opcode}, OpCodeNotInInstructionSet
}

// Fill ensures that every possible opcode has a valid instruction by filling any opcode
// gaps with the passed in instruction. The Opcode value of the Instruction pass in is
// ignored and replaced with the actual Opcode whose place it takes.
func (is InstructionSet) Fill(instruction Instruction) error {

	if err := is.validate(); err != nil {
		return err
	}

	for opcode := range 255 {
		if _, ok := is.instructions[Opcode(opcode)]; !ok {
			instruction.Opcode = Opcode(opcode)
			is.instructions[Opcode(opcode)] = instruction
		}
	}

	return nil
}

// Opcodes returns a sorted slice of all the opcodes in the instruction set.
func (is InstructionSet) Opcodes() []Opcode {

	opcodes := make([]Opcode, 0, len(is.instructions))
	for _, instruction := range is.instructions {
		opcodes = append(opcodes, instruction.Opcode)
	}

	sort.Slice(opcodes, func(i, j int) bool {
		return opcodes[i] < opcodes[j]
	})

	return opcodes
}

// NewInstructionSet returns a correctly initialised InstructionSet based on the
// passed in slice of Instructions.
func NewInstructionSet(is Instructions) (InstructionSet, error) {

	instructions := make(map[Opcode]Instruction, len(is))
	for _, instruction := range is {
		instructions[instruction.Opcode] = instruction
	}

	result := InstructionSet{
		instructions: instructions,
	}
	if err := result.validate(); err != nil {
		return InstructionSet{}, err
	}

	return result, nil
}
