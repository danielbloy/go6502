package processor

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewInstruction(t *testing.T) {
	adc := MnemonicOperation{
		Name:                 "Add with carry",
		AssemblyLanguageForm: "ADC",
		AffectedFlags:        Flags{Negative: true, Overflow: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            AddWithCarry,
	}
	brk := MnemonicOperation{
		Name:                 "Break",
		AssemblyLanguageForm: "BRK",
		AffectedFlags:        Flags{Break: true},
		Bytes:                1,
		Cycles:               7,
		Operation:            Break,
	}
	immediate := MnemonicAddressingMode{
		Name:                 "Immediate",
		AssemblyLanguageForm: "#$%02X",
		Bytes:                1,
		Cycles:               1,
		AddressingFunc:       Immediate,
	}
	implied := MnemonicAddressingMode{
		Name:                 "Implied",
		AssemblyLanguageForm: "",
		Bytes:                0,
		Cycles:               0,
		AddressingFunc:       Implied,
	}

	tests := []struct {
		name        string
		opcode      Opcode
		mnemonic    MnemonicOperation
		addressing  MnemonicAddressingMode
		cycleAdjust int
		want        Instruction
	}{
		{
			name: "Zero opcode in zero mnemonic cannot be found",
		},
		{
			name:       "Valid opcode and addressing mode returns instruction (ADC #)",
			opcode:     0x69,
			mnemonic:   adc,
			addressing: immediate,
			want: Instruction{
				Opcode:         0x69,
				AddressingFunc: immediate.AddressingFunc,
				Operation:      adc.Operation,
				Cycles:         1,
			},
		},
		{
			name:       "Valid opcode from valid mnemonic returns instruction (BRK)",
			opcode:     0x00,
			mnemonic:   brk,
			addressing: implied,
			want: Instruction{
				Opcode:         0x00,
				AddressingFunc: Implied,
				Operation:      adc.Operation,
				Cycles:         6,
			},
		},
		{
			name:        "Valid opcode from valid mnemonic with additional cycles returns instruction (BRK)",
			opcode:      0x00,
			mnemonic:    brk,
			addressing:  implied,
			cycleAdjust: 3,
			want: Instruction{
				Opcode:         0x00,
				AddressingFunc: Immediate,
				Operation:      adc.Operation,
				Cycles:         9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInstruction(
				Mnemonic{tt.opcode, tt.mnemonic, tt.addressing, tt.cycleAdjust})

			// DeepReflect won't work here as it cannot process function pointers. We
			// Just go for a close match.
			match := got.Opcode == tt.want.Opcode
			match = match && (got.AddressingFunc != nil) == (tt.want.AddressingFunc != nil)
			match = match && (got.Operation != nil) == (tt.want.Operation != nil)
			match = match && got.Cycles == tt.want.Cycles
			match = match && got.PageBoundaryPenalty == tt.want.PageBoundaryPenalty
			if !match {
				t.Errorf("NewInstruction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// MnemonicsMatch returns if the two mnemonics match or not.
func mnemonicsMatch(got, want Mnemonic) bool {
	// DeepReflect won't work here as it cannot process function pointers. We
	// Just go for a close match.
	match := got.Opcode == want.Opcode

	match = match && got.Operation.Name == want.Operation.Name
	match = match && got.Operation.Description == want.Operation.Description
	match = match && got.Operation.AssemblyLanguageForm == want.Operation.AssemblyLanguageForm
	match = match && got.Operation.AffectedFlags == want.Operation.AffectedFlags
	match = match && got.Operation.Bytes == want.Operation.Bytes
	match = match && got.Operation.Cycles == want.Operation.Cycles
	match = match && got.Operation.PageBoundaryPenalty == want.Operation.PageBoundaryPenalty

	match = match && got.Addressing.Name == want.Addressing.Name
	match = match && got.Addressing.AssemblyLanguageForm == want.Addressing.AssemblyLanguageForm
	match = match && got.Addressing.Bytes == want.Addressing.Bytes
	match = match && got.Addressing.Cycles == want.Addressing.Cycles

	match = match && got.CycleAdjust == want.CycleAdjust

	return match
}

func TestMnemonicFromOpCode(t *testing.T) {

	const (
		brk          = 0x00
		adcZeroPage  = 0x65
		adcAbsoluteX = 0x7D
		aslZeroPage  = 0x06
	)

	tests := []struct {
		name          string
		opcode        Opcode
		wantOperation MnemonicOperation
		wantAddress   MnemonicAddressingMode
		cycleAdjust   int
		wantErr       bool
	}{
		{name: "Unknown opcode", opcode: 0x80, wantErr: true},
		{name: "ADC zpg", opcode: adcZeroPage, wantOperation: Adc, wantAddress: Zpg},
		{name: "ADC abs, x", opcode: adcAbsoluteX, wantOperation: Adc, wantAddress: AbsX},
		{name: "BRK", opcode: brk, wantOperation: Brk, wantAddress: Imp},
		{name: "ASL zpg", opcode: aslZeroPage, wantOperation: Asl, wantAddress: Zpg, cycleAdjust: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MnemonicFromOpCode(tt.opcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("MnemonicFromOpCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			match := mnemonicsMatch(got,
				Mnemonic{tt.opcode, tt.wantOperation, tt.wantAddress, tt.cycleAdjust})
			if !match {
				t.Errorf("MnemonicFromOpCode() got = %v, want %v", got, tt.wantOperation)
			}
		})
	}
}

func TestMnemonicsFromInstructionName(t *testing.T) {
	tests := []struct {
		name        string
		instruction string
		want        []Mnemonic
		wantErr     bool
	}{

		{
			name:        "Empty string errors",
			instruction: "",
			wantErr:     true,
		},
		{
			name:        "Valid single opcode same case",
			instruction: "BRK",
			want:        []Mnemonic{mnemonics[0x00]},
		},
		{
			name:        "Valid single opcode different case",
			instruction: "nOp",
			want:        []Mnemonic{mnemonics[0xEA]},
		},
		{
			name:        "Valid multiple opcodes different case",
			instruction: "ADc",
			want: []Mnemonic{
				mnemonics[0x61], mnemonics[0x65],
				mnemonics[0x69], mnemonics[0x6D],
				mnemonics[0x71], mnemonics[0x75],
				mnemonics[0x79], mnemonics[0x7D],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MnemonicsFromInstructionName(tt.instruction)
			if (err != nil) != tt.wantErr {
				fmt.Println(got)
				t.Errorf("MnemonicsFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("MnemonicsFromName() got = %v, want %v", got, tt.want)
			}
			for i := range got {
				match := mnemonicsMatch(got[i], tt.want[i])
				if !match {
					t.Errorf("MnemonicsFromName() index = %v got = %v, want %v", i, got[i], tt.want[i])
				}

			}
		})
	}
}

func TestMnemonicsDetails(t *testing.T) {
	type expected struct {
		addressing string
		assembler  string
		opcode     Opcode
		bytes      uint
		cycles     uint
		penalty    bool
		flags      Flags
	}
	tests := []struct {
		instruction string
		wantFlags   Flags
		want        []expected
	}{
		{
			instruction: "ADC",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true, Overflow: true},
			want: []expected{
				{opcode: 0x61, bytes: 2, cycles: 6, assembler: "ADC ($FF,X)"},
				{opcode: 0x65, bytes: 2, cycles: 3, assembler: "ADC $FF"},
				{opcode: 0x69, bytes: 2, cycles: 2, assembler: "ADC #$FF"},
				{opcode: 0x6D, bytes: 3, cycles: 4, assembler: "ADC $FFFF"},
				{opcode: 0x71, bytes: 2, cycles: 5, assembler: "ADC ($FF),Y", penalty: true},
				{opcode: 0x75, bytes: 2, cycles: 4, assembler: "ADC $FF,X"},
				{opcode: 0x79, bytes: 3, cycles: 4, assembler: "ADC $FFFF,Y", penalty: true},
				{opcode: 0x7D, bytes: 3, cycles: 4, assembler: "ADC $FFFF,X", penalty: true},
			},
		},
		{
			instruction: "AND",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0x21, bytes: 2, cycles: 6, assembler: "AND ($FF,X)"},
				{opcode: 0x25, bytes: 2, cycles: 3, assembler: "AND $FF"},
				{opcode: 0x29, bytes: 2, cycles: 2, assembler: "AND #$FF"},
				{opcode: 0x2D, bytes: 3, cycles: 4, assembler: "AND $FFFF"},
				{opcode: 0x31, bytes: 2, cycles: 5, assembler: "AND ($FF),Y", penalty: true},
				{opcode: 0x35, bytes: 2, cycles: 4, assembler: "AND $FF,X"},
				{opcode: 0x39, bytes: 3, cycles: 4, assembler: "AND $FFFF,Y", penalty: true},
				{opcode: 0x3D, bytes: 3, cycles: 4, assembler: "AND $FFFF,X", penalty: true},
			},
		},
		{
			instruction: "ASL",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true},
			want: []expected{
				{opcode: 0x06, bytes: 2, cycles: 5, assembler: "ASL $FF"},
				{opcode: 0x0A, bytes: 1, cycles: 2, assembler: "ASL A"},
				{opcode: 0x0E, bytes: 3, cycles: 6, assembler: "ASL $FFFF"},
				{opcode: 0x16, bytes: 2, cycles: 6, assembler: "ASL $FF,X"},
				{opcode: 0x1e, bytes: 3, cycles: 7, assembler: "ASL $FFFF,X"},
			},
		},
		{
			instruction: "BCC",
			want: []expected{
				{opcode: 0x90, bytes: 2, cycles: 2, assembler: "BCC $FF", penalty: true},
			},
		},
		{
			instruction: "BCS",
			want: []expected{
				{opcode: 0xB0, bytes: 2, cycles: 2, assembler: "BCS $FF", penalty: true},
			},
		},
		{
			instruction: "BEQ",
			want: []expected{
				{opcode: 0xF0, bytes: 2, cycles: 2, assembler: "BEQ $FF", penalty: true},
			},
		},
		{
			instruction: "BIT",
			wantFlags:   Flags{Negative: true, Zero: true, Overflow: true},
			want: []expected{
				{opcode: 0x24, bytes: 2, cycles: 3, assembler: "BIT $FF"},
				{opcode: 0x2C, bytes: 3, cycles: 4, assembler: "BIT $FFFF"}},
		},
		{
			instruction: "BMI",
			want: []expected{
				{opcode: 0x30, bytes: 2, cycles: 2, assembler: "BMI $FF", penalty: true},
			},
		},
		{
			instruction: "BNE",
			want: []expected{
				{opcode: 0xD0, bytes: 2, cycles: 2, assembler: "BNE $FF", penalty: true},
			},
		},
		{
			instruction: "BPL",
			want: []expected{
				{opcode: 0x10, bytes: 2, cycles: 2, assembler: "BPL $FF", penalty: true},
			},
		},
		{
			instruction: "BRK",
			wantFlags:   Flags{Interrupt: true, Break: true},
			want: []expected{
				{opcode: 0x00, bytes: 1, cycles: 7, assembler: "BRK"},
			},
		},
		{
			instruction: "BVC",
			want: []expected{
				{opcode: 0x50, bytes: 2, cycles: 2, assembler: "BVC $FF", penalty: true},
			},
		},
		{
			instruction: "BVS",
			want: []expected{
				{opcode: 0x70, bytes: 2, cycles: 2, assembler: "BVS $FF", penalty: true},
			},
		},
		{
			instruction: "CLC",
			wantFlags:   Flags{Carry: true},
			want: []expected{
				{opcode: 0x18, bytes: 1, cycles: 2, assembler: "CLC"},
			},
		},
		{
			instruction: "CLD",
			wantFlags:   Flags{Decimal: true},
			want: []expected{
				{opcode: 0xD8, bytes: 1, cycles: 2, assembler: "CLD"},
			},
		},
		{
			instruction: "CLI",
			wantFlags:   Flags{Interrupt: true},
			want: []expected{
				{opcode: 0x58, bytes: 1, cycles: 2, assembler: "CLI"},
			},
		},
		{
			instruction: "CLV",
			wantFlags:   Flags{Overflow: true},
			want: []expected{
				{opcode: 0xB8, bytes: 1, cycles: 2, assembler: "CLV"},
			},
		},
		{
			instruction: "CMP",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true},
			want: []expected{
				{opcode: 0xC1, bytes: 2, cycles: 6, assembler: "CMP ($FF,X)"},
				{opcode: 0xC5, bytes: 2, cycles: 3, assembler: "CMP $FF"},
				{opcode: 0xC9, bytes: 2, cycles: 2, assembler: "CMP #$FF"},
				{opcode: 0xCD, bytes: 3, cycles: 4, assembler: "CMP $FFFF"},
				{opcode: 0xD1, bytes: 2, cycles: 5, assembler: "CMP ($FF),Y", penalty: true},
				{opcode: 0xD5, bytes: 2, cycles: 4, assembler: "CMP $FF,X"},
				{opcode: 0xD9, bytes: 3, cycles: 4, assembler: "CMP $FFFF,Y", penalty: true},
				{opcode: 0xDD, bytes: 3, cycles: 4, assembler: "CMP $FFFF,X", penalty: true},
			},
		},
		{
			instruction: "CPX",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true},
			want: []expected{
				{opcode: 0xE0, bytes: 2, cycles: 2, assembler: "CPX #$FF"},
				{opcode: 0xE4, bytes: 2, cycles: 3, assembler: "CPX $FF"},
				{opcode: 0xEC, bytes: 3, cycles: 4, assembler: "CPX $FFFF"},
			},
		},
		{
			instruction: "CPY",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true},
			want: []expected{
				{opcode: 0xC0, bytes: 2, cycles: 2, assembler: "CPY #$FF"},
				{opcode: 0xC4, bytes: 2, cycles: 3, assembler: "CPY $FF"},
				{opcode: 0xCC, bytes: 3, cycles: 4, assembler: "CPY $FFFF"},
			},
		},
		{
			instruction: "DEC",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xC6, bytes: 2, cycles: 5, assembler: "DEC $FF"},
				{opcode: 0xCE, bytes: 3, cycles: 7, assembler: "DEC $FFFF"},
				{opcode: 0xD6, bytes: 2, cycles: 6, assembler: "DEC $FF,X"},
				{opcode: 0xDE, bytes: 3, cycles: 7, assembler: "DEC $FFFF,X"},
			},
		},
		{
			instruction: "DEX",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xCA, bytes: 1, cycles: 2, assembler: "DEX"},
			},
		},
		{
			instruction: "DEY",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0x88, bytes: 1, cycles: 2, assembler: "DEY"},
			},
		},
		{
			instruction: "EOR",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0x41, bytes: 2, cycles: 6, assembler: "EOR ($FF,X)"},
				{opcode: 0x45, bytes: 2, cycles: 3, assembler: "EOR $FF"},
				{opcode: 0x49, bytes: 2, cycles: 2, assembler: "EOR #$FF"},
				{opcode: 0x4D, bytes: 3, cycles: 4, assembler: "EOR $FFFF"},
				{opcode: 0x51, bytes: 2, cycles: 5, assembler: "EOR ($FF),Y", penalty: true},
				{opcode: 0x55, bytes: 2, cycles: 4, assembler: "EOR $FF,X"},
				{opcode: 0x59, bytes: 3, cycles: 4, assembler: "EOR $FFFF,Y", penalty: true},
				{opcode: 0x5D, bytes: 3, cycles: 4, assembler: "EOR $FFFF,X", penalty: true},
			}},
		{
			instruction: "INC",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xE6, bytes: 2, cycles: 5, assembler: "INC $FF"},
				{opcode: 0xEE, bytes: 3, cycles: 7, assembler: "INC $FFFF"},
				{opcode: 0xF6, bytes: 2, cycles: 6, assembler: "INC $FF,X"},
				{opcode: 0xFE, bytes: 3, cycles: 7, assembler: "INC $FFFF,X"},
			},
		},
		{
			instruction: "INX",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xE8, bytes: 1, cycles: 2, assembler: "INX"},
			},
		},
		{
			instruction: "INY",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xC8, bytes: 1, cycles: 2, assembler: "INY"},
			},
		},
		{
			instruction: "JMP",
			want: []expected{
				{opcode: 0x4C, bytes: 3, cycles: 3, assembler: "JMP $FFFF"},
				{opcode: 0x6C, bytes: 3, cycles: 5, assembler: "JMP ($FFFF)"},
			},
		},
		{
			instruction: "JSR",
			want: []expected{
				{opcode: 0x20, bytes: 3, cycles: 6, assembler: "JSR $FFFF"},
			},
		},
		{
			instruction: "LDA",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xA1, bytes: 2, cycles: 6, assembler: "LDA ($FF,X)"},
				{opcode: 0xA5, bytes: 2, cycles: 3, assembler: "LDA $FF"},
				{opcode: 0xA9, bytes: 2, cycles: 2, assembler: "LDA #$FF"},
				{opcode: 0xAD, bytes: 3, cycles: 4, assembler: "LDA $FFFF"},
				{opcode: 0xB1, bytes: 2, cycles: 5, assembler: "LDA ($FF),Y", penalty: true},
				{opcode: 0xB5, bytes: 2, cycles: 4, assembler: "LDA $FF,X"},
				{opcode: 0xB9, bytes: 3, cycles: 4, assembler: "LDA $FFFF,Y", penalty: true},
				{opcode: 0xBD, bytes: 3, cycles: 4, assembler: "LDA $FFFF,X", penalty: true},
			},
		},
		{
			instruction: "LDX",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xA2, bytes: 2, cycles: 2, assembler: "LDX #$FF"},
				{opcode: 0xA6, bytes: 2, cycles: 3, assembler: "LDX $FF"},
				{opcode: 0xAE, bytes: 3, cycles: 4, assembler: "LDX $FFFF"},
				{opcode: 0xB6, bytes: 2, cycles: 4, assembler: "LDX $FF,Y"},
				{opcode: 0xBE, bytes: 3, cycles: 4, assembler: "LDX $FFFF,Y", penalty: true},
			},
		},
		{
			instruction: "LDY",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xA0, bytes: 2, cycles: 2, assembler: "LDY #$FF"},
				{opcode: 0xA4, bytes: 2, cycles: 3, assembler: "LDY $FF"},
				{opcode: 0xAC, bytes: 3, cycles: 4, assembler: "LDY $FFFF"},
				{opcode: 0xB4, bytes: 2, cycles: 4, assembler: "LDY $FF,X"},
				{opcode: 0xBC, bytes: 3, cycles: 4, assembler: "LDY $FFFF,X", penalty: true},
			},
		},
		{
			instruction: "LSR",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true},
			want: []expected{
				{opcode: 0x46, bytes: 2, cycles: 5, assembler: "LSR $FF"},
				{opcode: 0x4A, bytes: 1, cycles: 2, assembler: "LSR A"},
				{opcode: 0x4E, bytes: 3, cycles: 6, assembler: "LSR $FFFF"},
				{opcode: 0x56, bytes: 2, cycles: 6, assembler: "LSR $FF,X"},
				{opcode: 0x5E, bytes: 3, cycles: 7, assembler: "LSR $FFFF,X"},
			},
		},
		{
			instruction: "NOP",
			want: []expected{
				{opcode: 0xEA, bytes: 1, cycles: 2, assembler: "NOP"},
			},
		},
		{
			instruction: "ORA",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0x01, bytes: 2, cycles: 6, assembler: "ORA ($FF,X)"},
				{opcode: 0x05, bytes: 2, cycles: 3, assembler: "ORA $FF"},
				{opcode: 0x09, bytes: 2, cycles: 2, assembler: "ORA #$FF"},
				{opcode: 0x0D, bytes: 3, cycles: 4, assembler: "ORA $FFFF"},
				{opcode: 0x11, bytes: 2, cycles: 5, assembler: "ORA ($FF),Y", penalty: true},
				{opcode: 0x15, bytes: 2, cycles: 4, assembler: "ORA $FF,X"},
				{opcode: 0x19, bytes: 3, cycles: 4, assembler: "ORA $FFFF,Y", penalty: true},
				{opcode: 0x1D, bytes: 3, cycles: 4, assembler: "ORA $FFFF,X", penalty: true},
			},
		},
		{
			instruction: "PHA",
			want: []expected{
				{opcode: 0x48, bytes: 1, cycles: 3, assembler: "PHA"},
			},
		},
		{
			instruction: "PHP",
			want: []expected{
				{opcode: 0x08, bytes: 1, cycles: 3, assembler: "PHP"},
			},
		},
		{
			instruction: "PLA",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0x68, bytes: 1, cycles: 4, assembler: "PLA"},
			},
		},
		{
			instruction: "PLP",
			want: []expected{
				{opcode: 0x28, bytes: 1, cycles: 4, assembler: "PLP"},
			},
		},
		{
			instruction: "ROL",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true},
			want: []expected{
				{opcode: 0x26, bytes: 2, cycles: 5, assembler: "ROL $FF"},
				{opcode: 0x2A, bytes: 1, cycles: 2, assembler: "ROL A"},
				{opcode: 0x2E, bytes: 3, cycles: 6, assembler: "ROL $FFFF"},
				{opcode: 0x36, bytes: 2, cycles: 6, assembler: "ROL $FF,X"},
				{opcode: 0x3E, bytes: 3, cycles: 7, assembler: "ROL $FFFF,X"},
			},
		},
		{
			instruction: "ROR",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true},
			want: []expected{
				{opcode: 0x66, bytes: 2, cycles: 5, assembler: "ROR $FF"},
				{opcode: 0x6A, bytes: 1, cycles: 2, assembler: "ROR A"},
				{opcode: 0x6E, bytes: 3, cycles: 6, assembler: "ROR $FFFF"},
				{opcode: 0x76, bytes: 2, cycles: 6, assembler: "ROR $FF,X"},
				{opcode: 0x7E, bytes: 3, cycles: 7, assembler: "ROR $FFFF,X"},
			},
		},
		{
			instruction: "RTI",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true, Overflow: true, Interrupt: true, Break: true},
			want: []expected{
				{opcode: 0x40, bytes: 1, cycles: 6, assembler: "RTI"},
			},
		},
		{
			instruction: "RTS",
			want: []expected{
				{opcode: 0x60, bytes: 1, cycles: 6, assembler: "RTS"},
			},
		},
		{
			instruction: "SBC",
			wantFlags:   Flags{Negative: true, Zero: true, Carry: true, Overflow: true},
			want: []expected{
				{opcode: 0xE1, bytes: 2, cycles: 6, assembler: "SBC ($FF,X)"},
				{opcode: 0xE5, bytes: 2, cycles: 3, assembler: "SBC $FF"},
				{opcode: 0xE9, bytes: 2, cycles: 2, assembler: "SBC #$FF"},
				{opcode: 0xED, bytes: 3, cycles: 4, assembler: "SBC $FFFF"},
				{opcode: 0xF1, bytes: 2, cycles: 5, assembler: "SBC ($FF),Y", penalty: true},
				{opcode: 0xF5, bytes: 2, cycles: 4, assembler: "SBC $FF,X"},
				{opcode: 0xF9, bytes: 3, cycles: 4, assembler: "SBC $FFFF,Y", penalty: true},
				{opcode: 0xFD, bytes: 3, cycles: 4, assembler: "SBC $FFFF,X", penalty: true},
			},
		},
		{
			instruction: "SEC",
			wantFlags:   Flags{Carry: true},
			want: []expected{
				{opcode: 0x38, bytes: 1, cycles: 2, assembler: "SEC"},
			},
		},
		{
			instruction: "SED",
			wantFlags:   Flags{Decimal: true},
			want: []expected{
				{opcode: 0xF8, bytes: 1, cycles: 2, assembler: "SED"},
			},
		},
		{
			instruction: "SEI",
			wantFlags:   Flags{Interrupt: true},
			want: []expected{
				{opcode: 0x78, bytes: 1, cycles: 2, assembler: "SEI"},
			},
		},
		{
			instruction: "STA",
			want: []expected{
				{opcode: 0x81, bytes: 2, cycles: 6, assembler: "STA ($FF,X)"},
				{opcode: 0x85, bytes: 2, cycles: 3, assembler: "STA $FF"},
				{opcode: 0x8D, bytes: 3, cycles: 4, assembler: "STA $FFFF"},
				{opcode: 0x91, bytes: 2, cycles: 6, assembler: "STA ($FF),Y"},
				{opcode: 0x95, bytes: 2, cycles: 4, assembler: "STA $FF,X"},
				{opcode: 0x99, bytes: 3, cycles: 5, assembler: "STA $FFFF,Y"},
				{opcode: 0x9D, bytes: 3, cycles: 5, assembler: "STA $FFFF,X"},
			},
		},
		{
			instruction: "STX",
			want: []expected{
				{opcode: 0x86, bytes: 2, cycles: 3, assembler: "STX $FF"},
				{opcode: 0x8E, bytes: 3, cycles: 4, assembler: "STX $FFFF"},
				{opcode: 0x96, bytes: 2, cycles: 4, assembler: "STX $FF,Y"},
			},
		},
		{
			instruction: "STY",
			want: []expected{
				{opcode: 0x84, bytes: 2, cycles: 3, assembler: "STY $FF"},
				{opcode: 0x8C, bytes: 3, cycles: 4, assembler: "STY $FFFF"},
				{opcode: 0x94, bytes: 2, cycles: 4, assembler: "STY $FF,X"},
			},
		},
		{
			instruction: "TAX",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xAA, bytes: 1, cycles: 2, assembler: "TAX"},
			},
		},
		{
			instruction: "TAY",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xA8, bytes: 1, cycles: 2, assembler: "TAY"},
			},
		},
		{
			instruction: "TSX",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0xBA, bytes: 1, cycles: 2, assembler: "TSX"},
			},
		},
		{
			instruction: "TXA",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0x8A, bytes: 1, cycles: 2, assembler: "TXA"},
			},
		},
		{
			instruction: "TXS",
			want: []expected{
				{opcode: 0x9A, bytes: 1, cycles: 2, assembler: "TXS"},
			},
		},
		{
			instruction: "TYA",
			wantFlags:   Flags{Negative: true, Zero: true},
			want: []expected{
				{opcode: 0x98, bytes: 1, cycles: 2, assembler: "TYA"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.instruction, func(t *testing.T) {
			mnemonics, err := MnemonicsFromInstructionName(tt.instruction)
			if err != nil {
				t.Errorf(err.Error())
				return
			}

			if len(mnemonics) != len(tt.want) {
				t.Errorf("Unexpected number of mnemonics, got = %v, want = %v", len(mnemonics), len(tt.want))
			}

			details := make([]expected, len(mnemonics), len(mnemonics))
			for i, mnemonic := range mnemonics {

				detail := NewMnemonicDisplayDetails(mnemonic)

				// Replace the opcode holder with a hardcoded value of 0xFFFF or 0xFF
				if strings.Contains(detail.Assembler, "%04X") {
					detail.Assembler = fmt.Sprintf(detail.Assembler, 0xFFFF)
				}
				if strings.Contains(detail.Assembler, "%02X") {
					detail.Assembler = fmt.Sprintf(detail.Assembler, 0xFF)
				}

				if i < len(tt.want) {
					// Copy the wanted flags across.
					tt.want[i].flags = tt.wantFlags
					// We fill in this textual description as we are not super interested in it
					tt.want[i].addressing = detail.Addressing
				}

				details[i].addressing = detail.Addressing
				details[i].assembler = detail.Assembler
				details[i].opcode = detail.Opcode
				details[i].bytes = detail.Bytes
				details[i].cycles = detail.Cycles
				details[i].penalty = detail.PageBoundaryPenalty
				details[i].flags = detail.AffectedFlags
			}

			if !reflect.DeepEqual(details, tt.want) {
				t.Errorf("Unexpected instruction results got")
				t.Errorf("  Mnemonic: %v", tt.instruction)
				t.Errorf("  %-14v  %6v  %5v  %6v  %7v  %-24v  %v", "Assembler", "Opcode", "Bytes", "Cycles", "Penalty", "Addressing", "Flags")
				for _, d := range details {
					t.Errorf("  %-14v    0x%02X  %5v  %6v  %7v  %-24v  %v", d.assembler, d.opcode, d.bytes, d.cycles, d.penalty, d.addressing, d.flags)
				}
				t.Errorf("\nUnexpected instruction results want")
				t.Errorf("  Mnemonic: %v", tt.instruction)
				t.Errorf("  %-14v  %6v  %5v  %6v  %7v  %-24v  %v", "Assembler", "Opcode", "Bytes", "Cycles", "Penalty", "Addressing", "Flags")
				for _, w := range tt.want {
					t.Errorf("  %-14v    0x%02X  %5v  %6v  %7v  %-24v  %v", w.assembler, w.opcode, w.bytes, w.cycles, w.penalty, w.addressing, w.flags)
				}
			}
		})
	}
}

func TestCycleCounts(t *testing.T) {

	// Generate a full instruction set of known opcodes for this test.
	instructions := make(Instructions, 0, 0xFF)

	for _, mnemonic := range AllOpcodes() {
		instruction := NewInstruction(mnemonic)
		instructions = append(instructions, instruction)
	}

	is, err := NewInstructionSet(instructions)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		wantCycles uint
		testCpuMethodConfig
	}{
		{
			wantCycles: 1,
			testCpuMethodConfig: testCpuMethodConfig{
				name:      "START with ADC and test each cycle count.",
				startRam:  []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantRam:   []uint8{0, 0, 0, 0, 0, 0, 0, 0},
				wantState: State{PC: 1, SP: StackPointerStart},
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
					t.Errorf("TestCycleCounts() did not execute the expected number of cycles, got = %v, want = %v", cycles, tt.wantCycles)
				}
				return err
			}

			tt.instructionSet = is
			testCpuMethod(t, tt.testCpuMethodConfig, callback)
		})
	}
}
