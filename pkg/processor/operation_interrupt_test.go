package processor

import "testing"

func TestBreak(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Break with mostly zeros",
			startState: State{SP: 0xFB},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant | FlagBreak, 0x01, 0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Break with PC set",
			startState: State{PC: 0xF0A0, SP: 0xFB},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant | FlagBreak, 0xA1, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Break with PC set and flags",
			startState: State{PC: 0xF0A0, SP: 0xFB, P: FlagCarry | FlagZero},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagCarry | FlagZero | FlagInterrupt},
			wantRam:    []uint8{0, FlagCarry | FlagZero | FlagConstant | FlagBreak, 0xA1, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Break with different SP set",
			startState: State{PC: 0xF0A0, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xFA, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant | FlagBreak, 0xA1, 0xF0, 0x02, 0xA0},
		},
		{
			name:       "Break with different vector",
			startState: State{PC: 0xF0A0, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0xEE, 0xFF},
			wantState:  State{PC: 0xFFEE, SP: 0xFA, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant | FlagBreak, 0xA1, 0xF0, 0xEE, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, Break)
		})
	}
}

func TestInterrupt(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Interrupt with mostly zeros",
			startState: State{SP: 0xFB},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0x00, 0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with PC set",
			startState: State{PC: 0xF0A0, SP: 0xFB},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with PC set and flags",
			startState: State{PC: 0xF0A0, SP: 0xFB, P: FlagCarry | FlagZero},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagCarry | FlagZero | FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant | FlagCarry | FlagZero, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with different SP set",
			startState: State{PC: 0xF0A0, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xFA, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with different vector",
			startState: State{PC: 0xF0A0, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0xEE, 0xFF},
			wantState:  State{PC: 0xFFEE, SP: 0xFA, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0xEE, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, Interrupt)
		})
	}
}

func TestNmi(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Nmi with mostly zeros",
			startState: State{SP: 0xF3},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF0, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with PC set",
			startState: State{PC: 0xF0A0, SP: 0xF3},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF0, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0xA0, 0xF0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with PC set and flags",
			startState: State{PC: 0xF0A0, SP: 0xF3, P: FlagCarry | FlagZero},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF0, P: FlagCarry | FlagZero | FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant | FlagCarry | FlagZero, 0xA0, 0xF0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with different SP set",
			startState: State{PC: 0xF0A0, SP: 0xF5},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF2, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with different vector",
			startState: State{PC: 0xF0A0, SP: 0xF5},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xEE, 0xFF, 0, 0, 0, 0},
			wantState:  State{PC: 0xFFEE, SP: 0xF2, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0, 0, 0, 0, 0xEE, 0xFF, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, Nmi)
		})
	}
}

func TestReturnFromInterrupt(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "RTI from Break with mostly zeros",
			startState: State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			startRam:   []uint8{0, FlagConstant | FlagBreak, 0x01, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0x0001, SP: 0xFB, P: FlagConstant},
			wantRam:    []uint8{0, FlagConstant | FlagBreak, 0x01, 0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Interrupt with mostly zeros",
			startState: State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			startRam:   []uint8{0, FlagConstant, 0x00, 0x00, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0x0000, SP: 0xFB, P: FlagConstant},
			wantRam:    []uint8{0, FlagConstant, 0x00, 0x00, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Break with PC set",
			startState: State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			startRam:   []uint8{0, FlagConstant | FlagBreak, 0xA1, 0xF0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xF0A1, SP: 0xFB, P: FlagConstant},
			wantRam:    []uint8{0, FlagConstant | FlagBreak, 0xA1, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Interrupt with PC set",
			startState: State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			startRam:   []uint8{0, FlagConstant, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xF0A0, SP: 0xFB, P: FlagConstant},
			wantRam:    []uint8{0, FlagConstant, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Break with PC set and flags",
			startState: State{PC: 0xA002, SP: 0xF8, P: FlagCarry | FlagZero | FlagInterrupt},
			startRam:   []uint8{0, FlagCarry | FlagZero | FlagConstant | FlagBreak, 0xA1, 0xF0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xF0A1, SP: 0xFB, P: FlagConstant | FlagCarry | FlagZero},
			wantRam:    []uint8{0, FlagCarry | FlagZero | FlagConstant | FlagBreak, 0xA1, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Interrupt with PC set and flags",
			startState: State{PC: 0xA002, SP: 0xF8, P: FlagCarry | FlagZero | FlagInterrupt},
			startRam:   []uint8{0, FlagConstant | FlagCarry | FlagZero, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xF0A0, SP: 0xFB, P: FlagConstant | FlagCarry | FlagZero},
			wantRam:    []uint8{0, FlagConstant | FlagCarry | FlagZero, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Break with different SP set",
			startState: State{PC: 0xA002, SP: 0xFA, P: FlagInterrupt},
			startRam:   []uint8{0, 0, 0, FlagConstant | FlagBreak, 0xA1, 0xF0, 0x02, 0xA0},
			wantState:  State{PC: 0xF0A1, SP: 0xFD, P: FlagConstant},
			wantRam:    []uint8{0, 0, 0, FlagConstant | FlagBreak, 0xA1, 0xF0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Interrupt with different SP set",
			startState: State{PC: 0xA002, SP: 0xFA, P: FlagInterrupt},
			startRam:   []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0x02, 0xA0},
			wantState:  State{PC: 0xF0A0, SP: 0xFD, P: FlagConstant},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0x02, 0xA0},
		},
		{
			name:       "RTI from Break with different vector",
			startState: State{PC: 0xFFEE, SP: 0xFA, P: FlagInterrupt},
			startRam:   []uint8{0, 0, 0, FlagBreak | FlagOverflow | FlagDecimal, 0xA1, 0xF0, 0xEE, 0xFF},
			wantState:  State{PC: 0xF0A1, SP: 0xFD, P: FlagConstant | FlagOverflow | FlagDecimal},
			wantRam:    []uint8{0, 0, 0, FlagBreak | FlagOverflow | FlagDecimal, 0xA1, 0xF0, 0xEE, 0xFF},
		},
		{
			name:       "RTI from Interrupt with different vector",
			startState: State{PC: 0xFFEE, SP: 0xFA, P: FlagInterrupt},
			startRam:   []uint8{0, 0, 0, FlagOverflow | FlagDecimal, 0xA0, 0xF0, 0xEE, 0xFF},
			wantState:  State{PC: 0xF0A0, SP: 0xFD, P: FlagConstant | FlagOverflow | FlagDecimal},
			wantRam:    []uint8{0, 0, 0, FlagOverflow | FlagDecimal, 0xA0, 0xF0, 0xEE, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ReturnFromInterrupt)
		})
	}
}
