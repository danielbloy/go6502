package processor

import (
	"reflect"
	"testing"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		name  string
		state State
		want  string
	}{
		{
			name: "Zero State",
			want: "PC: 0x0000, SP: 0x00, A: 0x00, X: 0x00, Y: 0x00, P: n v - b d i z c",
		},
		{
			name:  "Set just Program counter",
			state: State{PC: 0xFEDC},
			want:  "PC: 0xFEDC, SP: 0x00, A: 0x00, X: 0x00, Y: 0x00, P: n v - b d i z c",
		},
		{
			name:  "Set just processor status",
			state: State{P: FlagCarry | FlagDecimal | FlagInterrupt},
			want:  "PC: 0x0000, SP: 0x00, A: 0x00, X: 0x00, Y: 0x00, P: n v - b D I z C",
		},
		{
			name:  "Mixture 1",
			state: State{PC: 0x1234, SP: 0xAB, A: 0xCD, X: 0xEF, Y: 0x99, P: FlagOverflow},
			want:  "PC: 0x1234, SP: 0xAB, A: 0xCD, X: 0xEF, Y: 0x99, P: n V - b d i z c",
		},
		{
			name:  "Mixture 2",
			state: State{PC: 0x00DD, SP: 0x02, A: 0x04, X: 0x06, Y: 0x0F, P: FlagNegative | FlagZero},
			want:  "PC: 0x00DD, SP: 0x02, A: 0x04, X: 0x06, Y: 0x0F, P: N v - b d i Z c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCpu(t *testing.T) {
	sampleInstructionSet, err := NewInstructionSet(Instructions{{}})
	if err != nil {
		panic(err)
	}
	memory, err := NewRepeatingRam(SixteenBytes)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		is      InstructionSet
		memory  Memory
		want    Cpu
		wantErr bool
	}{
		{
			name:    "Nil instruction set and memory should error",
			wantErr: true,
		},
		{
			name:    "Nil instruction set should error",
			wantErr: true,
			memory:  &memory,
		},
		{
			name:    "Nil memory should error",
			wantErr: true,
			is:      sampleInstructionSet,
		},
		{
			name:    "Empty instruction set should error",
			wantErr: true,
			is:      InstructionSet{instructions: map[Opcode]Instruction{}},
		},
		{
			name:   "Valid instruction set and memory should be fine",
			is:     sampleInstructionSet,
			memory: &memory,
			want:   Cpu{memory: &memory, instructionSet: sampleInstructionSet},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCpu(tt.is, tt.memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCpu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.State, tt.want.State) {
				t.Errorf("CompareCpus() got State = %v, want State %v", got.State, tt.want.State)
			}
			if !reflect.DeepEqual(got.memory, tt.want.memory) {
				t.Errorf("CompareCpus() got memory = %v, want memory %v", got.memory, tt.want.State)
			}

			// We can only really compare opcode as reflect.DeepEqual does not work with function pointers.
			wantOpcodes, _ := tt.want.Opcodes()
			gotOpcodes, _ := got.Opcodes()

			if !reflect.DeepEqual(wantOpcodes, gotOpcodes) {
				t.Errorf("CompareCpus() got opcodes = %v, want opcodes %v", wantOpcodes, gotOpcodes)
			}
		})
	}
}

// The following struct and harness function are used to test the individual Cpu
// methods by sorting out all the boilerplate.
type testCpuMethodConfig struct {
	name           string
	startState     State
	startRam       []uint8
	instructionSet InstructionSet
	wantState      State
	wantRam        []uint8
	wantErr        bool
}

func testCpuMethod(t *testing.T, tt testCpuMethodConfig, callback func(cpu *Cpu) error) {
	// Create CPU with populated memory and a default test instruction set if one is not specified.
	memory := NewPopulatedRam(RepeatingRamSize(len(tt.startRam)), tt.startRam)
	is := tt.instructionSet
	if is.validate() != nil {
		is = NewTestInstructionSet()
	}

	cpu, err := NewCpu(is, &memory)
	if err != nil {
		panic(err)
	}

	cpu.State = tt.startState

	// Call the required CPU method using the provided callback.
	if err = callback(&cpu); (err != nil) != tt.wantErr {
		t.Errorf("testCpuMethod() error = %v, wantErr %v", err, tt.wantErr)
	}

	// Validate expected final CPU State
	if !reflect.DeepEqual(cpu.State, tt.wantState) {
		t.Errorf("Unexpected CPU State got = %v, want = %v", cpu.State, tt.wantState)
	}

	// Validate expected final RAM State
	if !reflect.DeepEqual(memory.ram, tt.wantRam) {
		t.Errorf("Unexpected RAM State got = %v, want = %v", memory.ram, tt.wantRam)
	}
}

func TestCpu_Reset(t *testing.T) {
	callback := func(cpu *Cpu) error {
		return cpu.Reset()
	}
	tests := []testCpuMethodConfig{
		{
			name:      "Reset with no reset vector in place",
			startRam:  []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState: State{SP: StackPointerStart},
			wantRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:      "Reset with a reset vector 0x1234 in place",
			startRam:  []uint8{0, 0, 0, 0, 0x34, 0x12, 0, 0},
			wantState: State{PC: MakeAddress(0x34, 0x12), SP: StackPointerStart},
			wantRam:   []uint8{0, 0, 0, 0, 0x34, 0x12, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCpuMethod(t, tt, callback)
		})
	}

	// Check support for nil
	if err := (*Cpu)(nil).Reset(); err == nil {
		t.Errorf("Reset() did not raise an error when called on nil")
	}
}

func TestCpu_Step(t *testing.T) {
	tests := []struct {
		wantCycles uint
		testCpuMethodConfig
	}{
		{
			wantCycles: 1,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Single step after a reset with single default opcode.",
				startRam:  []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 1, SP: StackPointerStart},
			},
		},
		{
			wantCycles: 1,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Single step after a reset with an unknown opcode.",
				startRam:  []uint8{0xFF, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0xFF, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 1, SP: StackPointerStart},
				wantErr:   true,
			},
		},
		{
			wantCycles: 1,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Single step after a reset with a known opcode; no State mutation.",
				startRam:  []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 1, SP: StackPointerStart},
			},
		},
		{
			wantCycles: 1,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Single step after a reset with a known opcode; with State mutation.",
				startRam:  []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 1, SP: StackPointerStart, A: 0x01, X: 0x0F, Y: 0xF0, P: FlagBreak},
			},
		},
		{
			wantCycles: 3,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Single step after a reset with a known opcode; with State mutation and multiple cycles.",
				startRam:  []uint8{0x02, 0xA0, 0x12, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0x02, 0xA0, 0x12, 0, 0, 0, 0, 0},
				wantState: State{PC: 3, SP: StackPointerStart, A: 0x10, X: 0xA0, Y: 0x12, P: FlagOverflow},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callback := func(cpu *Cpu) error {
				if err := cpu.Reset(); err != nil {
					return err
				}
				cycles, err := cpu.Step()
				if cycles != tt.wantCycles {
					t.Errorf("Step() did not execute the expected number of cycles, got = %v, want = %v", cycles, tt.wantCycles)
				}
				return err
			}

			testCpuMethod(t, tt.testCpuMethodConfig, callback)
		})
	}

	// Check support for nil
	if cycles, err := (*Cpu)(nil).Step(); err == nil || cycles != 0 {
		if err != nil {
			t.Errorf("Step() did not raise an error when called on nil")
		}
		if cycles != 0 {
			t.Errorf("Step() did not return zero cycles when called on nil")
		}
	}
}

func TestCpu_Execute(t *testing.T) {

	tests := []struct {
		cycles     uint
		wantCycles uint
		testCpuMethodConfig
	}{
		{
			cycles:     1,
			wantCycles: 1,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Execute a single cycle default instruction",
				startRam:  []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 1, SP: StackPointerStart},
			},
		},
		{
			cycles:     2,
			wantCycles: 2,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Execute a two cycles with default instruction",
				startRam:  []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 2, SP: StackPointerStart},
			},
		},
		{
			cycles:     3,
			wantCycles: 3,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Execute three cycles with default instruction.",
				startRam:  []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 3, SP: StackPointerStart},
			},
		},
		{
			cycles:     2,
			wantCycles: 2,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Execute two cycles with known opcodes; with State mutation.",
				startRam:  []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0x01, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 2, SP: StackPointerStart, A: 0x01, X: 0x0F, Y: 0xF0, P: FlagBreak},
			},
		},
		{
			cycles:     2,
			wantCycles: 4,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "Single step after a reset with an known opcode; with State mutation and multiple cycles.",
				startRam:  []uint8{0x01, 0x02, 0xA0, 0x0A, 0, 0, 0, 0},
				wantRam:   []uint8{0x01, 0x02, 0xA0, 0x0A, 0, 0, 0, 0},
				wantState: State{PC: 4, SP: StackPointerStart, A: 0x10, X: 0xA0, Y: 0x0A, P: FlagOverflow},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Don't allow for infinite running.
			if tt.cycles == 0 {
				t.Errorf("Execute() cycles of zero was specified which will run for ever!")
				return
			}

			callback := func(cpu *Cpu) error {
				if err := cpu.Reset(); err != nil {
					return err
				}
				cycles, err := cpu.Execute(tt.cycles)
				if cycles != tt.wantCycles {
					t.Errorf("Execute() did not execute the expected number of cycles, got = %v, want = %v", cycles, tt.wantCycles)
				}
				return err
			}

			testCpuMethod(t, tt.testCpuMethodConfig, callback)
		})
	}

	// Check support for nil
	if cycles, err := (*Cpu)(nil).Execute(0); err == nil || cycles != 0 {
		if err != nil {
			t.Errorf("Execute() did not raise an error when called on nil")
		}
		if cycles != 0 {
			t.Errorf("Execute() did not return zero cycles when called on nil")
		}
	}
}

func TestCpu_Nmi(t *testing.T) {

	tests := []testCpuMethodConfig{
		{
			name:       "Nmi with mostly zeros",
			startState: State{SP: 0xF3},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF0, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with PC set",
			startState: State{PC: 0xF0A0, SP: 0xF3},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF0, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0xA0, 0xF0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with PC set and flags",
			startState: State{PC: 0xF0A0, SP: 0xF3, P: FlagCarry | FlagZero},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF0, P: FlagCarry | FlagZero | FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant | FlagCarry | FlagZero, 0xA0, 0xF0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with different SP set",
			startState: State{PC: 0xF0A0, SP: 0xF5},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
			wantState:  State{PC: 0xA002, SP: 0xF2, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0, 0, 0, 0, 0x02, 0xA0, 0, 0, 0, 0},
		},
		{
			name:       "Nmi with different vector",
			startState: State{PC: 0xF0A0, SP: 0xF5},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xEE, 0xFF, 0, 0, 0, 0},
			wantState:  State{PC: 0xFFEE, SP: 0xF2, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0, 0, 0, 0, 0xEE, 0xFF, 0, 0, 0, 0},
		},
		{
			name:       "Nmi when interrupt flag is set",
			startState: State{PC: 0xF0A0, SP: 0xF5, P: FlagInterrupt},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xEE, 0xFF, 0, 0, 0, 0},
			wantState:  State{PC: 0xFFEE, SP: 0xF2, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant | FlagInterrupt, 0xA0, 0xF0, 0, 0, 0, 0, 0xEE, 0xFF, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {

		callback := func(cpu *Cpu) error {
			return cpu.Nmi()
		}

		testCpuMethod(t, tt, callback)
	}

	// Check support for nil
	if err := (*Cpu)(nil).Nmi(); err == nil {
		t.Errorf("Nmi() did not raise an error when called on nil")
	}
}

func TestCpu_Interrupt(t *testing.T) {

	tests := []testCpuMethodConfig{
		{
			name:       "Interrupt with mostly zeros",
			startState: State{SP: 0xFB},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0x00, 0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with PC set",
			startState: State{PC: 0xF0A0, SP: 0xFB},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with PC set and flags",
			startState: State{PC: 0xF0A0, SP: 0xFB, P: FlagCarry | FlagZero},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xF8, P: FlagCarry | FlagZero | FlagInterrupt},
			wantRam:    []uint8{0, FlagConstant | FlagCarry | FlagZero, 0xA0, 0xF0, 0, 0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with different SP set",
			startState: State{PC: 0xF0A0, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0x02, 0xA0},
			wantState:  State{PC: 0xA002, SP: 0xFA, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0x02, 0xA0},
		},
		{
			name:       "Interrupt with different vector",
			startState: State{PC: 0xF0A0, SP: 0xFD},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0xEE, 0xFF},
			wantState:  State{PC: 0xFFEE, SP: 0xFA, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, FlagConstant, 0xA0, 0xF0, 0xEE, 0xFF},
		},
		{
			name:       "Interrupt is ignored when interrupt flag is set",
			startState: State{PC: 0xF0A0, SP: 0xFD, P: FlagInterrupt},
			startRam:   []uint8{0, 0, 0, 0, 0, 0, 0xEE, 0xFF},
			wantState:  State{PC: 0xF0A0, SP: 0xFD, P: FlagInterrupt},
			wantRam:    []uint8{0, 0, 0, 0, 0, 0, 0xEE, 0xFF},
		},
	}
	for _, tt := range tests {
		callback := func(cpu *Cpu) error {
			return cpu.Interrupt()
		}

		testCpuMethod(t, tt, callback)
	}

	// Check support for nil
	if err := (*Cpu)(nil).Interrupt(); err == nil {
		t.Errorf("Interrupt() did not raise an error when called on nil")
	}
}

func TestCpu_Opcodes(t *testing.T) {
	instr1, instr2, instr3, instr4, instr5, instr6 := SixInstructions()

	is1, err := NewInstructionSet(Instructions{instr3})
	if err != nil {
		panic(err)
	}

	is3, err := NewInstructionSet(Instructions{instr4, instr3, instr6})
	if err != nil {
		panic(err)
	}

	is6, err := NewInstructionSet(Instructions{instr6, instr5, instr4, instr3, instr2, instr1})
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name           string
		instructionSet InstructionSet
		want           []Opcode
	}{
		{
			name:           "Test with 1 instructions",
			instructionSet: is1,
			want:           []Opcode{instr3.Opcode},
		},
		{
			name:           "Test with 3 instructions",
			instructionSet: is3,
			want:           []Opcode{instr3.Opcode, instr4.Opcode, instr6.Opcode},
		},
		{
			name:           "Test with 6 instructions",
			instructionSet: is6,
			want:           []Opcode{instr1.Opcode, instr2.Opcode, instr3.Opcode, instr4.Opcode, instr5.Opcode, instr6.Opcode},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create CPU with populated memory and a default test instruction set if one is not specified.
			ram := []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80}
			memory := NewPopulatedRam(RepeatingRamSize(len(ram)), ram)

			cpu, err := NewCpu(tt.instructionSet, &memory)
			if err != nil {
				panic(err)
			}

			got, err := cpu.Opcodes()
			if err != nil {
				panic(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unexpected opcodes got = %v, want = %v", got, tt.want)
			}
		})
	}

	// Check support for nil
	if _, err := (*Cpu)(nil).Opcodes(); err == nil {
		t.Errorf("Opcodes() did not raise an error when called on nil")
	}
}

func TestCpu_Memory(t *testing.T) {

	// Create CPU with populated memory and a default test instruction set if one is not specified.
	ram := []uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80}
	memory := NewPopulatedRam(RepeatingRamSize(len(ram)), ram)
	is := NewTestInstructionSet()

	cpu, err := NewCpu(is, &memory)
	if err != nil {
		panic(err)
	}

	gotMemory, err := cpu.Memory()
	if gotMemory != &memory {
		t.Errorf("Unexpected RAM State got = %v, want = %v", memory.ram, ram)
	}

	// Check support for nil
	if _, err := (*Cpu)(nil).Memory(); err == nil {
		t.Errorf("Memory() did not raise an error when called on nil")
	}
}
