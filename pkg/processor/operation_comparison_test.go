package processor

import "testing"

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		start    Status
		register uint8
		value    uint8
		expect   Status
	}{
		{
			name:   "Everything is zero.",
			expect: FlagCarry | FlagZero,
		},
		{
			name:   "Everything is 0xFF.",
			expect: FlagCarry | FlagZero,
		},
		{
			name:     "Carry and Zero set when equal.",
			register: 0x32,
			value:    0x32,
			expect:   FlagCarry | FlagZero,
		},
		{
			name:     "Carry set when register is greater but not zero.",
			register: 0x32,
			value:    0x31,
			expect:   FlagCarry,
		},
		{
			name:     "Carry cleared when register is less and negative is set.",
			register: 0x31,
			value:    0x32,
			expect:   FlagNegative,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := tt.start
			compare(&status, tt.register, tt.value)

			if status != tt.expect {
				t.Errorf("compare() got = %v, expect %v", status.String(), tt.expect.String())
			}
		})
	}
}

func TestCompareWithA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Compare with A and ignore other values; match.",
			startState: State{PC: 0x1234, SP: 0x31, A: 0x20, X: 0x31, Y: 0x31, P: 0x0},
			addressing: Addressing{Value: 0x20},
			wantState:  State{PC: 0x1234, SP: 0x31, A: 0x20, X: 0x31, Y: 0x31, P: FlagZero | FlagCarry},
		},
		{
			name:       "Compare with A and ignore other values; no match.",
			startState: State{PC: 0x1234, SP: 0x31, A: 0x20, X: 0x31, Y: 0x31, P: 0x0},
			addressing: Addressing{Value: 0x31},
			wantState:  State{PC: 0x1234, SP: 0x31, A: 0x20, X: 0x31, Y: 0x31, P: FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, CompareWithA)
		})
	}
}

func TestCompareWithX(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Compare with X and ignore other values; match.",
			startState: State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x20, Y: 0x31, P: 0x0},
			addressing: Addressing{Value: 0x20},
			wantState:  State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x20, Y: 0x31, P: FlagZero | FlagCarry},
		},
		{
			name:       "Compare with X and ignore other values; no match.",
			startState: State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x20, Y: 0x31, P: 0x0},
			addressing: Addressing{Value: 0x31},
			wantState:  State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x20, Y: 0x31, P: FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, CompareWithX)
		})
	}
}

func TestCompareWithY(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Compare with Y and ignore other values; match.",
			startState: State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x31, Y: 0x20, P: 0x0},
			addressing: Addressing{Value: 0x20},
			wantState:  State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x31, Y: 0x20, P: FlagZero | FlagCarry},
		},
		{
			name:       "Compare with Y and ignore other values; no match.",
			startState: State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x31, Y: 0x20, P: 0x0},
			addressing: Addressing{Value: 0x31},
			wantState:  State{PC: 0x1234, SP: 0x31, A: 0x31, X: 0x31, Y: 0x20, P: FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, CompareWithY)
		})
	}
}
