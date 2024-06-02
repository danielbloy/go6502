package processor

import "errors"

var (
	InstructionSetEmpty       = errors.New("the instruction set is empty")
	OpCodeNotInInstructionSet = errors.New("the opcode is not present in the instruction set")

	MemoryMustBeProvided      = errors.New("a valid memory must be provided")
	InvalidMemorySizeProvided = errors.New("invalid memory size was provided")

	UninitialisedCpu = errors.New("the CPU has not been initialised correctly")

	NoAddressingModeFunction = errors.New("the instruction has no addressing mode function")
	NoOperationFunction      = errors.New("the instruction has no operation function")
)
