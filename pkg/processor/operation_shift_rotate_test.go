package processor

import "testing"

func TestArithmeticShiftLeft(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "ASL zero in memory, sets zero flag.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "ASL zero in memory with carry flag set, sets zero flag.",
			startState: State{P: FlagCarry},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "ASL 0xFF in memory, sets carry and negative flags.",
			addressing: Addressing{Value: 0xFF},
			wantState:  State{P: FlagNegative | FlagCarry},
			wantRam:    []uint8{0xFE, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ASL 0x01 in memory, sets no flags.",
			addressing: Addressing{Value: 0x01},
			wantRam:    []uint8{0x02, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ASL 0x80 in memory, sets carry flag and zero flag.",
			addressing: Addressing{Value: 0x80},
			wantState:  State{P: FlagCarry | FlagZero},
			wantRam:    []uint8{0x00, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ASL 0x81 in memory, sets carry flag only.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagCarry},
			wantRam:    []uint8{0x02, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ASL 0x40 in memory, sets negative flag only.",
			addressing: Addressing{Value: 0x40},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ASL 0x40 in memory with different effective address, sets negative flag only.",
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x40},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0, 0, 0, 0x80, 0, 0, 0, 0},
		},
		{
			name:       "ASL 0x40 in memory does not change other settings, sets negative flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x40},
			wantRam:    []uint8{0, 0, 0, 0x80, 0, 0, 0, 0},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagZero ^ FlagCarry},
		},
		{
			name:       "ASL zero in accumulator, sets zero flag.",
			wantState:  State{P: FlagZero},
			addressing: Addressing{Accumulator: true},
		},
		{
			name:       "ASL zero in accumulator with carry flag set, sets zero flag.",
			addressing: Addressing{Accumulator: true},
			startState: State{P: FlagCarry},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "ASL 0xFF in accumulator, sets carry and negative flags.",
			addressing: Addressing{Value: 0xFF, Accumulator: true},
			wantState:  State{A: 0xFE, P: FlagNegative | FlagCarry},
		},
		{
			name:       "ASL 0x01 in accumulator, sets no flags.",
			addressing: Addressing{Value: 0x01, Accumulator: true},
			wantState:  State{A: 0x02},
		},
		{
			name:       "ASL 0x80 in accumulator, sets carry flag and zero flag.",
			addressing: Addressing{Value: 0x80, Accumulator: true},
			wantState:  State{A: 0x00, P: FlagCarry | FlagZero},
		},
		{
			name:       "ASL 0x81 in accumulator, sets carry flag only.",
			addressing: Addressing{Value: 0x81, Accumulator: true},
			wantState:  State{A: 0x02, P: FlagCarry},
		},
		{
			name:       "ASL 0x40 in accumulator, sets negative flag only.",
			addressing: Addressing{Value: 0x40, Accumulator: true},
			wantState:  State{A: 0x80, P: FlagNegative},
		},
		{
			name:       "ASL 0x40 in accumulator does not change other settings, sets negative flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x00, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
			addressing: Addressing{Value: 0x40, Accumulator: true},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x80, X: 0x02, Y: 0x03, P: 0xFF ^ FlagZero ^ FlagCarry},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ArithmeticShiftLeft)
		})
	}
}

func TestLogicalShiftRight(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "LSR zero in memory, sets zero flag.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "LSR 0xFF in memory, sets carry flag.",
			addressing: Addressing{Value: 0xFF},
			wantState:  State{P: FlagCarry},
			wantRam:    []uint8{0x7F, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "LSR 0x01 in memory, sets zero and carry flags.",
			addressing: Addressing{Value: 0x01},
			wantState:  State{P: FlagCarry | FlagZero},
			wantRam:    []uint8{0x0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "LSR 0x80 in memory, sets no flags.",
			addressing: Addressing{Value: 0x80},
			wantRam:    []uint8{0x40, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "LSR 0x81 in memory, sets carry flag only.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagCarry},
			wantRam:    []uint8{0x40, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "LSR 0x40 in memory, sets no flags.",
			addressing: Addressing{Value: 0x40},
			wantRam:    []uint8{0x20, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "LSR 0x01 in memory with different effective address, sets carry and zero flags.",
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x01},
			wantState:  State{P: FlagCarry | FlagZero},
			wantRam:    []uint8{0, 0, 0, 0x00, 0, 0, 0, 0},
		},
		{
			name:       "LSR 0x81 in memory does not change other settings, sets carry flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x81},
			wantRam:    []uint8{0, 0, 0, 0x40, 0, 0, 0, 0},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagZero ^ FlagNegative},
		},
		{
			name:       "LSR zero in accumulator, sets zero flag.",
			addressing: Addressing{Accumulator: true},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "LSR 0xFF in accumulator, sets carry flag.",
			addressing: Addressing{Value: 0xFF, Accumulator: true},
			wantState:  State{A: 0x7F, P: FlagCarry},
		},
		{
			name:       "LSR 0x01 in accumulator, sets zero and carry flags.",
			addressing: Addressing{Value: 0x01, Accumulator: true},
			wantState:  State{A: 0x00, P: FlagZero | FlagCarry},
		},
		{
			name:       "LSR 0x80 in accumulator, sets no flags.",
			addressing: Addressing{Value: 0x80, Accumulator: true},
			wantState:  State{A: 0x40},
		},
		{
			name:       "LSR 0x81 in accumulator, sets carry flag.",
			addressing: Addressing{Value: 0x81, Accumulator: true},
			wantState:  State{A: 0x40, P: FlagCarry},
		},
		{
			name:       "LSR 0x40 in accumulator, sets no flags.",
			addressing: Addressing{Value: 0x40, Accumulator: true},
			wantState:  State{A: 0x20},
		},
		{
			name:       "LSR 0x81 in accumulator does not change other settings, sets carry flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x00, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
			addressing: Addressing{Value: 0x81, Accumulator: true},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x40, X: 0x02, Y: 0x03, P: 0xFF ^ FlagZero ^ FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, LogicalShiftRight)
		})
	}
}

func TestRotateLeft(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "ROL zero in memory, sets zero flag.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "ROL zero in memory with carry flag set.",
			startState: State{P: FlagCarry},
			wantState:  State{P: 0},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0xFF in memory, sets carry and negative flags.",
			addressing: Addressing{Value: 0xFF},
			startRam:   []uint8{0xFF, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{P: FlagNegative | FlagCarry},
			wantRam:    []uint8{0xFE, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0x01 in memory, sets no flags.",
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			addressing: Addressing{Value: 0x01},
			wantRam:    []uint8{0x02, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0x80 in memory, sets carry flag and zero flag.",
			addressing: Addressing{Value: 0x80},
			wantState:  State{P: FlagCarry | FlagZero},
			wantRam:    []uint8{0x00, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0x81 in memory, sets carry flag only.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagCarry},
			wantRam:    []uint8{0x02, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0x40 in memory, sets negative flag only.",
			addressing: Addressing{Value: 0x40},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0x40 in memory with carry, sets negative flag only.",
			startState: State{P: FlagCarry},
			addressing: Addressing{Value: 0x40},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0x81, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0x40 in memory with different effective address, sets negative flag only.",
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x40},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0, 0, 0, 0x80, 0, 0, 0, 0},
		},
		{
			name:       "ROL 0x40 in memory does not change other settings, sets negative flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x40},
			wantRam:    []uint8{0, 0, 0, 0x81, 0, 0, 0, 0},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagZero ^ FlagCarry},
		},
		{
			name:       "ROL zero in accumulator, sets zero flag.",
			addressing: Addressing{Accumulator: true},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "ROL zero in accumulator with carry set.",
			startState: State{P: FlagCarry},
			addressing: Addressing{Accumulator: true},
			wantState:  State{A: 0x01},
		},
		{
			name:       "ROL 0xFF in accumulator, sets carry and negative flags.",
			addressing: Addressing{Value: 0xFF, Accumulator: true},
			wantState:  State{A: 0xFE, P: FlagNegative | FlagCarry},
		},
		{
			name:       "ROL 0x01 in accumulator, sets no flags, doesn't change ram.",
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{Value: 0x01, Accumulator: true},
			wantState:  State{A: 0x02},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
		},
		{
			name:       "ROL 0x80 in accumulator, sets carry flag and zero flag.",
			addressing: Addressing{Value: 0x80, Accumulator: true},
			wantState:  State{A: 0x00, P: FlagCarry | FlagZero},
		},
		{
			name:       "ROL 0x81 in accumulator, sets carry flag only.",
			addressing: Addressing{Value: 0x81, Accumulator: true},
			wantState:  State{A: 0x02, P: FlagCarry},
		},
		{
			name:       "ROL 0x40 in accumulator, sets negative flag only.",
			addressing: Addressing{Value: 0x40, Accumulator: true},
			wantState:  State{A: 0x80, P: FlagNegative},
		},
		{
			name:       "ROL 0x40 in accumulator with carry, sets negative flag only.",
			startState: State{P: FlagCarry},
			addressing: Addressing{Value: 0x40, Accumulator: true},
			wantState:  State{A: 0x81, P: FlagNegative},
		},
		{
			name:       "ROL 0x40 in accumulator does not change other settings, sets negative flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x00, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
			addressing: Addressing{Value: 0x40, Accumulator: true},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x81, X: 0x02, Y: 0x03, P: 0xFF ^ FlagZero ^ FlagCarry},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, RotateLeft)
		})
	}
}

func TestRotateRight(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "ROR zero in memory, sets zero flag.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "ROR zero in memory with carry flag set.",
			startState: State{P: FlagCarry},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x01 in memory sets carry flag.",
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			addressing: Addressing{Value: 0x01},
			wantState:  State{P: FlagCarry | FlagZero},
			wantRam:    []uint8{0x00, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x01 in memory with carry keeps carry flag and sets negative.",
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			startState: State{P: FlagCarry},
			addressing: Addressing{Value: 0x01},
			wantState:  State{P: FlagCarry | FlagNegative},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},

		{
			name:       "ROR 0xFF in memory, sets carry and clears negative flags.",
			addressing: Addressing{Value: 0xFF},
			startRam:   []uint8{0xFF, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{P: FlagCarry},
			wantRam:    []uint8{0x7F, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x01 in memory, sets carry and zero flag.",
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			addressing: Addressing{Value: 0x01},
			wantState:  State{P: FlagCarry | FlagZero},
			wantRam:    []uint8{0x00, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x80 in memory, sets no flags.",
			addressing: Addressing{Value: 0x80},
			wantRam:    []uint8{0x40, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x81 in memory, sets carry flag only.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagCarry},
			wantRam:    []uint8{0x40, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x40 in memory with carry, sets negative flag only.",
			startState: State{P: FlagCarry},
			addressing: Addressing{Value: 0x40},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0xA0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x40 in memory with carry different effective address, sets negative flag only.",
			startState: State{P: FlagCarry},
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x40},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0, 0, 0, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x01 in memory does not change other settings, sets carry flag and clears negative flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagCarry},
			addressing: Addressing{EffectiveAddress: 0x03, Value: 0x01},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x01, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
			wantRam:    []uint8{0, 0, 0, 0x00, 0, 0, 0, 0},
		},
		{
			name:       "ROR zero in accumulator, sets zero flag.",
			addressing: Addressing{Accumulator: true},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "ROR zero in accumulator with carry flag set, also does not change ram.",
			startState: State{P: FlagCarry},
			startRam:   []uint8{0x12, 0, 0, 0, 0, 0, 0, 0},
			addressing: Addressing{Accumulator: true},
			wantState:  State{A: 0x80, P: FlagNegative},
			wantRam:    []uint8{0x12, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x01 in accumulator sets carry flag.",
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			addressing: Addressing{Value: 0x01, Accumulator: true},
			wantState:  State{P: FlagCarry | FlagZero},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "ROR 0x01 in accumulator with carry keeps carry flag and sets negative, also does not change ram.",
			startRam:   []uint8{0x34, 0, 0, 0, 0, 0, 0, 0},
			startState: State{A: 0x50, P: FlagCarry},
			addressing: Addressing{Value: 0x01, Accumulator: true},
			wantState:  State{A: 0x80, P: FlagCarry | FlagNegative},
			wantRam:    []uint8{0x34, 0, 0, 0, 0, 0, 0, 0},
		},

		{
			name:       "ROR 0xFF in accumulator, sets carry and clears negative flags.",
			addressing: Addressing{Value: 0xFF, Accumulator: true},
			wantState:  State{A: 0x7F, P: FlagCarry},
		},
		{
			name:       "ROR 0x01 in accumulator, sets carry and zero flag, doesn't change ram.",
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{Value: 0x01, Accumulator: true},
			wantState:  State{A: 0x00, P: FlagCarry | FlagZero},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
		},
		{
			name:       "ROR 0x80 in accumulator, sets no flags.",
			addressing: Addressing{Value: 0x80, Accumulator: true},
			wantState:  State{A: 0x40},
		},
		{
			name:       "ROR 0x81 in accumulator, sets carry flag only.",
			addressing: Addressing{Value: 0x81, Accumulator: true},
			wantState:  State{A: 0x40, P: FlagCarry},
		},
		{
			name:       "ROR 0x40 in accumulator with carry, sets negative flag only.",
			startState: State{P: FlagCarry},
			addressing: Addressing{Value: 0x40, Accumulator: true},
			wantState:  State{A: 0xA0, P: FlagNegative},
		},
		{
			name:       "ROR 0x01 in accumulator does not change other settings, sets carry flag and clears negative flag only.",
			startState: State{PC: 0xFEAD, SP: 0x0F, A: 0x04, X: 0x02, Y: 0x03, P: 0xFF ^ FlagCarry},
			addressing: Addressing{Value: 0x01, Accumulator: true},
			wantState:  State{PC: 0xFEAD, SP: 0x0F, A: 0x00, X: 0x02, Y: 0x03, P: 0xFF ^ FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, RotateRight)
		})
	}
}
