package processor

import "testing"

func TestJump(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Jump to 0x0000",
			startState: State{PC: 0x010A},
			addressing: Addressing{EffectiveAddress: 0x0000},
			wantState:  State{PC: 0x0000},
		},
		{
			name:       "Jump to 0xFFFF",
			addressing: Addressing{EffectiveAddress: 0xFFFF},
			startState: State{PC: 0x010A},
			wantState:  State{PC: 0xFFFF},
		},
		{
			name:       "Jump to 0xABCD",
			addressing: Addressing{EffectiveAddress: 0xABCD},
			startState: State{PC: 0x010A},
			wantState:  State{PC: 0xABCD},
		},
		{
			name:       "Jump to 0x1234",
			addressing: Addressing{EffectiveAddress: 0x1234},
			startState: State{PC: 0x010A},
			wantState:  State{PC: 0x1234},
		},

		{
			name:       "Jump to 0xABCD; validate SP, A, X, Y and Flags are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
			addressing: Addressing{EffectiveAddress: 0xABCD},
			wantState:  State{PC: 0xABCD, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, Jump)
		})
	}
}

func TestJumpSubRoutine(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "JSR from 0x0000 to 0x0000",
			startState: State{PC: 0x010A, SP: 0xFF},
			addressing: Addressing{EffectiveAddress: 0x0000},
			wantState:  State{PC: 0x0000, SP: 0xFD},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x09, 0x01},
		},
		{
			name:       "JSR from 0x010A to 0xFFFF",
			addressing: Addressing{EffectiveAddress: 0xFFFF},
			startState: State{PC: 0x010A, SP: 0xFF},
			wantState:  State{PC: 0xFFFF, SP: 0xFD},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x09, 0x01},
		},
		{
			name:       "JSR from 0x010A to 0xABCD",
			addressing: Addressing{EffectiveAddress: 0xABCD},
			startState: State{PC: 0x010A, SP: 0xFE},
			wantState:  State{PC: 0xABCD, SP: 0xFC},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0x09, 0x01, 0},
		},
		{
			name:       "JSR from 0x010A to 0x1234, with wrap around of SP - 1",
			addressing: Addressing{EffectiveAddress: 0x1234},
			startState: State{PC: 0x010A, SP: 0x01},
			wantState:  State{PC: 0x1234, SP: 0xFF},
			wantRam:    []uint8{0x09, 0x01, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "JSR from 0x010A to 0x1234, with wrap around of SP - 2",
			addressing: Addressing{EffectiveAddress: 0x1234},
			startState: State{PC: 0x010A, SP: 0x00},
			wantState:  State{PC: 0x1234, SP: 0xFE},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0x09},
		},
		{
			name:       "JSR from 0x1234 to 0xABCD; A, X, Y and Flags are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
			addressing: Addressing{EffectiveAddress: 0xABCD},
			wantState:  State{PC: 0xABCD, SP: 0x0E, A: 0x20, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
			wantRam:    []uint8{0x12, 0, 0, 0, 0, 0, 0, 0x33},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, JumpSubRoutine)
		})
	}
}

func TestReturnFromSubroutine(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "RTS from 0x0000 to 0x010A",
			startState: State{PC: 0x0000, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x09, 0x01},
			wantState:  State{PC: 0x010A, SP: 0xFF},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x09, 0x01},
		},
		{
			name:       "RTS from 0xFFFF to 0x010A",
			startState: State{PC: 0xFFFF, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x09, 0x01},
			wantState:  State{PC: 0x010A, SP: 0xFF},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x09, 0x01},
		},
		{
			name:       "RTS from 0xABCD to 0x010A",
			startState: State{PC: 0xABCD, SP: 0xFC},
			startRam:   []uint8{0, 0, 0, 0, 0, 0x09, 0x01, 0},
			wantState:  State{PC: 0x010A, SP: 0xFE},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0x09, 0x01, 0},
		},
		{
			name:       "RTS from 0x1234 to 0x010A, with wrap around of SP - 1",
			startState: State{PC: 0x1234, SP: 0xFF},
			startRam:   []uint8{0x09, 0x01, 0, 0, 0, 0, 0, 0},
			wantState:  State{PC: 0x010A, SP: 0x01},
			wantRam:    []uint8{0x09, 0x01, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "RTS from 0x1234 to 0x010A, with wrap around of SP - 2",
			startState: State{PC: 0x1234, SP: 0xFE},
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0x09},
			wantState:  State{PC: 0x010A, SP: 0x00},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0x09},
		},
		{
			name:       "RTS from 0xABCD to 0x1234; A, X, Y and Flags are not touched.",
			startState: State{PC: 0xABCD, SP: 0x0E, A: 0x20, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
			startRam:   []uint8{0x12, 0, 0, 0, 0, 0, 0, 0x33},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
			wantRam:    []uint8{0x12, 0, 0, 0, 0, 0, 0, 0x33},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ReturnFromSubroutine)
		})
	}
}
