package processor

import "testing"

func TestDecrement(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Decrement when Memory is zero.",
			wantState: State{P: FlagNegative},
			wantRam:   []uint8{0xFF, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Decrement when Memory is 1.",
			addressing: Addressing{Value: 0x01},
			wantState:  State{P: FlagZero},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Decrement when Memory is 0x80; sign should not be set.",
			startState: State{P: FlagNegative},
			addressing: Addressing{Value: 0x80},
			wantState:  State{P: FlagNone},
			wantRam:    []uint8{0x7F, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Decrement when Memory is 0x81; sign should be set.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Decrement using different effective address values 0x0001.",
			addressing: Addressing{EffectiveAddress: 0x0001, Value: 0x51},
			wantRam:    []uint8{0, 0x50, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Decrement using different effective address values 0x0002.",
			addressing: Addressing{EffectiveAddress: 0x0002, Value: 0x51},
			wantRam:    []uint8{0, 0, 0x50, 0, 0, 0, 0, 0},
		},
		{
			name:       "Decrement using different effective address values 0xFFFF.",
			addressing: Addressing{EffectiveAddress: 0xFFFF, Value: 0x51},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x50},
		},
		{
			name:       "Decrement when Memory is zero; validate PC, SP, Y and A are not touched..",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40},
			wantState:  State{P: FlagNegative, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40},
			wantRam:    []uint8{0xFF, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, Decrement)
		})
	}
}

func TestDecrementX(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Decrement when X is zero.",
			startState: State{X: 0x00},
			wantState:  State{X: 0xFF, P: FlagNegative},
		},
		{
			name:       "Decrement when X is 1.",
			startState: State{X: 0x01},
			wantState:  State{X: 0x0, P: FlagZero},
		},
		{
			name:       "Decrement when X is 0x80; sign should not be set.",
			startState: State{X: 0x80, P: FlagNegative},
			wantState:  State{X: 0x7F, P: 0},
		},
		{
			name:       "Decrement when X is 0x81; sign should be set.",
			startState: State{X: 0x81},
			wantState:  State{X: 0x80, P: FlagNegative},
		},
		{
			name:       "Decrement when X is zero; validate PC, SP, Y and A are not touched.",
			startState: State{X: 0x00, PC: 0x1234, SP: 0x10, A: 0x20, Y: 0x30},
			wantState:  State{X: 0xFF, P: FlagNegative, PC: 0x1234, SP: 0x10, A: 0x20, Y: 0x30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, DecrementX)
		})
	}
}

func TestDecrementY(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Decrement when Y is zero.",
			startState: State{Y: 0x00},
			wantState:  State{Y: 0xFF, P: FlagNegative},
		},
		{
			name:       "Decrement when Y is 1.",
			startState: State{Y: 0x01},
			wantState:  State{Y: 0x0, P: FlagZero},
		},
		{
			name:       "Decrement when Y is 0x80; sign should not be set.",
			startState: State{Y: 0x80, P: FlagNegative},
			wantState:  State{Y: 0x7F, P: 0},
		},
		{
			name:       "Decrement when Y is 0x81; sign should be set.",
			startState: State{Y: 0x81},
			wantState:  State{Y: 0x80, P: FlagNegative},
		},
		{
			name:       "Decrement when Y is zero; validate PC, SP, X and A are not touched.",
			startState: State{Y: 0x00, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30},
			wantState:  State{Y: 0xFF, P: FlagNegative, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, DecrementY)
		})
	}
}

func TestIncrement(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:    "Increment when Memory is zero.",
			wantRam: []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Increment when Memory is 0xFF.",
			addressing: Addressing{Value: 0xFF},
			wantState:  State{P: FlagZero},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Increment when Memory is 0x7E; sign should not be set.",
			startState: State{P: FlagNegative},
			addressing: Addressing{Value: 0x7E},
			wantState:  State{P: 0},
			wantRam:    []uint8{0x7F, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Increment when Memory is 0x7F; sign should be set.",
			addressing: Addressing{Value: 0x7F},
			wantState:  State{P: FlagNegative},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Increment using different effective address values 0x0001.",
			addressing: Addressing{EffectiveAddress: 0x0001, Value: 0x51},
			wantRam:    []uint8{0, 0x52, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Increment using different effective address values 0x0002.",
			addressing: Addressing{EffectiveAddress: 0x0002, Value: 0x51},
			wantRam:    []uint8{0, 0, 0x52, 0, 0, 0, 0, 0},
		},
		{
			name:       "Increment using different effective address values 0xFFFF.",
			addressing: Addressing{EffectiveAddress: 0xFFFF, Value: 0x51},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x52},
		},
		{
			name:       "Increment when Memory is zero; validate PC, SP, Y and A are not touched..",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40},
			wantState:  State{P: FlagNone, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, Increment)
		})
	}
}

func TestIncrementX(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Increment when X is zero.",
			startState: State{X: 0x00},
			wantState:  State{X: 0x01, P: FlagNone},
		},
		{
			name:       "Increment when X is 0xFF.",
			startState: State{X: 0xFF},
			wantState:  State{X: 0x00, P: FlagZero},
		},
		{
			name:       "Increment when X is 0x7E; sign should not be set.",
			startState: State{X: 0x7E, P: FlagNegative},
			wantState:  State{X: 0x7F, P: 0},
		},
		{
			name:       "Increment when X is 0x7F; sign should be set.",
			startState: State{X: 0x7F},
			wantState:  State{X: 0x80, P: FlagNegative},
		},
		{
			name:       "Increment when X is 0xFF; validate PC, SP, Y and A are not touched.",
			startState: State{X: 0xFF, PC: 0x1234, SP: 0x10, A: 0x20, Y: 0x30},
			wantState:  State{X: 0x00, P: FlagZero, PC: 0x1234, SP: 0x10, A: 0x20, Y: 0x30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, IncrementX)
		})
	}
}

func TestIncrementY(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:       "Increment when Y is zero.",
			startState: State{Y: 0x00},
			wantState:  State{Y: 0x01, P: FlagNone},
		},
		{
			name:       "Increment when Y is 0xFF.",
			startState: State{Y: 0xFF},
			wantState:  State{Y: 0x00, P: FlagZero},
		},
		{
			name:       "Increment when Y is 0x7E; sign should not be set.",
			startState: State{Y: 0x7E, P: FlagNegative},
			wantState:  State{Y: 0x7F, P: 0},
		},
		{
			name:       "Increment when Y is 0x7F; sign should be set.",
			startState: State{Y: 0x7F},
			wantState:  State{Y: 0x80, P: FlagNegative},
		},
		{
			name:       "Increment when Y is 0xFF; validate PC, SP, X and A are not touched.",
			startState: State{Y: 0xFF, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30},
			wantState:  State{Y: 0x00, P: FlagZero, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, IncrementY)
		})
	}
}
