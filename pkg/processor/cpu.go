package processor

import (
	"fmt"
	"math"
)

const StackPointerStart = 0xFD

// State represents the entire state of the Cpu at a specific point.
type State struct {
	PC Address
	SP uint8
	A  uint8
	X  uint8
	Y  uint8
	P  Status
}

func (s State) String() string {
	return fmt.Sprintf("PC: 0x%04X, SP: 0x%02X, A: 0x%02X, X: 0x%02X, Y: 0x%02X, P: %v", s.PC, s.SP, s.A, s.X, s.Y, s.P.String())
}

// Cpu represents the actual Cpu
type Cpu struct {
	State          State
	memory         Memory
	instructionSet InstructionSet
}

// NewCpu returns an initialised Cpu that supports the provided instruction set
// and is connected to the supplied memory. A Reset() should be called on the
// newly constructed CPU to set up the reset vector and stack pointer.
func NewCpu(is InstructionSet, memory Memory) (Cpu, error) {
	if err := is.validate(); err != nil {
		return Cpu{}, err
	}
	if memory == nil {
		return Cpu{}, MemoryMustBeProvided
	}
	return Cpu{memory: memory, instructionSet: is}, nil
}

// Reset should be called before execution begins. It clears all flags, clears
// registers X, Y and A, sets SP to 0xFD and sets PC to the reset vector that
// is stored in RAM at 0xFFFC and 0xFFFD.
func (c *Cpu) Reset() error {
	if c == nil {
		return UninitialisedCpu
	}

	start, err := ReadResetVectorFromMemory(c.memory)
	if err != nil {
		return err
	}

	c.State = State{PC: start, SP: StackPointerStart}

	return nil
}

// Step a single machine instruction. This will read the opcode, advance the
// program counter and then execute the instruction. If the opcode read is not
// present in the CPUs instruction set then an error is returned. If there is
// an error executing the instruction then an error is returned and the state
// changes relating to the instruction execution are not applied. In all cases
// the program counter is incremented by at least 1 byte. Details of the number
// of CPU cycles that have elapsed are returned; this will always be at least
// 1 for a valid Cpu instance.
func (c *Cpu) Step() (uint, error) {
	if c == nil {
		return 0, UninitialisedCpu
	}

	opcode := Opcode(c.memory.Read(c.State.PC))
	c.State.PC++

	instruction, err := c.instructionSet.Get(opcode)
	if err != nil {
		return 1, err
	}

	// If there is an error executing the instruction (which should not happen)
	// then we return an error and do not apply the instruction state changes.
	newState, cycles, err := instruction.Execute(c.State, c.memory)
	if err != nil {
		return cycles + 1, err
	}
	c.State = newState

	return cycles + 1, nil
}

// Execute will execute instructions until the specified number of cycles have been
// passed; returning the actual number of cycles that have cycled. The number of cycles
// actually executed may be more than those specified if the last instruction executed
// takes it over the limit. Specifying a value of zero for cycles will let the CPU run
// continuously. If an unknown instruction is executed then Execute also stops.
func (c *Cpu) Execute(cycles uint) (uint, error) {
	if c == nil {
		return 0, UninitialisedCpu
	}

	if cycles == 0 {
		cycles = math.MaxUint - 100 // We need this to avoid overflow
	}
	elapsedCycles := uint(0)

	for elapsedCycles < cycles {
		stepCycles, err := c.Step()
		elapsedCycles += stepCycles
		if err != nil {
			return elapsedCycles, err
		}
	}

	return elapsedCycles, nil
}

// Nmi will trigger a non-maskable interrupt in the 6502 core. The current PC value
// is pushed to the stack (high byte first, low byte second). The status register is
// then pushed onto the stack. The Interrupt flag is set then the NMI vector stored
// at address 0xFFFA (low byte) and 0xFFFB (high byte) is loaded into the PC ready
// to execute. No actual instructions are executed.
func (c *Cpu) Nmi() error {
	if c == nil {
		return UninitialisedCpu
	}
	addressing := Addressing{Memory: c.memory}
	state, err := Nmi(c.State, addressing)
	if err != nil {
		return err
	}
	c.State = state
	return nil
}

// Interrupt will trigger a hardware interrupt in the 6502 core. The current PC value
// is pushed to the stack (high byte first, low byte second). The status register is
// then pushed onto the stack. The Interrupt flag is set then the IRQ vector stored
// at address 0xFFFE (low byte) and 0xFFFF (high byte) is loaded into the PC ready
// to execute. No actual instructions are executed. If the processor status flag has
// the Interrupt flag set when calling this method, it does nothing.
func (c *Cpu) Interrupt() error {
	if c == nil {
		return UninitialisedCpu
	}

	// If the interrupt disable flag is set then ignore the request.
	if c.State.P.ToFlags().Interrupt {
		return nil
	}

	addressing := Addressing{Memory: c.memory}
	state, err := Interrupt(c.State, addressing)
	if err != nil {
		return err
	}
	c.State = state
	return nil
}

// Opcodes returns a sorted slice of all the opcodes the Cpu has in its instruction set.
func (c *Cpu) Opcodes() ([]Opcode, error) {
	if c == nil {
		return []Opcode{}, UninitialisedCpu
	}

	return c.instructionSet.Opcodes(), nil
}

// Memory returns a reference to the memory attached to the CPU.
func (c *Cpu) Memory() (Memory, error) {
	if c == nil {
		return nil, UninitialisedCpu
	}

	return c.memory, nil
}
