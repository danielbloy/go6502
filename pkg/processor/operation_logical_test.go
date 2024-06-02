package processor

import "testing"

func TestAndWithAccumulator(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "AND zero to zero, set zero flag.",
			addressing: Addressing{},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "AND 0xFF to 0xFF, set negative flag.",
			startState: State{A: 0xFF},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{
			name:       "AND two mutually exclusive values, sets zero flag.",
			startState: State{A: 0x27},
			addressing: Addressing{Value: 0xFF ^ 0x27},
			wantState:  State{A: 0x00, P: FlagZero},
		},
		{
			name:       "AND two positive values.",
			startState: State{A: 0x27},
			addressing: Addressing{Value: 0x12},
			wantState:  State{A: 0x02},
		},
		{
			name:       "AND two negative values.",
			startState: State{A: 0x87},
			addressing: Addressing{Value: 0x92},
			wantState:  State{A: 0x82, P: FlagNegative},
		},
		{
			name:       "Simple positive values, checking nothing else changes.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x57, X: 0x30, Y: 0x40, P: 0xFF},
			addressing: Addressing{EffectiveAddress: 0xABCD, Value: 0x66},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x46, X: 0x30, Y: 0x40, P: (0xFF ^ FlagZero) ^ FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, AndWithA)
		})
	}
}

func TestExclusiveOrWithA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "XOR zero to zero, sets zero flag.",
			addressing: Addressing{},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "XOR zero to 0xFF, set negative flag.",
			startState: State{A: 0xFF},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{
			name:       "XOR 0xF0 to 0x0F, set negative flag.",
			startState: State{A: 0xF0},
			addressing: Addressing{Value: 0x0F},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{
			name:       "XOR two mutually exclusive values, sets negative flag.",
			startState: State{A: 0x27},
			addressing: Addressing{Value: 0xFF ^ 0x27},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{
			name:       "XOR two identical values, sets zero flag.",
			startState: State{A: 0x27},
			addressing: Addressing{Value: 0x27},
			wantState:  State{A: 0x00, P: FlagZero},
		},
		{
			name:       "XOR two positive values, clears Zero and Negative flags but not others.",
			startState: State{A: 0x27, P: FlagNegative | FlagZero | FlagCarry | FlagDecimal},
			addressing: Addressing{Value: 0x12},
			wantState:  State{A: 0x35, P: FlagCarry | FlagDecimal},
		},
		{
			name:       "XOR two negative values.",
			startState: State{A: 0x87},
			addressing: Addressing{Value: 0x92},
			wantState:  State{A: 0x15, P: FlagNone},
		},
		{
			name:       "Simple positive values, checking nothing else changes.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x57, X: 0x30, Y: 0x40, P: 0xFF},
			addressing: Addressing{EffectiveAddress: 0xABCD, Value: 0x66},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x31, X: 0x30, Y: 0x40, P: (0xFF ^ FlagZero) ^ FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ExclusiveOrWithA)
		})
	}
}

func TestOrWithA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "ORA zero to zero, sets zero flag.",
			addressing: Addressing{},
			wantState:  State{P: FlagZero},
		},
		{
			name:       "ORA zero to 0xFF, set negative flag.",
			startState: State{A: 0xFF},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{
			name:       "ORA 0xF0 to 0x0F, set negative flag.",
			startState: State{A: 0xF0},
			addressing: Addressing{Value: 0x0F},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{
			name:       "ORA two values with common bits set, clears negative flag.",
			startState: State{A: 0x27, P: FlagNegative},
			addressing: Addressing{Value: 0x63},
			wantState:  State{A: 0x67},
		},
		{
			name:       "ORA two identical values, clears zero flag.",
			startState: State{A: 0x27, P: FlagZero},
			addressing: Addressing{Value: 0x27},
			wantState:  State{A: 0x27},
		},
		{
			name:       "ORA two positive values, clears zero and negative flags.",
			startState: State{A: 0x27, P: FlagZero | FlagNegative},
			addressing: Addressing{Value: 0x12},
			wantState:  State{A: 0x37},
		},
		{
			name:       "ORA two negative values, clears zero and sets negative flag but ignores others.",
			startState: State{A: 0x82, P: FlagZero | FlagCarry | FlagDecimal},
			addressing: Addressing{Value: 0x97},
			wantState:  State{A: 0x97, P: FlagNegative | FlagCarry | FlagDecimal},
		},
		{
			name:       "Simple positive values, checking nothing else changes (except zero and negative flags cleared).",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x57, X: 0x30, Y: 0x40, P: 0xFF},
			addressing: Addressing{EffectiveAddress: 0xABCD, Value: 0x66},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x77, X: 0x30, Y: 0x40, P: 0xFF ^ (FlagZero | FlagNegative)},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, OrWithA)
		})
	}
}
