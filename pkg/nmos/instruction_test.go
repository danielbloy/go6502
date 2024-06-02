package nmos

import (
	"go6502/pkg/processor"
	"reflect"
	"testing"
)

const (
	brk          = 0x00
	adcZeroPage  = 0x65
	adcImmediate = 0x69
	nop          = 0xEA
)

// The following struct and harness function are used to test an instruction test
// for a Cpu. A Cpu is created and initialised for each test case with the instruction
// set defined by the function. The test case is then executed to ensure the operation
// executes as expected. The test case starts by resetting the CPU and then performing
// the defined number of instruction steps (1 if zero is specified).
type testCpuInstructionSet struct {
	name       string
	steps      uint
	startState processor.State
	startRam   []uint8
	wantCycles uint
	wantState  processor.State
	wantRam    []uint8
	wantErr    bool
}

func testInstructionSet(t *testing.T, tests []testCpuInstructionSet, instructionSet func() processor.InstructionSet) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callback := func(cpu *processor.Cpu) error {
				if err := cpu.Reset(); err != nil {
					return err
				}

				// If a start state was supplied, override the reset CPU.
				if tt.startState != (processor.State{}) {
					cpu.State = tt.startState
				}

				// Perform at least 1 step
				if tt.steps == 0 {
					tt.steps = 1
				}

				// Perform the steps
				cycles := uint(0)
				for range tt.steps {
					delta, err := cpu.Step()
					cycles += delta
					if err != nil {
						return err
					}
				}

				if cycles != tt.wantCycles {
					t.Errorf("TestInstructionSet() did not execute the expected number of cycles, got = %v, want = %v", cycles, tt.wantCycles)
				}
				return nil
			}

			// Create CPU with populated memory and a default test instruction set if one is not specified.
			memory := processor.NewPopulatedRam(processor.RepeatingRamSize(len(tt.startRam)), tt.startRam)
			is := instructionSet()

			cpu, err := processor.NewCpu(is, &memory)
			if err != nil {
				panic(err)
			}

			cpu.State = tt.startState

			// Call the required CPU method using the provided callback.
			if err = callback(&cpu); (err != nil) != tt.wantErr {
				t.Errorf("testCpuMethod() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Validate expected final CPU State
			if !reflect.DeepEqual(cpu.State, tt.wantState) {
				t.Errorf("Unexpected CPU state got = %v, want = %v", cpu.State, tt.wantState)
			}

			// Validate expected final RAM state
			wantMemory := processor.NewPopulatedRam(processor.RepeatingRamSize(len(tt.wantRam)), tt.wantRam)
			if !reflect.DeepEqual(memory, wantMemory) {
				t.Errorf("Unexpected RAM state got = %v, want = %v", memory, wantMemory)
			}
		})
	}
}

func instructionSetTests6502() []testCpuInstructionSet {
	tests := []testCpuInstructionSet{
		{
			name:       "ADC immediate ; $20 to $0",
			startRam:   []uint8{adcImmediate, 0x20, 0, 0, 0, 0, 0, 0},
			wantCycles: 2,
			wantRam:    []uint8{adcImmediate, 0x20, 0, 0, 0, 0, 0, 0},
			wantState:  processor.State{PC: 2, A: 0x20, SP: processor.StackPointerStart},
		},
		{
			name:       "ADC immediate ; $21 to $20",
			steps:      2,
			startRam:   []uint8{adcImmediate, 0x20, adcImmediate, 0x21, 0, 0, 0, 0},
			wantCycles: 4,
			wantRam:    []uint8{adcImmediate, 0x20, adcImmediate, 0x21, 0, 0, 0, 0},
			wantState:  processor.State{PC: 4, A: 0x41, SP: processor.StackPointerStart},
		},
		{
			name:       "ADC zpg ; value at $0F ($20) to $0",
			startRam:   []uint8{adcZeroPage, 0x0F, 0, 0, 0, 0, 0, 0x20},
			wantCycles: 3,
			wantRam:    []uint8{adcZeroPage, 0x0F, 0, 0, 0, 0, 0, 0x20},
			wantState:  processor.State{PC: 2, A: 0x20, SP: processor.StackPointerStart},
		},
		{
			name:       "ADC zpg ; value $0F (0x20) to $11",
			steps:      2,
			startRam:   []uint8{adcImmediate, 0x11, adcZeroPage, 0x0F, 0, 0, 0, 0x20},
			wantCycles: 5,
			wantRam:    []uint8{adcImmediate, 0x11, adcZeroPage, 0x0f, 0, 0, 0, 0x20},
			wantState:  processor.State{PC: 4, A: 0x31, SP: processor.StackPointerStart},
		},
		{
			name:       "BRK",
			startState: processor.State{PC: 0xF0A0, SP: 0xFB},
			startRam:   []uint8{brk, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantCycles: 7,
			wantState:  processor.State{PC: 0xA002, SP: 0xF8, P: processor.FlagInterrupt},
			wantRam:    []uint8{brk, processor.FlagConstant | processor.FlagBreak, 0xA2, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "NOP",
			wantCycles: 2,
			startRam:   []uint8{nop, 0, 0, 0, 0, 0, 0, 0},
			wantRam:    []uint8{nop, 0, 0, 0, 0, 0, 0, 0},
			wantState:  processor.State{PC: 1, SP: processor.StackPointerStart},
		},
	}

	return tests
}

// This tests some 6502 instructions, validating the CPU and memory
// state is correct afterward. It makes use of the testCpuMethod() function.
// The exhaustive tests are done using the Klaus2m5 test suite.
func TestNew6502InstructionSet(t *testing.T) {

	instructionSet := func() processor.InstructionSet {
		result, err := New6502InstructionSet()
		if err != nil {
			panic(err)
		}
		return result
	}

	testInstructionSet(t, instructionSetTests6502(), instructionSet)
}

// This tests some 65C02 instructions, validating the CPU and memory
// state is correct afterward. It makes use of the testCpuMethod() function.
// The exhaustive tests are done using the Klaus2m5 test suite.
func TestNew65C02InstructionSet(t *testing.T) {

	instructionSetTests65C02 := func() []testCpuInstructionSet {
		tests := instructionSetTests6502()

		// TODO: Add some tests for specific 65C02 instructions.

		return tests
	}

	instructionSet := func() processor.InstructionSet {
		result, err := New65C02InstructionSet()
		if err != nil {
			panic(err)
		}
		return result
	}

	testInstructionSet(t, instructionSetTests65C02(), instructionSet)
}
