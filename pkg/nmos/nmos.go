package nmos

import "go6502/pkg/processor"

// TODO: NewExtendedCpu(memory Memory): Returns a Cpu with the full extended instruction set.
// TODO: Optionally fill in unused opcodes with Nops or Traps to catch issues.

// New6502Cpu returns a Cpu with the standard 6502 instruction set.
func New6502Cpu(memory processor.Memory) (processor.Cpu, error) {
	is, err := New6502InstructionSet()
	if err != nil {
		return processor.Cpu{}, err
	}
	return processor.NewCpu(is, memory)
}

// New65C02Cpu returns a Cpu with the standard 65C02 instruction set.
func New65C02Cpu(memory processor.Memory) (processor.Cpu, error) {
	is, err := New65C02InstructionSet()
	if err != nil {
		return processor.Cpu{}, err
	}
	return processor.NewCpu(is, memory)
}
