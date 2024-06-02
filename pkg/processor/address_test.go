package processor

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_MakeAddress(t *testing.T) {
	tests := []struct {
		name string
		low  uint8
		high uint8
		want Address
	}{
		{name: "All zeros"},
		{name: "Only a low byte", low: 0x13, want: 0x0013},
		{name: "Only a high byte", high: 0xAB, want: 0xAB00},
		{name: "Both a high and low byte", low: 0x42, high: 0x13, want: 0x1342},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeAddress(tt.low, tt.high); got != tt.want {
				t.Errorf("MakeAddress() = %v, want %v", got, tt.want)
			}

			// Do the reverse.
			if low, high := SplitAddress(tt.want); low != tt.low || high != tt.high {
				t.Errorf("SplitAddress() = %v:%v, want = %v:%v", low, high, tt.low, tt.high)
			}
		})
	}
}

func Test_SplitAddress(t *testing.T) {
	tests := []struct {
		name    string
		address Address
		low     uint8
		high    uint8
	}{
		{name: "All zeros"},
		{name: "Only a low byte", address: 0x000A, low: 0x0A},
		{name: "Only a high byte", address: 0xAB00, high: 0xAB},
		{name: "Both a high and low byte", address: 0x1342, low: 0x42, high: 0x13},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			low, high := SplitAddress(tt.address)
			if low != tt.low {
				t.Errorf("SplitAddress() low = %v, want %v", low, tt.low)
			}
			if high != tt.high {
				t.Errorf("SplitAddress() high = %v, want %v", high, tt.high)
			}

			// Do the reverse with MakeAddress
			if got := MakeAddress(tt.low, tt.high); got != tt.address {
				t.Errorf("MakeAddress() = %v, want %v", got, tt.address)
			}
		})
	}
}

func TestAddressing_Save(t *testing.T) {
	tests := []struct {
		name             string
		startRam         []uint8
		startState       State
		accumulator      bool
		effectiveAddress Address
		wantRam          []uint8
		wantState        State
		wantErr          bool
	}{
		{
			name:    "Error when memory is nil and it is memory write",
			wantErr: true,
		},
		{
			name:        "No error when memory is nil and it is accumulator write",
			accumulator: true,
			wantState:   State{A: 0x1E},
			wantErr:     false,
		},
		{
			name:        "Test accumulator write returns modified state but not modified memory",
			startRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState:  State{A: 0xF2},
			accumulator: true,
			wantRam:     []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState:   State{A: 0x1E},
			wantErr:     false,
		},
		{
			name:             "Test memory write returns modified memory but not modified accumulator",
			startRam:         []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState:       State{A: 0xF2},
			effectiveAddress: 0xFFF3,
			accumulator:      false,
			wantRam:          []uint8{0, 0, 0, 0x1E, 0, 0, 0, 0},
			wantState:        State{A: 0xF2},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var memory Memory
			var ram RepeatingRam
			if tt.startRam != nil {
				ram = NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
				memory = &ram
			}

			as := Addressing{
				Accumulator:      tt.accumulator,
				EffectiveAddress: tt.effectiveAddress,
				Memory:           memory,
			}
			got, err := as.Store(tt.startState, 0x1E)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantState) {
				t.Errorf("Store() got = %v, want %v", got, tt.wantState)
			}

			// Validate expected final RAM state
			if tt.startRam != nil && !reflect.DeepEqual(ram.ram, tt.wantRam) {
				t.Errorf("Store() RAM state got = %v, want = %v", ram.ram, tt.wantRam)
			}
		})
	}
}

func Push8bitValueTest(t *testing.T, name string, setConstant bool, callback func(Addressing, State, uint8) (State, error)) {

	modifier := uint8(0)
	if setConstant {
		modifier = FlagConstant
	}

	tests := []struct {
		name       string
		startRam   []uint8
		startState State
		wantRam    []uint8
		wantState  State
		wantErr    bool
	}{
		{
			name:    "Error when memory is nil and it is memory write",
			wantErr: true,
		},
		{
			name:       fmt.Sprintf("Test %v returns modified memory - a", name),
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0xFF},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x1E + modifier},
			wantState:  State{SP: 0xFE},
			wantErr:    false,
		},
		{
			name:       fmt.Sprintf("Test %v returns modified memory - b", name),
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0xFE},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x1E + modifier, 0},
			wantState:  State{SP: 0xFD},
			wantErr:    false,
		},
		{
			name:       fmt.Sprintf("Test %v wraps around", name),
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0x00},
			wantRam:    []uint8{0x1E + modifier, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{SP: 0xFF},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var memory Memory
			var ram RepeatingRam
			if tt.startRam != nil {
				ram = NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
				memory = &ram
			}

			as := Addressing{Memory: memory}
			got, err := callback(as, tt.startState, 0x1E)

			if (err != nil) != tt.wantErr {
				t.Errorf("%v error = %v, wantErr %v", name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantState) {
				t.Errorf("%v got = %v, want %v", name, got, tt.wantState)
			}

			// Validate expected final RAM state
			if tt.startRam != nil && !reflect.DeepEqual(ram.ram, tt.wantRam) {
				t.Errorf("%v RAM state got = %v, want = %v", name, ram.ram, tt.wantRam)
			}
		})
	}
}

func TestAddressing_PushByte(t *testing.T) {
	callback := func(addressing Addressing, state State, value uint8) (State, error) {
		return addressing.PushByte(state, value)
	}
	Push8bitValueTest(t, "PushByte", false, callback)
}

func TestAddressing_PushStatus(t *testing.T) {
	callback := func(addressing Addressing, state State, value uint8) (State, error) {
		return addressing.PushStatus(state, Status(value))
	}
	Push8bitValueTest(t, "PushStatus", true, callback)
}

func TestAddressing_PushAddress(t *testing.T) {
	tests := []struct {
		name       string
		startRam   []uint8
		startState State
		wantRam    []uint8
		wantState  State
		wantErr    bool
	}{
		{
			name:    "Error when memory is nil and it is memory write",
			wantErr: true,
		},
		{
			name:       "Test push address returns modified memory - a",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0xFF},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x1E, 0x2A},
			wantState:  State{SP: 0xFD},
			wantErr:    false,
		},
		{
			name:       "Test push address returns modified memory - b",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0xFE},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0x1E, 0x2A, 0},
			wantState:  State{SP: 0xFC},
			wantErr:    false,
		},
		{
			name:       "Test push address wraps around - a",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0x01},
			wantRam:    []uint8{0x1E, 0x2A, 0, 0, 0, 0, 0, 0},
			wantState:  State{SP: 0xFF},
			wantErr:    false,
		},
		{
			name:       "Test push address wraps around - b",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0x00},
			wantRam:    []uint8{0x2A, 0, 0, 0, 0, 0, 0, 0x1E},
			wantState:  State{SP: 0xFE},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var memory Memory
			var ram RepeatingRam
			if tt.startRam != nil {
				ram = NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
				memory = &ram
			}

			as := Addressing{Memory: memory}
			got, err := as.PushAddress(tt.startState, 0x2A1E)

			if (err != nil) != tt.wantErr {
				t.Errorf("PushAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantState) {
				t.Errorf("PushAddress() got = %v, want %v", got, tt.wantState)
			}

			// Validate expected final RAM state
			if tt.startRam != nil && !reflect.DeepEqual(ram.ram, tt.wantRam) {
				t.Errorf("PushAddress() RAM state got = %v, want = %v", ram.ram, tt.wantRam)
			}
		})
	}
}

func TestAddressing_PullByte(t *testing.T) {

	tests := []struct {
		name       string
		startRam   []uint8
		startState State
		wantValue  uint8
		wantRam    []uint8
		wantState  State
		wantErr    bool
	}{
		{
			name:    "PullByte error when memory is nil and it is memory read",
			wantErr: true,
		},
		{
			name:       "PullByte returns modified memory - a",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0x1E},
			startState: State{SP: 0xFE},
			wantValue:  0x1E,
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x1E},
			wantState:  State{SP: 0xFF},
			wantErr:    false,
		},
		{
			name:       "PullByte returns modified memory - b",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x1E, 0},
			startState: State{SP: 0xFD},
			wantValue:  0x1E,
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x1E, 0},
			wantState:  State{SP: 0xFE},
			wantErr:    false,
		},
		{
			name:       "PullByte wraps around",
			startRam:   []uint8{0x1E, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0xFF},
			wantValue:  0x1E,
			wantRam:    []uint8{0x1E, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{SP: 0x00},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var memory Memory
			var ram RepeatingRam
			if tt.startRam != nil {
				ram = NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
				memory = &ram
			}

			addressing := Addressing{Memory: memory}
			state, value, err := addressing.PullByte(tt.startState)

			if (err != nil) != tt.wantErr {
				t.Errorf("PullByte error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if value != tt.wantValue {
				t.Errorf("PullByte value got = %v, value want %v", value, tt.wantValue)
			}

			if !reflect.DeepEqual(state, tt.wantState) {
				t.Errorf("PullByte got = %v, want %v", state, tt.wantState)
			}

			// Validate expected final RAM tate
			if tt.startRam != nil && !reflect.DeepEqual(ram.ram, tt.wantRam) {
				t.Errorf("PullByte RAM state got = %v, want = %v", ram.ram, tt.wantRam)
			}
		})
	}
}

func TestAddressing_PullStatus(t *testing.T) {
	tests := []struct {
		name       string
		startRam   []uint8
		startState State
		wantRam    []uint8
		wantState  State
		wantErr    bool
	}{
		{
			name:    "PullStatus error when memory is nil and it is memory read",
			wantErr: true,
		},
		{
			name:       "PullStatus returns modified memory - 1",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, FlagDecimal | FlagOverflow},
			startState: State{SP: 0xFE},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, FlagDecimal | FlagOverflow},
			wantState:  State{SP: 0xFF, P: FlagConstant | FlagDecimal | FlagOverflow},
			wantErr:    false,
		},
		{
			name:       "PullStatus returns modified memory - 2",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, FlagBreak | FlagDecimal | FlagOverflow, 0},
			startState: State{SP: 0xFD},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, FlagBreak | FlagDecimal | FlagOverflow, 0},
			wantState:  State{SP: 0xFE, P: FlagConstant | FlagDecimal | FlagOverflow},
			wantErr:    false,
		},
		{
			name:       "PullStatus wraps around",
			startRam:   []uint8{FlagBreak | FlagDecimal | FlagOverflow, 0, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0xFF},
			wantRam:    []uint8{FlagBreak | FlagDecimal | FlagOverflow, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{SP: 0x00, P: FlagConstant | FlagDecimal | FlagOverflow},
			wantErr:    false,
		},
		{
			name:       "PullStatus will set constant",
			startState: State{SP: 0xFE, P: 0x00},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState:  State{SP: 0xFF, P: FlagConstant},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:       "PullStatus ignores break",
			startState: State{SP: 0xFE, P: 0x00},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant},
			wantState:  State{SP: 0xFF, P: FlagConstant},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant},
		},
		{
			name:       "PullStatus sets three flags and clears others",
			startState: State{SP: 0xFE, P: 0xFF},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagNegative | FlagCarry},
			wantState:  State{SP: 0xFF, P: FlagConstant | FlagNegative | FlagCarry},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagNegative | FlagCarry},
		},
		{
			name:       "PullStatus sets four flags and clears others",
			startState: State{SP: 0xFD, P: FlagNegative | FlagCarry | FlagOverflow},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagZero | FlagDecimal | FlagOverflow, 0},
			wantState:  State{SP: 0xFE, P: FlagConstant | FlagZero | FlagDecimal | FlagOverflow},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, FlagBreak | FlagConstant | FlagZero | FlagDecimal | FlagOverflow, 0},
		},
		{
			name:       "PullStatus sets status, PC, A, X, and Y are not touched",
			startState: State{SP: 0xFC, PC: 0x1234, A: 0x10, X: 0x30, Y: 0x40, P: 0x00},
			startRam:   []uint8{0, 0, 0, 0, 0, FlagBreak | FlagInterrupt, 0, 0},
			wantState:  State{SP: 0xFD, PC: 0x1234, A: 0x10, X: 0x30, Y: 0x40, P: FlagConstant | FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, 0, 0, FlagBreak | FlagInterrupt, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var memory Memory
			var ram RepeatingRam
			if tt.startRam != nil {
				ram = NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
				memory = &ram
			}

			addressing := Addressing{Memory: memory}
			state, err := addressing.PullStatus(tt.startState)

			if (err != nil) != tt.wantErr {
				t.Errorf("PullStatus error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(state, tt.wantState) {
				t.Errorf("PullStatus got = %v, want %v", state, tt.wantState)
			}

			// Validate expected final RAM state
			if tt.startRam != nil && !reflect.DeepEqual(ram.ram, tt.wantRam) {
				t.Errorf("PullStatus RAM state got = %v, want = %v", ram.ram, tt.wantRam)
			}
		})
	}
}

func TestAddressing_PullAddress(t *testing.T) {
	tests := []struct {
		name       string
		startRam   []uint8
		startState State
		wantRam    []uint8
		wantState  State
		wantErr    bool
	}{
		{
			name:    "Error when memory is nil and it is memory write",
			wantErr: true,
		},
		{
			name:       "Test pull address returns correct address - a",
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x1E, 0x2A},
			startState: State{SP: 0xFD},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0x1E, 0x2A},
			wantState:  State{SP: 0xFF},
			wantErr:    false,
		},
		{
			name:       "Test pull address returns correct address - b",
			startRam:   []uint8{0, 0, 0, 0, 0, 0x1E, 0x2A, 0},
			startState: State{SP: 0xFC},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0x1E, 0x2A, 0},
			wantState:  State{SP: 0xFE},
			wantErr:    false,
		},
		{
			name:       "Test pull address wraps around - a",
			startRam:   []uint8{0x1E, 0x2A, 0, 0, 0, 0, 0, 0},
			startState: State{SP: 0xFF},
			wantRam:    []uint8{0x1E, 0x2A, 0, 0, 0, 0, 0, 0},
			wantState:  State{SP: 0x01},
			wantErr:    false,
		},
		{
			name:       "Test pull address wraps around - b",
			startRam:   []uint8{0x2A, 0, 0, 0, 0, 0, 0, 0x1E},
			startState: State{SP: 0xFE},
			wantRam:    []uint8{0x2A, 0, 0, 0, 0, 0, 0, 0x1E},
			wantState:  State{SP: 0x00},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var memory Memory
			var ram RepeatingRam
			if tt.startRam != nil {
				ram = NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
				memory = &ram
			}

			as := Addressing{Memory: memory}
			state, address, err := as.PullAddress(tt.startState)

			if (err != nil) != tt.wantErr {
				t.Errorf("PullAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(state, tt.wantState) {
				t.Errorf("PullAddress() state got = %v, want %v", state, tt.wantState)
			}

			if err == nil && address != 0x2A1E {
				t.Errorf("PullAddress() address got = %04X, want %04X", address, 0x2A1E)
			}

			// Validate expected final RAM state
			if tt.startRam != nil && !reflect.DeepEqual(ram.ram, tt.wantRam) {
				t.Errorf("PullAddress() RAM state got = %v, want = %v", ram.ram, tt.wantRam)
			}
		})
	}
}

// The following struct and harness function are used to test the individual
// addressing functions.
type testAddressingConfig struct {
	name    string
	start   State
	ram     []uint8
	want    Addressing
	wantErr bool
}

func testAddressing(t *testing.T, tt testAddressingConfig, addressing AddressingFunc) {
	// Default the start RAM as empty if not supplied
	if len(tt.ram) <= 0 {
		tt.ram = []uint8{0, 0, 0, 0, 0, 0, 0, 0}
	}
	memory := NewPopulatedRam(RepeatingRamSize(len(tt.ram)), tt.ram)
	tt.want.Memory = &memory
	got, err := addressing(tt.start, &memory)

	if (err != nil) != tt.wantErr {
		t.Errorf("testAddressing() error = %v, wantErr %v", err, tt.wantErr)
		return
	}

	// Validate the final state.
	if !reflect.DeepEqual(got, tt.want) {
		t.Errorf("testAddressing() state got = %v, want = %v", got, tt.want)
	}

	// Validate the RAM has not changed
	if !reflect.DeepEqual(memory.ram, tt.ram) {
		t.Errorf("testAddressing() RAM state got = %v, want = %v", memory.ram, tt.ram)
	}

	// Test with nil Memory
	_, _ = addressing(tt.start, nil)
}

func TestAbsolute(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in absolute address of 0x0000",
			want: Addressing{ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and zero memory results in absolute address of 0x0000",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{ProgramCounterChange: 2},
		},
		{
			name: "Zero input State and non-zero memory results in absolute address of 0x0201 and Value of 0x02",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0201, Value: 0x02, ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and non-zero memory results in absolute address of 0xF5F6 and Value of 0xF2",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0xF5F6, Value: 0xF2, ProgramCounterChange: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, Absolute)
		})
	}
}

func TestAbsoluteX(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in absolute address of 0x0000",
			want: Addressing{ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and zero memory results in absolute address of X 0x0071",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{EffectiveAddress: 0x0071, ProgramCounterChange: 2},
		},
		{
			name: "Zero input State and non-zero memory results in absolute address of 0x0201 and Value of 0x02",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0201, Value: 0x02, ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and non-zero memory results in absolute address of 0xF5F6 + 0x11 (X) and Value of 0xF1",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0xF607, Value: 0xF1, ProgramCounterChange: 2, PageBoundaryCrossed: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, AbsoluteX)
		})
	}
}

func TestAbsoluteY(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in absolute address of 0x0000",
			want: Addressing{ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and zero memory results in absolute address of Y 0x00AB",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{EffectiveAddress: 0x00AB, ProgramCounterChange: 2},
		},
		{
			name: "Zero input State and non-zero memory results in absolute address of 0x0201 and Value of 0x02",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0201, Value: 0x02, ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and non-zero memory results in absolute address of 0xF5F6 + 0xFF (Y) and Value of 0xF3",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0xF6F5, Value: 0xF3, ProgramCounterChange: 2, PageBoundaryCrossed: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, AbsoluteY)
		})
	}
}

func TestAccumulator(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in zero result",
			want: Addressing{Accumulator: true},
		},
		{
			name:  "Non-zero input State and zero memory results in Accumulator result",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{Accumulator: true, Value: 0x3F},
		},
		{
			name: "Zero input State and non-zero memory results in zero result",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{Accumulator: true},
		},
		{
			name:  "Non-zero input State and non-zero memory results in Accumulator result",
			start: State{PC: 0x20, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{Accumulator: true, Value: 0x12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, Accumulator)
		})
	}
}

func TestImmediate(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in reading a zero",
			want: Addressing{ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and zero memory results effective address of PC and PC-delta of 1",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{EffectiveAddress: 0xF0, ProgramCounterChange: 1},
		},
		{
			name: "Zero input State and non-zero memory results in the memory Value at 0x00 being used",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{Value: 0x01, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and non-zero memory results in Value at PC being used",
			start: State{PC: 0x03, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0x03, Value: 0xF5, ProgramCounterChange: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, Immediate)
		})
	}
}

func TestImplied(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in zero result",
		},
		{
			name:  "Non-zero input State and zero memory results in zero result",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
		},
		{
			name: "Zero input State and non-zero memory results in zero result",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
		{
			name:  "Non-zero input State and non-zero memory results in zero result",
			start: State{PC: 0x20, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, Implied)
		})
	}
}

func TestIndirect(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in effective address of 0x0000 from 0x0000",
			want: Addressing{ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and zero memory results in effective address of 0x0000 from 0x0000",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{ProgramCounterChange: 2},
		},
		{
			name: "Zero input State and non-zero memory results in effective address of 0x0302 from 0x0201 and Value of 0x03",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0302, Value: 0x03, ProgramCounterChange: 2},
		},
		{
			name:  "Non-zero input State and non-zero memory results in effective address of 0xF1F2 from 0xF5F6 and Value of 0xF6",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0xF1F2, Value: 0xF6, ProgramCounterChange: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, Indirect)
		})
	}

	t.Run("Test page boundary wrap around where indirect $xxFF and $xxFF+1 resolves as $xxFF and $xx00 rather than $xxFF and $x(x+1)00", func(t *testing.T) {
		ram, err := NewRepeatingRam(OneKiloByte)
		if err != nil {
			panic(err)
		}

		data := map[Address]uint8{
			0x0000: 0xFF,
			0x0001: 0x01, // Indirect address = 0x01FF
			0x0100: 0x02, // Low byte wrap = 0x0100 = 0x02
			0x01FF: 0xAA, // Low byte      = 0x01FF = 0xAA
			0x0200: 0x03, // Low byte + 1  = 0x0200
			0x02AA: 0x0F,
			0x03AA: 0x0A,
		}
		err = WriteDataToMemory(&ram, data)
		if err != nil {
			panic(err)
		}

		got, _ := Indirect(State{}, &ram)
		want := Addressing{EffectiveAddress: 0x02AA, Value: 0x0F, ProgramCounterChange: 2, Memory: &ram}

		// Validate the final State.
		if !reflect.DeepEqual(got, want) {
			t.Errorf("TestIndirect() State got = %v, want = %v", got, want)
		}
	})
}

func TestIndirectX(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in effective address of 0x00 from 0x00",
			want: Addressing{ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and zero memory results in effective address of 0x00 from 0x00 + 0x71",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{ProgramCounterChange: 1},
		},
		{
			name: "Zero input State and non-zero memory results in effective address of 0x0302 from 0x01 + 0x00 and Value of 0x03",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0302, Value: 0x03, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and non-zero memory results in effective address of 0xF8F1 from 0xF6 + 0x11 ($107 cut to $07) and Value of 0xF7",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0xF8F1, Value: 0xF7, ProgramCounterChange: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, IndirectX)
		})
	}

	t.Run("Test zero page boundary wrap around", func(t *testing.T) {
		ram, err := NewRepeatingRam(OneKiloByte)
		if err != nil {
			panic(err)
		}

		data := map[Address]uint8{
			0x0000: 0xF6,
			0x0007: 0x02, // 0xF6 + 0x11 (X) = 0x07 with wraparound.
			0x0008: 0x01,
			0x0107: 0x04, // 0xF6 + 0x11 (X) = 0x107 without wraparound, but 0x07 with wraparound
			0x0108: 0x03,
			0x0102: 0xFF, // Expected value
			0x0304: 0x10, // Not wanted value
		}
		err = WriteDataToMemory(&ram, data)
		if err != nil {
			panic(err)
		}

		got, _ := IndirectX(State{X: 0x11}, &ram)
		want := Addressing{EffectiveAddress: 0x0102, Value: 0xFF, ProgramCounterChange: 1, Memory: &ram}

		// Validate the final State.
		if !reflect.DeepEqual(got, want) {
			t.Errorf("IndirectX() State got = %v, want = %v", got, want)
		}
	})
}

func TestIndirectY(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in effective address of 0x00 from 0x00 + 0x00",
			want: Addressing{ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and zero memory results in effective address of 0x00AB from 0x00 + 0xAB",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{EffectiveAddress: 0xAB, ProgramCounterChange: 1},
		},
		{
			name: "Zero input State and non-zero memory results in effective address of 0x0302 and Value 0x03; read byte 0x00 which has Value 0x01. Address is then read from 0x01 (0x02) and 0x02 (0x03)",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0302, Value: 0x03, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and non-zero memory results in effective address of 0xF2F1 and Value 0xF7; read byte 0x22 which has Value 0xF6. Address is then read from 0xF6 (0xF2) and 0xF6 (0xF1 - 0xF1F2) to which Y (0xFF) is added (0xF2F1)",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0xF2F1, Value: 0xF7, ProgramCounterChange: 1, PageBoundaryCrossed: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, IndirectY)
		})
	}

	t.Run("Test zero page boundary wrap around", func(t *testing.T) {
		ram, err := NewRepeatingRam(OneKiloByte)
		if err != nil {
			panic(err)
		}

		data := map[Address]uint8{
			0x0001: 0xFF, // 0xFF + 0x1 = 0x100 without wraparound, but 0x00 with wraparound
			0x00FF: 0x02, // Expected address (before Y is added): 0x0102
			0x0000: 0x01,
			0x0100: 0x04, // Unexpected address (before Y is added): 0x0402
			0x0113: 0xAA, // Expected value
			0x0413: 0x10, // Not wanted value
		}
		err = WriteDataToMemory(&ram, data)
		if err != nil {
			panic(err)
		}

		got, _ := IndirectY(State{PC: 0x0001, Y: 0x11}, &ram)
		want := Addressing{EffectiveAddress: 0x0113, Value: 0xAA, ProgramCounterChange: 1, Memory: &ram}

		// Validate the final State.
		if !reflect.DeepEqual(got, want) {
			t.Errorf("IndirectY() State got = %v, want = %v", got, want)
		}
	})
}

func TestRelative(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in effective address of 0x0001",
			want: Addressing{EffectiveAddress: 0x0001, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and zero memory results in effective address of 0x00F0 + 1",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{EffectiveAddress: 0x00F1, ProgramCounterChange: 1},
		},
		{
			name: "Zero input State and non-zero memory results in effective address of 0x0000 + 1 + 0x01 and Value of 0x01",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0002, Value: 0x01, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and non-zero memory results in effective address of 0x0022 + 1 + 0xFB (-5) = $001E, a program counter change of 0xFFFC (-4) and Value of 0xFB",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xFB, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0x001E, Value: 0xFB, ProgramCounterChange: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, Relative)
		})
	}
}

func TestZeroPage(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in effective address of 0x0000",
			want: Addressing{ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and zero memory results in effective address of 0x0000",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{ProgramCounterChange: 1},
		},
		{
			name: "Zero input State and non-zero memory results in effective address of 0x0001 and Value of 0x02",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0001, Value: 0x02, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and non-zero memory results in effective address of 0x00F6 and Value of 0xF2",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0x00F6, Value: 0xF2, ProgramCounterChange: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, ZeroPage)
		})
	}
}

func TestZeroPageX(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in effective address of 0x0000",
			want: Addressing{ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and zero memory results in effective address of 0x0000 + 0x71",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{EffectiveAddress: 0x0071, ProgramCounterChange: 1},
		},
		{
			name: "Zero input State and non-zero memory results in effective address of 0x0001 and Value of 0x02",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0001, Value: 0x02, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and non-zero memory results in effective address of 0x00F6 + 0x11 = $0107 = $0007 and Value of 0xF1",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0x0007, Value: 0xF1, ProgramCounterChange: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, ZeroPageX)
		})
	}
}

func TestZeroPageY(t *testing.T) {
	tests := []testAddressingConfig{
		{
			name: "Zero input State and memory results in effective address of 0x0000",
			want: Addressing{ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and zero memory results in effective address of 0x0000 + 0xAB",
			start: State{PC: 0xF0, A: 0x3F, X: 0x71, Y: 0xAB, SP: 0xE4, P: FlagCarry | FlagZero | FlagInterrupt},
			want:  Addressing{EffectiveAddress: 0x00AB, ProgramCounterChange: 1},
		},
		{
			name: "Zero input State and non-zero memory results in effective address of 0x0001 and Value of 0x02",
			ram:  []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: Addressing{EffectiveAddress: 0x0001, Value: 0x02, ProgramCounterChange: 1},
		},
		{
			name:  "Non-zero input State and non-zero memory results in effective address of 0x00F6 + 0x1FF = $01F5 = $00F5 and Value of 0xF3",
			start: State{PC: 0x22, A: 0x12, X: 0x11, Y: 0xFF, SP: StackPointerStart, P: FlagOverflow | FlagBreak},
			ram:   []uint8{0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1},
			want:  Addressing{EffectiveAddress: 0x00F5, Value: 0xF3, ProgramCounterChange: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAddressing(t, tt, ZeroPageY)
		})
	}
}
