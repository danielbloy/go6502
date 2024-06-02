package nmos

import (
	"go6502/pkg/processor"
)

// New6502InstructionSet returns a correctly initialised InstructionSet
// for the 6502 CPU.
func New6502InstructionSet() (processor.InstructionSet, error) {

	instructions := make(processor.Instructions, 0, 0xFF)

	for _, mnemonic := range processor.AllOpcodes() {
		instruction := processor.NewInstruction(mnemonic)
		instructions = append(instructions, instruction)
	}
	return processor.NewInstructionSet(instructions)
}

// New65C02InstructionSet returns a correctly initialised InstructionSet
// for the 65C02 CPU.
func New65C02InstructionSet() (processor.InstructionSet, error) {
	instructions, err := New6502InstructionSet()
	if err != nil {
		return processor.InstructionSet{}, err
	}

	// TODO: Add in 65C02 specific instructions.

	return instructions, nil
}
