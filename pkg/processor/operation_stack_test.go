package processor

import "testing"

func TestPushA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "PHA 0x00",
			startState: State{A: 0x00, SP: 0xFF},
			wantState:  State{A: 0x00, SP: 0xFE},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "PHA 0x01",
			startState: State{A: 0x01, SP: 0xFF},
			wantState:  State{A: 0x01, SP: 0xFE},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x01},
		},
		{
			name:       "PHA 0x10",
			startState: State{A: 0x10, SP: 0xFE},
			wantState:  State{A: 0x10, SP: 0xFD},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x10, 0},
		},
		{
			name:       "PHA 0xAC, PC, X, Y and Flags are not touched",
			startState: State{A: 0xAC, SP: 0xFD, PC: 0x1234, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
			wantState:  State{A: 0xAC, SP: 0xFC, PC: 0x1234, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0xAC, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, PushA)
		})
	}
}

func TestPushP(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "PHP zero always writes constant and break",
			startState: State{SP: 0xFF, P: 0x00},
			wantState:  State{SP: 0xFE, P: 0x00},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant},
		},
		{
			name:       "PHP two flags combines with constant and break",
			startState: State{SP: 0xFF, P: FlagNegative | FlagCarry},
			wantState:  State{SP: 0xFE, P: FlagNegative | FlagCarry},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagNegative | FlagCarry},
		},
		{
			name:       "PHP three flags combines with constant and break",
			startState: State{SP: 0xFE, P: FlagZero | FlagDecimal | FlagOverflow},
			wantState:  State{SP: 0xFD, P: FlagZero | FlagDecimal | FlagOverflow},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagZero | FlagDecimal | FlagOverflow, 0},
		},
		{
			name:       "PHP break and interrupt, PC, A, X, Y and Flags are not touched",
			startState: State{SP: 0xFD, PC: 0x1234, A: 0x10, X: 0x30, Y: 0x40, P: FlagBreak | FlagInterrupt},
			wantState:  State{SP: 0xFC, PC: 0x1234, A: 0x10, X: 0x30, Y: 0x40, P: FlagBreak | FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagInterrupt, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, PushP)
		})
	}
}

func TestPullA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "PLA 0x00 sets zero flag",
			startState: State{A: 0x00, SP: 0xFE},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{A: 0x00, SP: 0xFF, P: FlagZero},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "PLA 0x81 sets sign flag",
			startState: State{A: 0x00, SP: 0xFE},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0x81},
			wantState:  State{A: 0x81, SP: 0xFF, P: FlagNegative},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x81},
		},
		{
			name:       "PLA 0x00 sets zero flag clears negative flag",
			startState: State{A: 0x00, SP: 0xFE, P: FlagNegative},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{A: 0x00, SP: 0xFF, P: FlagZero},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "PLA 0x81 sets sign flag, clears zero flag",
			startState: State{A: 0x00, SP: 0xFE, P: FlagZero},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0x81},
			wantState:  State{A: 0x81, SP: 0xFF, P: FlagNegative},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x81},
		},
		{
			name:       "PLA 0x10, clears zero and negative flags",
			startState: State{A: 0xFF, SP: 0xFD, P: FlagZero | FlagNegative},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x10, 0},
			wantState:  State{A: 0x10, SP: 0xFE},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x10, 0},
		},
		{
			name:       "PLA 0x10, clears zero and negative flags but ignores all others",
			startState: State{A: 0xFF, SP: 0xFD, P: FlagZero | FlagNegative | FlagCarry | FlagOverflow | FlagDecimal | FlagInterrupt | FlagConstant},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x10, 0},
			wantState:  State{A: 0x10, SP: 0xFE, P: FlagCarry | FlagOverflow | FlagDecimal | FlagInterrupt | FlagConstant},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x10, 0},
		},
		{
			name:       "PLA 0xAC, PC, X, and Y are not touched",
			startState: State{A: 0x13, SP: 0xFC, PC: 0x1234, X: 0x30, Y: 0x40, P: FlagZero | FlagNegative | FlagCarry | FlagDecimal},
			startRam:   []uint8{0, 0, 0, 0, 0, 0xAC, 0, 0},
			wantState:  State{A: 0xAC, SP: 0xFD, PC: 0x1234, X: 0x30, Y: 0x40, P: FlagNegative | FlagCarry | FlagDecimal},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0xAC, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, PullA)
		})
	}
}

func TestPullP(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "PLP 0x00 will set constant",
			startState: State{SP: 0xFE, P: 0x00},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{SP: 0xFF, P: FlagConstant},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "PLP ignores break",
			startState: State{SP: 0xFE, P: 0x00},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant},
			wantState:  State{SP: 0xFF, P: FlagConstant},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant},
		},
		{
			name:       "PLP sets three flags and clears others",
			startState: State{SP: 0xFE, P: 0xFF},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagNegative | FlagCarry},
			wantState:  State{SP: 0xFF, P: FlagConstant | FlagNegative | FlagCarry},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagNegative | FlagCarry},
		},
		{
			name:       "PLP sets four flags and clears others",
			startState: State{SP: 0xFD, P: FlagNegative | FlagCarry | FlagOverflow},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagZero | FlagDecimal | FlagOverflow, 0},
			wantState:  State{SP: 0xFE, P: FlagConstant | FlagZero | FlagDecimal | FlagOverflow},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagZero | FlagDecimal | FlagOverflow, 0},
		},
		{
			name:       "PLP sets status, PC, A, X, and Y are not touched",
			startState: State{SP: 0xFC, PC: 0x1234, A: 0x10, X: 0x30, Y: 0x40, P: 0x00},
			startRam:   []uint8{0, 0, 0, 0, 0, FlagBreak | FlagInterrupt, 0, 0},
			wantState:  State{SP: 0xFD, PC: 0x1234, A: 0x10, X: 0x30, Y: 0x40, P: FlagConstant | FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, 0, 0, FlagBreak | FlagInterrupt, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, PullP)
		})
	}
}
