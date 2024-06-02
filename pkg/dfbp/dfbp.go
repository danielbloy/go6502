package dfbp

import (
	"go6502/pkg/processor"
)

// A 6502 with the addition of a 16-bit registers and multiplication
// instructions using the special 16-bit register. The purpose of this
// processor is to demonstrate how a custom instruction set can be
// specified and provide additional custom registers.

type Cpu struct {
	processor.Cpu

	B uint16
}

// NewDfbpCpu returns a Dfbp Cpu with a custom instruction set.
func NewDfbpCpu(memory processor.Memory) (Cpu, error) {
	is, err := NewDfbpInstructionSet()
	if err != nil {
		return Cpu{}, err
	}

	// TODO: Need to add an opcode for our multiplication
	// TODO: Need to work out how an operation can access B.

	cpu, err := processor.NewCpu(is, memory)
	if err != nil {
		return Cpu{}, err
	}

	return Cpu{cpu, 0}, nil
}
