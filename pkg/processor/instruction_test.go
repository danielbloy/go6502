package processor

import (
	"reflect"
	"testing"
)

// NewTestInstructionSet returns an instruction set that is only useful for testing and
// contains the following opcodes:
//
//   - 0x0: NOP that does not change the state of the CPU.
//     -- Requires no additional bytes or cycles.
//
//   - 0x1: Sets: A: 0x01, X: 0x0F, Y: 0xF0, P: FlagBreak
//     -- Requires no additional bytes or cycles.
//
//   - 0x2: Sets: A: 0x10, X: <low byte>, Y: <high byte>, P: FlagOverflow
//     -- Requires 2 additional bytes (<low byte>, <high byte>) and 2 additional cycles.
//     -- Uses absolute addressing which adds 2 to the PC
//
//   - 0x3: Sets X: <low byte>, Y: <high byte> and writes <low byte>, <high byte> to 0x00 and 0x01.
//     --  Requires 4 additional cycles and increases PC by 2 itself (i.e. not via AddressingFunc).
//
//   - 0x4: NOP that does not change the state of the CPU.
//     -- Requires 5 cycles + an additional cycle is page boundary crossed.
//     -- AddressingFunc mode always crosses a page boundary.
//
//   - 0x5: NOP that does not change the state of the CPU.
//     -- Requires 5 cycles regardless of whether a page boundary is crossed.
//     -- AddressingFunc mode always crosses a page boundary.
func NewTestInstructionSet() InstructionSet {

	CrossBoundary := func(_ State, _ Memory) (Addressing, error) {
		return Addressing{
			ProgramCounterChange: 2,
			Value:                0,
			PageBoundaryCrossed:  true,
		}, nil
	}

	testFunc1 := func(state State, _ Addressing) (State, error) {
		state.A = 0x01
		state.X = 0x0F
		state.Y = 0xF0
		state.P = FlagBreak
		return state, nil
	}

	testFunc2 := func(state State, addressing Addressing) (State, error) {
		state.A = 0x10
		state.X = uint8(addressing.EffectiveAddress & 0x00FF)
		state.Y = uint8(addressing.EffectiveAddress >> 8)
		state.P = FlagOverflow
		return state, nil
	}

	// Reads the next 2 byes and places them in X, Y and RAM 0x0 and 0x1
	testFunc3 := func(state State, as Addressing) (State, error) {
		state.X = as.Memory.Read(state.PC)
		state.PC++
		state.Y = as.Memory.Read(state.PC)
		state.PC++
		as.Memory.Write(0x0, state.X)
		as.Memory.Write(0x1, state.Y)
		return state, nil
	}

	instructions := Instructions{
		Instruction{
			Opcode:         0x0,
			AddressingFunc: Implied,
			Operation:      NoOperation,
		},
		Instruction{
			Opcode:         0x1,
			AddressingFunc: Implied,
			Operation:      testFunc1,
		},
		Instruction{
			Opcode:         0x2,
			AddressingFunc: Absolute,
			Operation:      testFunc2,
			Cycles:         2,
		},
		Instruction{
			Opcode:         0x3,
			AddressingFunc: Implied,
			Operation:      testFunc3,
			Cycles:         4,
		},
		Instruction{
			Opcode:              0x4,
			AddressingFunc:      CrossBoundary,
			Operation:           NoOperation,
			Cycles:              5,
			PageBoundaryPenalty: true,
		},
		Instruction{
			Opcode:              0x5,
			AddressingFunc:      CrossBoundary,
			Operation:           NoOperation,
			Cycles:              5,
			PageBoundaryPenalty: false,
		},
	}

	is, err := NewInstructionSet(instructions)
	if err != nil {
		panic(err)
	}
	return is
}

// Returns the first 3 instructions in the test instruction set.
func ThreeInstructions() (Instruction, Instruction, Instruction) {
	instructionSet := NewTestInstructionSet()
	instr1, err := instructionSet.Get(0)
	if err != nil {
		panic(err)
	}
	instr2, err := instructionSet.Get(1)
	if err != nil {
		panic(err)
	}
	instr3, err := instructionSet.Get(2)
	if err != nil {
		panic(err)
	}
	return instr1, instr2, instr3
}

// Returns the first 4 instructions in the test instruction set.
func FourInstructions() (Instruction, Instruction, Instruction, Instruction) {
	nop, instr1, instr2 := ThreeInstructions()
	instructionSet := NewTestInstructionSet()
	instr3, err := instructionSet.Get(3)
	if err != nil {
		panic(err)
	}
	return nop, instr1, instr2, instr3
}

// Returns the first 6 instructions in the test instruction set.
func SixInstructions() (Instruction, Instruction, Instruction, Instruction, Instruction, Instruction) {
	nop, instr1, instr2, instr3 := FourInstructions()
	instructionSet := NewTestInstructionSet()
	instr4, err := instructionSet.Get(4)
	if err != nil {
		panic(err)
	}
	instr5, err := instructionSet.Get(5)
	if err != nil {
		panic(err)
	}
	return nop, instr1, instr2, instr3, instr4, instr5
}

func TestInstruction_String(t *testing.T) {
	tests := []struct {
		name                string
		OpCode              Opcode
		Cycles              uint
		PageBoundaryPenalty bool
		want                string
	}{
		{
			name: "Zero value",
			want: "Opcode: $00, Cycles: 0, PBP: false",
		},
		{
			name:   "Opcode has a non-zero value",
			OpCode: 0x13,
			want:   "Opcode: $13, Cycles: 0, PBP: false",
		},
		{
			name:   "Cycles has non-zero value",
			Cycles: 3,
			want:   "Opcode: $00, Cycles: 3, PBP: false",
		},
		{
			name:                "PageBoundary is true",
			PageBoundaryPenalty: true,
			want:                "Opcode: $00, Cycles: 0, PBP: true",
		},
		{
			name:   "Mixed values 1",
			OpCode: 0xFA,
			Cycles: 7,
			want:   "Opcode: $FA, Cycles: 7, PBP: false",
		},
		{
			name:                "Mixed values 2",
			OpCode:              0x9F,
			Cycles:              7,
			PageBoundaryPenalty: true,
			want:                "Opcode: $9F, Cycles: 7, PBP: true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Instruction{
				Opcode:              tt.OpCode,
				Cycles:              tt.Cycles,
				PageBoundaryPenalty: tt.PageBoundaryPenalty,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstruction_Execute(t *testing.T) {

	nop, instr1, instr2, instr3, instr4, instr5 := SixInstructions()

	noAddressing := nop
	noAddressing.AddressingFunc = nil

	noOperation := nop
	noOperation.Operation = nil

	tests := []struct {
		name        string
		instruction Instruction
		state       State
		startRam    []uint8
		wantState   State
		wantRam     []uint8
		wantCycles  uint
		wantErr     bool
	}{
		{
			name:        "Nil Memory results in an error",
			instruction: nop,
			wantErr:     true,
		},
		{
			name:        "Nil addressing results in an error",
			instruction: noAddressing,
			startRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:     true,
		},
		{
			name:        "Nil operation results in an error",
			instruction: noOperation,
			startRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:     true,
		},
		{
			name:        "All arguments passed correctly, nop",
			instruction: nop,
			state:       State{PC: 0x3},
			startRam:    []uint8{0xA0, 0x12, 0, 0, 0, 0, 0, 0},
			wantState:   State{PC: 0x3},
			wantRam:     []uint8{0xA0, 0x12, 0, 0, 0, 0, 0, 0},
		},
		{
			name:        "All arguments passed correctly, instruction 1",
			instruction: instr1,
			state:       State{PC: 2},
			startRam:    []uint8{0, 0, 0xA0, 0x12, 0, 0, 0, 0},
			wantState:   State{PC: 2, A: 0x01, X: 0x0F, Y: 0xF0, P: FlagBreak},
			wantRam:     []uint8{0, 0, 0xA0, 0x12, 0, 0, 0, 0},
		},
		{
			name:        "All arguments passed correctly, instruction 2",
			instruction: instr2,
			state:       State{PC: 2, A: 0xFF, X: 0xFF, Y: 0xFF},
			startRam:    []uint8{0, 0, 0xA0, 0x12, 0, 0, 0, 0},
			wantState:   State{PC: 4, A: 0x10, X: 0xA0, Y: 0x12, P: FlagOverflow},
			wantRam:     []uint8{0, 0, 0xA0, 0x12, 0, 0, 0, 0},
			wantCycles:  2,
		},
		{
			name:        "All arguments passed correctly, instruction 3 which writes to Memory 0",
			instruction: instr3,
			state:       State{PC: 4, A: 0xFF, X: 0xFF, Y: 0xFF, P: FlagCarry | FlagZero},
			startRam:    []uint8{0, 0, 0, 0, 0x0A, 0x0B, 0, 0},
			wantState:   State{PC: 6, A: 0xFF, X: 0x0A, Y: 0x0B, P: FlagCarry | FlagZero},
			wantRam:     []uint8{0x0A, 0x0B, 0, 0, 0x0A, 0x0B, 0, 0},
			wantCycles:  4,
		},
		{
			name:        "Page boundary crossed with susceptible operation",
			instruction: instr4,
			state:       State{PC: 0x5},
			startRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState:   State{PC: 0x7},
			wantRam:     []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantCycles:  6,
		},
		{
			name:        "Page boundary crossed with non-susceptible operation",
			instruction: instr5,
			state:       State{PC: 0x5},
			startRam:    []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantState:   State{PC: 0x7},
			wantRam:     []uint8{0, 0, 0, 0, 0, 0, 0, 0},
			wantCycles:  5,
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

			gotState, gotCycles, err := tt.instruction.Execute(tt.state, memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(gotState, tt.wantState) {
				t.Errorf("Execute() State got = %v, want %v", gotState, tt.wantState)
			}
			// Validate expected final RAM state
			if !reflect.DeepEqual(ram.ram, tt.wantRam) {
				t.Errorf("Execute() RAM state got = %v, want = %v", ram.ram, tt.wantRam)
			}
			if gotCycles != tt.wantCycles {
				t.Errorf("Execute() cycles got = %v, want %v", gotCycles, tt.wantCycles)
			}
		})
	}
}

func TestInstructionSet_validate(t *testing.T) {
	tests := []struct {
		name    string
		is      InstructionSet
		wantErr bool
	}{
		{
			name:    "Invalid InstructionSet should error",
			wantErr: true,
		},
		{
			name: "Valid InstructionSet wont error",
			is: InstructionSet{
				instructions: map[Opcode]Instruction{
					0: {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.is.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstructionSet_Get(t *testing.T) {
	instructionSet := InstructionSet{
		instructions: map[Opcode]Instruction{
			0x00: {
				Opcode: 0x00,
				Cycles: 1,
			},
			0xF0: {
				Opcode:              0xF0,
				Cycles:              2,
				PageBoundaryPenalty: true,
			},
			0xFF: {
				Opcode:              0xFF,
				Cycles:              3,
				PageBoundaryPenalty: true,
			},
		},
	}
	tests := []struct {
		name    string
		is      InstructionSet
		opcode  Opcode
		want    Instruction
		wantErr bool
	}{
		{
			name:    "Invalid InstructionSet should error",
			wantErr: true,
		},
		{
			name:    "Invalid instruction code should error 1",
			is:      instructionSet,
			opcode:  0x1F,
			wantErr: true,
		},
		{
			name:    "Invalid instruction code should error 2",
			is:      instructionSet,
			opcode:  0xFE,
			wantErr: true,
		},
		{
			name:   "Valid instruction code should return correct instruction 1",
			is:     instructionSet,
			opcode: 0x00,
			want:   instructionSet.instructions[0x00],
		},
		{
			name:   "Valid instruction code should return correct instruction 2",
			is:     instructionSet,
			opcode: 0xFF,
			want:   instructionSet.instructions[0xFF],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.is.Get(tt.opcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstructionSet_Fill(t *testing.T) {
	instructionSet := InstructionSet{
		instructions: map[Opcode]Instruction{
			0x00: {
				Opcode: 0x00,
				Cycles: 1,
			},
			0xF0: {
				Opcode:              0xF0,
				Cycles:              2,
				PageBoundaryPenalty: true,
			},
			0xFF: {
				Opcode:              0xFF,
				Cycles:              3,
				PageBoundaryPenalty: true,
			},
		},
	}
	tests := []struct {
		name    string
		is      InstructionSet
		fill    Instruction
		wantErr bool
	}{
		{
			name:    "Invalid InstructionSet should error",
			wantErr: true,
		},
		{
			name: "Fill with zero instruction",
			is:   instructionSet,
		},
		{
			name: "Fill with new instructions",
			is:   instructionSet,
			fill: Instruction{
				Opcode:              0xAA,
				Cycles:              20,
				PageBoundaryPenalty: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.is.Fill(tt.fill)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fill() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			if len(tt.is.instructions) != 256 {
				t.Errorf("Fill() did not result in a full instruction set. have = %v", len(tt.is.instructions))
			}

			// Loop through every instruction and ensure that they are as expected.
			for opcode := range 255 {
				instruction := tt.is.instructions[Opcode(opcode)]
				if instruction.Opcode != Opcode(opcode) {
					t.Errorf("Fill() unexpected opcode; got = %02X, want %02X", instruction.Opcode, opcode)
				}

				// No need to test the originally provided opcodes.
				if opcode == 0x00 || opcode == 0xF0 || opcode == 0xFF {
					return
				}

				if instruction.Cycles != tt.fill.Cycles {
					t.Errorf("Fill() unexpected Cycles; got = %v, want %v", instruction.Cycles, tt.fill.Cycles)
				}
				if instruction.PageBoundaryPenalty != tt.fill.PageBoundaryPenalty {
					t.Errorf("Fill() unexpected PBP ; got = %v, want %v", instruction.PageBoundaryPenalty, tt.fill.PageBoundaryPenalty)
				}
			}
		})
	}
}

func TestInstructionSet_Opcodes(t *testing.T) {
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

			got := tt.instructionSet.Opcodes()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unexpected opcodes got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestNewInstructionSet(t *testing.T) {
	tests := []struct {
		name    string
		is      Instructions
		want    InstructionSet
		wantErr bool
	}{
		{
			name:    "Nil instructions should result in error",
			is:      nil,
			wantErr: true,
		},
		{
			name:    "Empty instructions should result in error",
			is:      Instructions{},
			wantErr: true,
		},
		{
			name: "Single instruction is fine 1",
			is:   Instructions{Instruction{}},
			want: InstructionSet{
				instructions: map[Opcode]Instruction{
					0x00: {},
				},
			},
		},
		{
			name: "Single instruction is fine 2",
			is:   Instructions{Instruction{Opcode: 0xF0, Cycles: 3}},
			want: InstructionSet{
				instructions: map[Opcode]Instruction{
					0xF0: {Opcode: 0xF0, Cycles: 3},
				},
			},
		},
		{
			name: "Multiple instruction is fine",
			is: Instructions{
				Instruction{Opcode: 0x00, Cycles: 3},
				Instruction{Opcode: 0xF0, Cycles: 7, PageBoundaryPenalty: true},
				Instruction{Opcode: 0xFF, Cycles: 1},
			},
			want: InstructionSet{
				instructions: map[Opcode]Instruction{
					0x00: {Opcode: 0x00, Cycles: 3},
					0xF0: {Opcode: 0xF0, Cycles: 7, PageBoundaryPenalty: true},
					0xFF: {Opcode: 0xFF, Cycles: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInstructionSet(tt.is)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInstructionSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInstructionSet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
