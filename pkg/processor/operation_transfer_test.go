package processor

import "testing"

func TestLoadA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Load zero.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Load 0x01.",
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x01, P: FlagNone},
		},
		{
			name:       "Load 0x7E.",
			addressing: Addressing{Value: 0x7E},
			wantState:  State{A: 0x7E, P: FlagNone},
		},
		{
			name:       "Load 0x80.",
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x80, P: FlagNegative},
		},
		{
			name:       "Load 0x81.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{A: 0x81, P: FlagNegative},
		},
		{
			name:       "Load value; validate PC, SP, X and Y are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40},
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagNegative, PC: 0x1234, SP: 0x10, A: 0x81, X: 0x30, Y: 0x40},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, LoadA)
		})
	}
}

func TestLoadX(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Load zero.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Load 0x01.",
			addressing: Addressing{Value: 0x01},
			wantState:  State{X: 0x01, P: FlagNone},
		},
		{
			name:       "Load 0x7E.",
			addressing: Addressing{Value: 0x7E},
			wantState:  State{X: 0x7E, P: FlagNone},
		},
		{
			name:       "Load 0x80.",
			addressing: Addressing{Value: 0x80},
			wantState:  State{X: 0x80, P: FlagNegative},
		},
		{
			name:       "Load 0x81.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{X: 0x81, P: FlagNegative},
		},
		{
			name:       "Load value; validate PC, SP, A and Y are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40},
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagNegative, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x81, Y: 0x40},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, LoadX)
		})
	}
}

func TestLoadY(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Load zero.",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Load 0x01.",
			addressing: Addressing{Value: 0x01},
			wantState:  State{Y: 0x01, P: FlagNone},
		},
		{
			name:       "Load 0x7E.",
			addressing: Addressing{Value: 0x7E},
			wantState:  State{Y: 0x7E, P: FlagNone},
		},
		{
			name:       "Load 0x80.",
			addressing: Addressing{Value: 0x80},
			wantState:  State{Y: 0x80, P: FlagNegative},
		},
		{
			name:       "Load 0x81.",
			addressing: Addressing{Value: 0x81},
			wantState:  State{Y: 0x81, P: FlagNegative},
		},
		{
			name:       "Load value; validate PC, SP, A and X are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x40},
			addressing: Addressing{Value: 0x81},
			wantState:  State{P: FlagNegative, PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x81},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, LoadY)
		})
	}
}

func TestStoreA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Store zero",
		},
		{
			name:       "Store 0x01.",
			startState: State{A: 0x01, P: FlagNone},
			wantState:  State{A: 0x01, P: FlagNone},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Store 0x80.",
			startState: State{A: 0x80, P: FlagAll},
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{A: 0x80, P: FlagAll},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Store 0x81.",
			startState: State{PC: 0xFFFF, A: 0x81, P: FlagZero | FlagNegative | FlagDecimal},
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{EffectiveAddress: 0x03},
			wantState:  State{PC: 0xFFFF, A: 0x81, P: FlagZero | FlagNegative | FlagDecimal},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x81, 0x50, 0x60, 0x70, 0x80},
		},
		{
			name:       "Store value; validate PC, SP, A, X and Y are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x81, X: 0x30, Y: 0x40, P: FlagAll},
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{EffectiveAddress: 0x04},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x81, X: 0x30, Y: 0x40, P: FlagAll},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x40, 0x81, 0x60, 0x70, 0x80},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, StoreA)
		})
	}
}

func TestStoreX(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Store zero",
		},
		{
			name:       "Store 0x01.",
			startState: State{X: 0x01, P: FlagNone},
			wantState:  State{X: 0x01, P: FlagNone},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Store 0x80.",
			startState: State{X: 0x80, P: FlagAll},
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{X: 0x80, P: FlagAll},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Store 0x81.",
			startState: State{PC: 0xFFFF, X: 0x81, P: FlagZero | FlagNegative | FlagDecimal},
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{EffectiveAddress: 0x03},
			wantState:  State{PC: 0xFFFF, X: 0x81, P: FlagZero | FlagNegative | FlagDecimal},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x81, 0x50, 0x60, 0x70, 0x80},
		},
		{
			name:       "Store value; validate PC, SP, A, X and Y are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x81, Y: 0x40, P: FlagAll},
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{EffectiveAddress: 0x04},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x81, Y: 0x40, P: FlagAll},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x40, 0x81, 0x60, 0x70, 0x80},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, StoreX)
		})
	}
}

func TestStoreY(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Store zero",
		},
		{
			name:       "Store 0x01.",
			startState: State{Y: 0x01, P: FlagNone},
			wantState:  State{Y: 0x01, P: FlagNone},
			wantRam:    []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Store 0x80.",
			startState: State{Y: 0x80, P: FlagAll},
			startRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{Y: 0x80, P: FlagAll},
			wantRam:    []uint8{0x80, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "Store 0x81.",
			startState: State{PC: 0xFFFF, Y: 0x81, P: FlagZero | FlagNegative | FlagDecimal},
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{EffectiveAddress: 0x03},
			wantState:  State{PC: 0xFFFF, Y: 0x81, P: FlagZero | FlagNegative | FlagDecimal},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x81, 0x50, 0x60, 0x70, 0x80},
		},
		{
			name:       "Store value; validate PC, SP, A, X and Y are not touched.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x81, P: FlagAll},
			startRam:   []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
			addressing: Addressing{EffectiveAddress: 0x04},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x81, P: FlagAll},
			wantRam:    []uint8{0x10, 0x20, 0x30, 0x40, 0x81, 0x60, 0x70, 0x80},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, StoreY)
		})
	}
}

func TestTransferAtoX(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Transfer zero sets zero flag",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Transfer 0x80 sets sign flag",
			startState: State{A: 0x80},
			wantState:  State{A: 0x80, X: 0x80, P: FlagNegative},
		},
		{
			name:       "Transfer 0x01 clears sign and zero flag",
			startState: State{A: 0x01, P: FlagAll},
			wantState:  State{A: 0x01, X: 0x01, P: FlagAll ^ FlagZero ^ FlagNegative},
		},
		{
			name:       "Transfer does not affect PC, SP, A and Y.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x82, X: 0x30, Y: 0x40, P: FlagAll},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x82, X: 0x82, Y: 0x40, P: FlagAll ^ FlagZero},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, TransferAtoX)
		})
	}
}

func TestTransferAtoY(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Transfer zero sets zero flag",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Transfer 0x80 sets sign flag",
			startState: State{A: 0x80},
			wantState:  State{A: 0x80, Y: 0x80, P: FlagNegative},
		},
		{
			name:       "Transfer 0x01 clears sign and zero flag",
			startState: State{A: 0x01, P: FlagAll},
			wantState:  State{A: 0x01, Y: 0x01, P: FlagAll ^ FlagZero ^ FlagNegative},
		},
		{
			name:       "Transfer does not affect PC, SP, A and X.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x82, X: 0x30, Y: 0x40, P: FlagAll},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x82, X: 0x30, Y: 0x82, P: FlagAll ^ FlagZero},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, TransferAtoY)
		})
	}
}

func TestTransferSPtoX(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Transfer zero sets zero flag",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Transfer 0x80 sets sign flag",
			startState: State{SP: 0x80},
			wantState:  State{SP: 0x80, X: 0x80, P: FlagNegative},
		},
		{
			name:       "Transfer 0x01 clears sign and zero flag",
			startState: State{SP: 0x01, P: FlagAll},
			wantState:  State{SP: 0x01, X: 0x01, P: FlagAll ^ FlagZero ^ FlagNegative},
		},
		{
			name:       "Transfer does not affect PC, SP, A and Y.",
			startState: State{PC: 0x1234, SP: 0x82, A: 0x20, X: 0x30, Y: 0x40, P: FlagAll},
			wantState:  State{PC: 0x1234, SP: 0x82, A: 0x20, X: 0x82, Y: 0x40, P: FlagAll ^ FlagZero},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, TransferSPtoX)
		})
	}
}

func TestTransferXtoA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Transfer zero sets zero flag",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Transfer 0x80 sets sign flag",
			startState: State{X: 0x80},
			wantState:  State{X: 0x80, A: 0x80, P: FlagNegative},
		},
		{
			name:       "Transfer 0x01 clears sign and zero flag",
			startState: State{X: 0x01, P: FlagAll},
			wantState:  State{X: 0x01, A: 0x01, P: FlagAll ^ FlagZero ^ FlagNegative},
		},
		{
			name:       "Transfer does not affect PC, SP, X and Y.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x82, Y: 0x40, P: FlagAll},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x82, X: 0x82, Y: 0x40, P: FlagAll ^ FlagZero},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, TransferXtoA)
		})
	}
}

func TestTransferXtoSP(t *testing.T) {
	tests := []testOperationConfig{
		{
			name: "Transfer zero does nothing",
		},
		{
			name:       "Transfer 0x80 does not set sign flag",
			startState: State{X: 0x80},
			wantState:  State{X: 0x80, SP: 0x80},
		},
		{
			name:       "Transfer 0x01 does not clear sign and zero flags",
			startState: State{X: 0x01, P: FlagAll},
			wantState:  State{X: 0x01, SP: 0x01, P: FlagAll},
		},
		{
			name:       "Transfer does not affect PC, A, X and Y.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x82, Y: 0x40, P: FlagAll},
			wantState:  State{PC: 0x1234, SP: 0x82, A: 0x20, X: 0x82, Y: 0x40, P: FlagAll},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, TransferXtoSP)
		})
	}
}

func TestTransferYtoA(t *testing.T) {
	tests := []testOperationConfig{
		{
			name:      "Transfer zero sets zero flag",
			wantState: State{P: FlagZero},
		},
		{
			name:       "Transfer 0x80 sets sign flag",
			startState: State{Y: 0x80},
			wantState:  State{Y: 0x80, A: 0x80, P: FlagNegative},
		},
		{
			name:       "Transfer 0x01 clears sign and zero flag",
			startState: State{Y: 0x01, P: FlagAll},
			wantState:  State{Y: 0x01, A: 0x01, P: FlagAll ^ FlagZero ^ FlagNegative},
		},
		{
			name:       "Transfer does not affect PC, SP, X and Y.",
			startState: State{PC: 0x1234, SP: 0x10, A: 0x20, X: 0x30, Y: 0x82, P: FlagAll},
			wantState:  State{PC: 0x1234, SP: 0x10, A: 0x82, X: 0x30, Y: 0x82, P: FlagAll ^ FlagZero},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, TransferYtoA)
		})
	}
}
