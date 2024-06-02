package processor

import "testing"

func TestBranchOnCarryClear(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Carry flag clear causes branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagCarry},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: 0xFF ^ FlagCarry},
		},
		{
			name:       "Carry flag set does not branch.",
			startState: State{PC: 0x0FFF, P: FlagCarry},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: FlagCarry},
		},
		{
			name:       "Branching does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagCarry},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagCarry},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnCarryClear)
		})
	}
}

func TestBranchOnCarrySet(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Carry flag set causes branch.",
			startState: State{PC: 0x0FFF, P: FlagCarry},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: FlagCarry},
		},
		{
			name:       "Carry flag clear does not branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagCarry},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: 0xFF ^ FlagCarry},
		},
		{
			name:       "Branching does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnCarrySet)
		})
	}
}

func TestBranchOnEqual(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Zero flag set causes branch.",
			startState: State{PC: 0x0FFF, P: FlagZero},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: FlagZero},
		},
		{
			name:       "Zero flag clear does not branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagZero},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: 0xFF ^ FlagZero},
		},
		{
			name:       "Branching does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnEqual)
		})
	}
}

func TestBranchOnMinus(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Negative flag set causes branch.",
			startState: State{PC: 0x0FFF, P: FlagNegative},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: FlagNegative},
		},
		{
			name:       "Negative flag clear does not branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagNegative},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: 0xFF ^ FlagNegative},
		},
		{
			name:       "Negative does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnMinus)
		})
	}
}

func TestBranchOnNotEqual(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Zero flag clear causes branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagZero},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: 0xFF ^ FlagZero},
		},
		{
			name:       "Zero flag set does not branch.",
			startState: State{PC: 0x0FFF, P: FlagZero},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: FlagZero},
		},
		{
			name:       "Branching does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagZero},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagZero},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnNotEqual)
		})
	}
}
func TestBranchOnOverflowClear(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Overflow flag clear causes branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagOverflow},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: 0xFF ^ FlagOverflow},
		},
		{
			name:       "Overflow flag set does not branch.",
			startState: State{PC: 0x0FFF, P: FlagOverflow},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: FlagOverflow},
		},
		{
			name:       "Branching does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagOverflow},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagOverflow},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnOverflowClear)
		})
	}
}
func TestBranchOnOverflowSet(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Overflow flag set causes branch.",
			startState: State{PC: 0x0FFF, P: FlagOverflow},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: FlagOverflow},
		},
		{
			name:       "Overflow flag clear does not branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagOverflow},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: 0xFF ^ FlagOverflow},
		},
		{
			name:       "Overflow does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnOverflowSet)
		})
	}
}
func TestBranchOnPlus(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Negative flag clear causes branch.",
			startState: State{PC: 0x0FFF, P: 0xFF ^ FlagNegative},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, P: 0xFF ^ FlagNegative},
		},
		{
			name:       "Negative flag set does not branch.",
			startState: State{PC: 0x0FFF, P: FlagNegative},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0x0FFF, P: FlagNegative},
		},
		{
			name:       "Branching does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagNegative},
			addressing: Addressing{EffectiveAddress: 0xFF00},
			wantState:  State{PC: 0xFF00, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, BranchOnPlus)
		})
	}
}
