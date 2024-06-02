package processor

import "testing"

func TestClearCarry(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Clear on Zero value does nothing.",
		},
		{
			name:       "Clear clears value.",
			startState: State{P: FlagCarry},
		},
		{
			name:       "Clear does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagCarry},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ClearCarry)
		})
	}
}

func TestClearDecimal(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Clear on Zero value does nothing.",
		},
		{
			name:       "Clear clears value.",
			startState: State{P: FlagDecimal},
		},
		{
			name:       "Clear does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagDecimal},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ClearDecimal)
		})
	}
}

func TestClearInterrupt(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Clear on Zero value does nothing.",
		},
		{
			name:       "Clear clears value.",
			startState: State{P: FlagInterrupt},
		},
		{
			name:       "Clear does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagInterrupt},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ClearInterrupt)
		})
	}
}

func TestClearOverflow(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Clear on Zero value does nothing.",
		},
		{
			name:       "Clear clears value.",
			startState: State{P: FlagOverflow},
		},
		{
			name:       "Clear does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagOverflow},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, ClearOverflow)
		})
	}
}

func TestSetCarry(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Set on Zero value sets carry.",
			wantState: State{P: FlagCarry},
		},
		{
			name:       "Set carry when already set does nothing.",
			startState: State{P: FlagCarry},
			wantState:  State{P: FlagCarry},
		},
		{
			name:       "Set does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagCarry},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, SetCarry)
		})
	}
}

func TestSetDecimal(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Set on Zero value sets decimal.",
			wantState: State{P: FlagDecimal},
		},
		{
			name:       "Set decimal when already set does nothing.",
			startState: State{P: FlagDecimal},
			wantState:  State{P: FlagDecimal},
		},
		{
			name:       "Set does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagDecimal},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, SetDecimal)
		})
	}
}

func TestSetInterrupt(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Set on Zero value sets interrupt.",
			wantState: State{P: FlagInterrupt},
		},
		{
			name:       "Set interrupt when already set does nothing.",
			startState: State{P: FlagInterrupt},
			wantState:  State{P: FlagInterrupt},
		},
		{
			name:       "Set does not change any other values.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF ^ FlagInterrupt},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40, P: 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, SetInterrupt)
		})
	}
}
