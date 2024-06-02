package processor

import (
	"reflect"
	"testing"
)

// Because of the significant number of test, they have been split up into multiple
// file using the groups based on: Instructions by Type from:
//    https://www.masswerk.at/6502/6502_instruction_set.html#bytype

// The following struct and harness function are used to test the individual
// operation functions.
type testOperationConfig struct {
	name       string
	startState State
	startRam   []uint8
	addressing Addressing
	wantState  State
	wantRam    []uint8
	wantErr    bool
}

func testOperation(t *testing.T, tt testOperationConfig, operation Operation) {
	// Default the start and end RAMs as empty if not supplied
	if len(tt.startRam) <= 0 {
		tt.startRam = []uint8{0, 0, 0, 0, 0, 0, 0, 0}
	}
	if len(tt.wantRam) <= 0 {
		tt.wantRam = []uint8{0, 0, 0, 0, 0, 0, 0, 0}
	}
	memory := NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
	tt.addressing.Memory = &memory
	got, err := operation(tt.startState, tt.addressing)

	if (err != nil) != tt.wantErr {
		t.Errorf("testOperation() error = %v, wantErr %v", err, tt.wantErr)
		return
	}

	// Validate the final State.
	if !reflect.DeepEqual(got, tt.wantState) {
		t.Errorf("testOperation() State got = %v, want = %v", got, tt.wantState)
	}

	// Validate expected final RAM state
	if !reflect.DeepEqual(memory.ram, tt.wantRam) {
		t.Errorf("testOperation() RAM state got = %v, want = %v", memory.ram, tt.wantRam)
	}
}

func TestNoOperation(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "In State is the same as the return State 1",
		},
		{
			name:       "In State is the same as the return State 2",
			startState: State{PC: 0x20, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			wantState:  State{PC: 0x20, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
		},
		{
			name:       "In State is the same as the return State 2",
			startState: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			wantState:  State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, NoOperation)
		})
	}
}

func TestTestBitsInMemoryWithAccumulator(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "BIT zero in memory and zero in accumulator, sets zero flag only.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "BIT zero in memory and zero in accumulator, sets zero flag and clears negative and sign.",
			startState: State{P: 0xFF},
			wantState:  State{P: 0xFF ^ FlagNegative ^ FlagOverflow},
		},
		{
			name:       "BIT 0x80 in memory, sets negative flag only.",
			startState: State{A: 0xFF},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{
			name:       "BIT 0x40 in memory, sets overflow flag only.",
			startState: State{A: 0xFF},
			addressing: Addressing{Value: 0x40},
			wantState:  State{A: 0xFF, P: FlagOverflow},
		},
		{
			name:       "BIT 0x3F in memory and zero in accumulator, clears zero flag.",
			startState: State{A: 0xFF, P: FlagZero},
			addressing: Addressing{Value: 0x3F},
			wantState:  State{A: 0xFF, P: 0x00},
		},
		{
			name:       "BIT 0x3F in memory, clears negative flag.",
			startState: State{A: 0xFF, P: FlagZero},
			addressing: Addressing{Value: 0x3F},
			wantState:  State{A: 0xFF, P: 0x00},
		},
		{
			name:       "BIT 0x3F in memory, clears overflow flag.",
			startState: State{A: 0xFF, P: FlagZero},
			addressing: Addressing{Value: 0x3F},
			wantState:  State{A: 0xFF, P: 0x00},
		},
		{
			name:       "BIT result of Accumulator AND memory results in zero, sets Zero - 1.",
			startState: State{A: 0xFF},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0xFF, P: FlagZero},
		},
		{
			name:       "BIT result of Accumulator AND memory results in zero, sets Zero - 2.",
			startState: State{A: 0xF0},
			addressing: Addressing{Value: 0x0F},
			wantState:  State{A: 0xF0, P: FlagZero},
		},
		{
			name:       "BIT result of Accumulator AND memory results in zero, sets Zero - 1.",
			startState: State{A: 0xC3},
			addressing: Addressing{Value: 0xFF ^ 0xC3},
			wantState:  State{A: 0xC3, P: FlagZero},
		},
		{
			name:       "BIT result of Accumulator AND memory results in non-zero, clears Zero - 1.",
			startState: State{A: 0xFF, P: FlagZero},
			addressing: Addressing{Value: 0x03},
			wantState:  State{A: 0xFF, P: 0x00},
		},
		{
			name:       "BIT result of Accumulator AND memory results in non-zero, clears Zero - 2.",
			startState: State{A: 0xC1, P: FlagZero},
			addressing: Addressing{Value: 0x13},
			wantState:  State{A: 0xC1, P: 0x00},
		},
		{
			name:       "BIT result of Accumulator AND memory results in zero, negative and overflow are set too.",
			startState: State{A: 0x0F},
			addressing: Addressing{Value: 0xC0},
			wantState:  State{A: 0x0F, P: FlagZero | FlagNegative | FlagOverflow},
		},
		{
			name:       "BIT result of Accumulator AND memory results in non-zero, negative and overflow are cleared too.",
			startState: State{A: 0xF1, P: FlagZero | FlagNegative | FlagOverflow},
			addressing: Addressing{Value: 0x03},
			wantState:  State{A: 0xF1, P: 0x00},
		},
		{
			name:       "BIT operation does not change other settings, sets zero, negative and overflow flags only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x00, X: 0x02, Y: 0x03, P: 0xFF ^ FlagZero ^ FlagNegative ^ FlagOverflow},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x00, X: 0x02, Y: 0x03, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, TestBitsInMemoryWithAccumulator)
		})
	}
}
