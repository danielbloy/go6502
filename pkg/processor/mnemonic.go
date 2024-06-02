package processor

import (
	"fmt"
	"sort"
	"strings"
)

// NewInstruction Builds an Instruction for the given opcode from Mnemonic data
// (which can come from largely static data).
func NewInstruction(mnemonic Mnemonic) Instruction {

	// 1 is removed from the operation cycles because the CPU uses 1 cycle to read the opcode from memory.
	cycles := int(mnemonic.Operation.Cycles)
	if cycles > 0 {
		cycles--
	}
	cycles += int(mnemonic.Addressing.Cycles)
	cycles += mnemonic.CycleAdjust

	result := Instruction{
		Opcode:              mnemonic.Opcode,
		AddressingFunc:      mnemonic.Addressing.AddressingFunc,
		Operation:           mnemonic.Operation.Operation,
		Cycles:              uint(cycles),
		PageBoundaryPenalty: mnemonic.Operation.PageBoundaryPenalty && mnemonic.Addressing.PageBoundaryPenalty,
	}

	return result
}

// Mnemonic combines an operation code with an operation and addressing mode as well as additional
// timing information. This data can be used to generate Instructions.
type Mnemonic struct {
	Opcode      Opcode
	Operation   MnemonicOperation
	Addressing  MnemonicAddressingMode
	CycleAdjust int // The number of cycles to adjust the combination from Operation and Addressing.
}

// MnemonicOperation details from:
//   - http://www.6502.org/tutorials/6502opcodes.html
//   - https://www.masswerk.at/6502/6502_instruction_set.html
//
// The contents of this file are intended to be static data representing all the
// known supported operation codes for the plain old vanilla 6502.
type MnemonicOperation struct {
	Name                 string
	Description          string
	AssemblyLanguageForm string
	AffectedFlags        Flags
	Bytes                uint // The number of bytes for the Opcode (but not addressing), usually 1.
	Cycles               uint // The number of cycles for the Operation (but not addressing), usually 1.
	PageBoundaryPenalty  bool // See note below
	// TODO BranchTakenPenalty   bool // See note below
	Operation Operation

	// NOTE: If PageBoundaryPenalty is true and the corresponding PageBoundaryPenalty value in the
	//       MnemonicAddressingMode is true then the operation is susceptible to a page boundary penalty.
	// NOTE: If BranchTakenPenalty is true and the branch was taken, then an additional 1 cycle
	//       penalty is incurred in the CPU.
}

type MnemonicAddressingMode struct {
	Name                 string
	AssemblyLanguageForm string
	Bytes                uint // The number of bytes following the opcode required for addressing.
	Cycles               uint // The number of cycles required for the addressing; not including any penalties.
	PageBoundaryPenalty  bool // See note above
	AddressingFunc       AddressingFunc
}

// MnemonicFromOpCode searches through all known mnemonics and returns the
// Mnemonic that contains the given opcode. If the opcode cannot be found
// then an error is returned.
func MnemonicFromOpCode(opcode Opcode) (Mnemonic, error) {
	mnemonic, ok := mnemonics[opcode]
	if !ok {
		return Mnemonic{}, fmt.Errorf("the operation code %20x cannot be found", opcode)
	}
	return mnemonic, nil
}

// MnemonicsFromInstructionName returns all of the mnemonics for that given
// assembly language name. The results are returned in opcode order, smallest
// to highest. The comparison is case insensitive.
func MnemonicsFromInstructionName(name string) ([]Mnemonic, error) {
	result := make([]Mnemonic, 0)
	for _, mnemonic := range mnemonics {
		if strings.EqualFold(mnemonic.Operation.AssemblyLanguageForm, name) {
			result = append(result, mnemonic)
		}
	}
	if len(result) == 0 {
		return result, fmt.Errorf("could not find any mnemonics named `%v`", name)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Opcode < result[j].Opcode })
	return result, nil
}

// AllOpcodes returns Mnemonic representations of all the known opcodes.
//
// NOTE: Currently this is only the known legal opcodes.
func AllOpcodes() []Mnemonic {

	// Here we make a copy of the data to avoid it being mutated.
	result := make([]Mnemonic, len(opcodes))
	for i, opcode := range opcodes {
		result[i] = opcode
	}

	return result
}

var (
	mnemonics map[Opcode]Mnemonic
)

// Populates legal Mnemonics from the pre-defined array of opcodes.
func init() {

	mnemonics = make(map[Opcode]Mnemonic, len(opcodes))
	for _, opcode := range opcodes {
		mnemonics[opcode.Opcode] = opcode
	}
}

// MnemonicDisplayDetails flattens the Mnemonic details into
// a simple structure that can be used to display the single
// instructions details. The tables this can be used to
// generate are expected to look like this:
// https://www.masswerk.at/6502/6502_instruction_set.html#BPL
type MnemonicDisplayDetails struct {
	Name                string
	Description         string
	Mnemonic            string
	AffectedFlags       Flags
	Addressing          string
	Assembler           string
	Opcode              Opcode
	Bytes               uint
	Cycles              uint
	PageBoundaryPenalty bool
}

// NewMnemonicDisplayDetails converts a Mnemonic instance into MnemonicDisplayDetails which
// is designed to be presented to a human.
func NewMnemonicDisplayDetails(mnemonic Mnemonic) MnemonicDisplayDetails {
	instruction := NewInstruction(mnemonic)
	// The TrimSpace call is needed to remove any additional whitespace at the end of the Assembler
	// which will occur in cases where the addressing does not have an assembly language form such
	// as the implied addressing form.
	assembler := strings.TrimSpace(
		fmt.Sprintf("%v %v", mnemonic.Operation.AssemblyLanguageForm, mnemonic.Addressing.AssemblyLanguageForm))

	result := MnemonicDisplayDetails{
		Name:          mnemonic.Operation.Name,
		Description:   mnemonic.Operation.Description,
		Mnemonic:      mnemonic.Operation.AssemblyLanguageForm,
		AffectedFlags: mnemonic.Operation.AffectedFlags,
		Addressing:    mnemonic.Addressing.Name,
		Assembler:     assembler,
		Opcode:        mnemonic.Opcode,
		Bytes:         mnemonic.Operation.Bytes + mnemonic.Addressing.Bytes,
		// The +1 is because the instruction object does not include the cycle the CPU
		// uses to read the opcode but the as it is subtracted in NewInstruction().
		Cycles:              instruction.Cycles + 1,
		PageBoundaryPenalty: instruction.PageBoundaryPenalty,
	}

	return result
}

var (
	Abs = MnemonicAddressingMode{
		Name:                 "Absolute",
		AssemblyLanguageForm: "$%04X",
		Bytes:                2,
		Cycles:               3,
		AddressingFunc:       Absolute,
	}
	AbsX = MnemonicAddressingMode{
		Name:                 "Absolute, X-indexed",
		AssemblyLanguageForm: "$%04X,X",
		Bytes:                2,
		Cycles:               3,
		PageBoundaryPenalty:  true,
		AddressingFunc:       AbsoluteX,
	}
	AbsY = MnemonicAddressingMode{
		Name:                 "Absolute, Y-indexed",
		AssemblyLanguageForm: "$%04X,Y",
		Bytes:                2,
		Cycles:               3,
		PageBoundaryPenalty:  true,
		AddressingFunc:       AbsoluteY,
	}
	Acc = MnemonicAddressingMode{
		Name:                 "Accumulator",
		AssemblyLanguageForm: "A",
		Bytes:                0,
		Cycles:               0,
		AddressingFunc:       Accumulator,
	}
	Imm = MnemonicAddressingMode{
		Name:                 "Immediate",
		AssemblyLanguageForm: "#$%02X",
		Bytes:                1,
		Cycles:               1,
		AddressingFunc:       Immediate,
	}
	Imp = MnemonicAddressingMode{
		Name:                 "Implied",
		AssemblyLanguageForm: "",
		Bytes:                0,
		Cycles:               0,
		AddressingFunc:       Implied,
	}
	Ind = MnemonicAddressingMode{
		Name:                 "Indirect",
		AssemblyLanguageForm: "($%04X)",
		Bytes:                2,
		Cycles:               4,
		AddressingFunc:       Indirect,
	}
	IndX = MnemonicAddressingMode{
		Name:                 "X-indexed, Indirect",
		AssemblyLanguageForm: "($%02X,X)",
		Bytes:                1,
		Cycles:               5,
		AddressingFunc:       IndirectX,
	}
	IndY = MnemonicAddressingMode{
		Name:                 "Indirect, Y-indexed",
		AssemblyLanguageForm: "($%02X),Y",
		Bytes:                1,
		Cycles:               4,
		PageBoundaryPenalty:  true,
		AddressingFunc:       IndirectY,
	}
	Rel = MnemonicAddressingMode{
		Name:                 "Relative",
		AssemblyLanguageForm: "$%02X",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		AddressingFunc:       Relative,
	}
	Zpg = MnemonicAddressingMode{
		Name:                 "Zero Page",
		AssemblyLanguageForm: "$%02X",
		Bytes:                1,
		Cycles:               2,
		AddressingFunc:       ZeroPage,
	}
	ZpgX = MnemonicAddressingMode{
		Name:                 "Zero Page, X-indexed",
		AssemblyLanguageForm: "$%02X,X",
		Bytes:                1,
		Cycles:               3,
		AddressingFunc:       ZeroPageX,
	}
	ZpgY = MnemonicAddressingMode{
		Name:                 "Zero Page, Y-indexed",
		AssemblyLanguageForm: "$%02X,Y",
		Bytes:                1,
		Cycles:               3,
		AddressingFunc:       ZeroPageY,
	}
)

var (
	Adc = MnemonicOperation{
		Name:                 "Add with carry",
		Description:          "Add memory to accumulator with carry. Results are dependant on the setting of the decimal flag. In decimal mode, addition is carried out on the assumption that the values involved are packed BCD (Binary Coded Decimal). There is no way to add without carry.",
		AssemblyLanguageForm: "ADC",
		AffectedFlags:        Flags{Negative: true, Overflow: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            AddWithCarry,
	}
	And = MnemonicOperation{
		Name:                 "Bitwise AND with A",
		Description:          "Perform a bitwise AND of the accumulator with a value in memory.",
		AssemblyLanguageForm: "AND",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            AndWithA,
	}
	Asl = MnemonicOperation{
		Name:                 "Arithmetic shift left",
		Description:          "ASL shifts all bits left one position. 0 is shifted into bit 0 and the original bit 7 is shifted into the Carry.",
		AssemblyLanguageForm: "ASL",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            ArithmeticShiftLeft,
	}
	Bit = MnemonicOperation{
		Name:                 "Test bits",
		Description:          "BIT sets the Z flag as though the value in the address tested were ANDed with the accumulator. The N and V flags are set to match bits 7 and 6 respectively in the value stored at the tested address.",
		AssemblyLanguageForm: "BIT",
		AffectedFlags:        Flags{Negative: true, Overflow: true, Zero: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  false,
		Operation:            TestBitsInMemoryWithAccumulator,
	}
	Bcc = MnemonicOperation{
		Name:                 "Branch on carry clear",
		Description:          "Branch relative to current program counter if the carry flag is clear (i.e. carry = 0).",
		AssemblyLanguageForm: "BCC",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnCarryClear,
	}
	Bcs = MnemonicOperation{
		Name:                 "Branch on carry set",
		Description:          "Branch relative to current program counter if the carry flag is set (i.e. carry = 1).",
		AssemblyLanguageForm: "BCS",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnCarrySet,
	}
	Beq = MnemonicOperation{
		Name:                 "Branch on equal (result zero)",
		Description:          "Branch relative to current program counter if the zero flag is clear (i.e. zero = 0).",
		AssemblyLanguageForm: "BEQ",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnEqual,
	}
	Bmi = MnemonicOperation{
		Name:                 "Branch on result negative",
		Description:          "Branch relative to current program counter if the negative flag is set (i.e. negative = 1).",
		AssemblyLanguageForm: "BMI",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnMinus,
	}
	Bne = MnemonicOperation{
		Name:                 "Branch on not equal (result not zero)",
		Description:          "Branch relative to current program counter if the zero flag is set (i.e. zero = 1).",
		AssemblyLanguageForm: "BNE",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnNotEqual,
	}
	Bpl = MnemonicOperation{
		Name:                 "Branch on result plus",
		Description:          "Branch relative to current program counter if the negative flag is clear (i.e. negative = 0).",
		AssemblyLanguageForm: "BPL",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnPlus,
	}
	Brk = MnemonicOperation{
		Name:                 "Break",
		Description:          "BRK causes a non-maskable interrupt and increments the program counter by one. Therefore an RTI will go to the address of the BRK +2 so that BRK may be used to replace a two-byte instruction for debugging and the subsequent RTI will be correct.",
		AssemblyLanguageForm: "BRK",
		AffectedFlags:        Flags{Interrupt: true, Break: true},
		Bytes:                1,
		Cycles:               7,
		PageBoundaryPenalty:  false,
		Operation:            Break,
	}
	Bvc = MnemonicOperation{
		Name:                 "Branch on overflow clear",
		Description:          "Branch relative to current program counter if the overflow flag is clear (i.e. overflow = 0).",
		AssemblyLanguageForm: "BVC",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnOverflowClear,
	}
	Bvs = MnemonicOperation{
		Name:                 "Branch on overflow set",
		Description:          "Branch relative to current program counter if the overflow flag is set (i.e. overflow = 1).",
		AssemblyLanguageForm: "BVS",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true, // TODO: Needs an additional penalty cycle if branch taken
		Operation:            BranchOnOverflowSet,
	}
	Clc = MnemonicOperation{
		Name:                 "Clear carry flag",
		Description:          "Clears the carry flag (i.e. carry = 0).",
		AssemblyLanguageForm: "CLC",
		AffectedFlags:        Flags{Carry: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            ClearCarry,
	}
	Cld = MnemonicOperation{
		Name:                 "Clear decimal flag",
		Description:          "Clears the decimal mode flag (i.e. decimal = 0).",
		AssemblyLanguageForm: "CLD",
		AffectedFlags:        Flags{Decimal: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            ClearDecimal,
	}
	Cli = MnemonicOperation{
		Name:                 "Clear interrupt flag",
		Description:          "Clear the interrupt disable flag (i.e. interrupt = 0).",
		AssemblyLanguageForm: "CLI",
		AffectedFlags:        Flags{Interrupt: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            ClearInterrupt,
	}
	Clv = MnemonicOperation{
		Name:                 "Clear overflow flag",
		Description:          "Clears the overflow flag (i.e. overflow = 0).",
		AssemblyLanguageForm: "CLV",
		AffectedFlags:        Flags{Overflow: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            ClearOverflow,
	}
	Cmp = MnemonicOperation{
		Name: "Compare memory with A",
		Description: "Sets the flags as if a subtraction has been carried out. If the value in the accumulator is equal " +
			"or greater than the compared value, the carry will be set (i.e. carry = 1). The equal (zero) and negative flags " +
			"will be set based on equality or lack thereof and the sign (i.e. A >= $80) of the accumulator.",
		AssemblyLanguageForm: "CMP",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            CompareWithA,
	}
	Cpx = MnemonicOperation{
		Name: "Compare memory with X",
		Description: "Sets the flags as if a subtraction has been carried out. If the value in the X is equal " +
			"or greater than the compared value, the carry will be set (i.e. carry = 1). The equal (zero) and negative flags " +
			"will be set based on equality or lack thereof and the sign (i.e. X >= $80) of the X register.",
		AssemblyLanguageForm: "CPX",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  false,
		Operation:            CompareWithX,
	}
	Cpy = MnemonicOperation{
		Name: "Compare memory with Y",
		Description: "Sets the flags as if a subtraction has been carried out. If the value in the Y is equal " +
			"or greater than the compared value, the carry will be set (i.e. carry = 1). The equal (zero) and negative flags " +
			"will be set based on equality or lack thereof and the sign (i.e. X >= $80) of the Y register.",
		AssemblyLanguageForm: "CPY",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            CompareWithY,
	}
	Dec = MnemonicOperation{
		Name:                 "Decrement memory",
		Description:          "Decrement memory by 1.",
		AssemblyLanguageForm: "DEC",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               3,
		PageBoundaryPenalty:  false,
		Operation:            Decrement,
	}
	Dex = MnemonicOperation{
		Name:                 "Decrement X",
		Description:          "Decrement X by 1.",
		AssemblyLanguageForm: "DEX",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            DecrementX,
	}
	Dey = MnemonicOperation{
		Name:                 "Decrement Y",
		Description:          "Decrement Y by 1.",
		AssemblyLanguageForm: "DEY",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            DecrementY,
	}
	Eor = MnemonicOperation{
		Name:                 "Bitwise exclusive OR with A",
		Description:          "Perform a bitwise exclusive OR of the accumulator with a value in memory.",
		AssemblyLanguageForm: "EOR",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            ExclusiveOrWithA,
	}
	Inc = MnemonicOperation{
		Name:                 "Increment memory",
		Description:          "Increment memory by 1.",
		AssemblyLanguageForm: "INC",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               3,
		PageBoundaryPenalty:  false,
		Operation:            Increment,
	}
	Inx = MnemonicOperation{
		Name:                 "Increment X",
		Description:          "Increment X by 1.",
		AssemblyLanguageForm: "INX",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            IncrementX,
	}
	Iny = MnemonicOperation{
		Name:                 "Increment Y",
		Description:          "Increment Y by 1.",
		AssemblyLanguageForm: "INY",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            IncrementY,
	}
	Jmp = MnemonicOperation{
		Name:                 "Jump",
		Description:          "Jump to a new memory location.",
		AssemblyLanguageForm: "JMP",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  false,
		Operation:            Jump,
	}
	Jsr = MnemonicOperation{
		Name:                 "Jump to subroutine",
		Description:          "Jump to subroutine, saving the return address on the stack.",
		AssemblyLanguageForm: "JSR",
		Bytes:                1,
		Cycles:               3,
		PageBoundaryPenalty:  false,
		Operation:            JumpSubRoutine,
	}
	Lda = MnemonicOperation{
		Name:                 "Load A",
		Description:          "Load accumulator with value.",
		AssemblyLanguageForm: "LDA",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            LoadA,
	}
	Ldx = MnemonicOperation{
		Name:                 "Load X",
		Description:          "Load X with value.",
		AssemblyLanguageForm: "LDX",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            LoadX,
	}
	Ldy = MnemonicOperation{
		Name:                 "Load Y",
		Description:          "Load Y with value.",
		AssemblyLanguageForm: "LDY",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            LoadY,
	}
	Lsr = MnemonicOperation{
		Name:                 "Logical shift right",
		Description:          "Shifts all bits right one position. Zero is shifted into bit 7 and bit 0 is shifted into the carry flag.",
		AssemblyLanguageForm: "LSR",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            LogicalShiftRight,
	}
	Nop = MnemonicOperation{
		Name:                 "No Operation",
		Description:          "NOP can be used to reserve space for future modifications or adjust timing loops.",
		AssemblyLanguageForm: "NOP",
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            NoOperation,
	}
	Ora = MnemonicOperation{
		Name:                 "Bitwise OR with A",
		Description:          "Perform a bitwise OR of the accumulator with a value in memory.",
		AssemblyLanguageForm: "ORA",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            OrWithA,
	}
	Pha = MnemonicOperation{
		Name:                 "Push A on stack",
		Description:          "Push the accumulator onto the top of the stack.",
		AssemblyLanguageForm: "PHA",
		Bytes:                1,
		Cycles:               3,
		PageBoundaryPenalty:  false,
		Operation:            PushA,
	}
	Php = MnemonicOperation{
		Name:                 "Push P on stack",
		Description:          "Push the processor status flag (with the break and constant bits set to 1) onto the top of the stack.",
		AssemblyLanguageForm: "PHP",
		Bytes:                1,
		Cycles:               3,
		PageBoundaryPenalty:  false,
		Operation:            PushP,
	}
	Pla = MnemonicOperation{
		Name:                 "Pull a from stack",
		Description:          "Pull the accumulator from the top of the stack.",
		AssemblyLanguageForm: "PLA",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               4,
		PageBoundaryPenalty:  false,
		Operation:            PullA,
	}
	Plp = MnemonicOperation{
		Name:                 "Pull P from stack",
		Description:          "Pull the processor status flag (with the break and constant bits ignored) from the top of the stack.",
		AssemblyLanguageForm: "PLP",
		Bytes:                1,
		Cycles:               4,
		PageBoundaryPenalty:  false,
		Operation:            PullP,
	}
	Rol = MnemonicOperation{
		Name:                 "Rotate left",
		Description:          "Rotates all bits left one position. The carry flag is shifted into bit 0 and bit 7 is shifted into the carry flag.",
		AssemblyLanguageForm: "ROL",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            RotateLeft,
	}
	Ror = MnemonicOperation{
		Name:                 "Rotate right",
		Description:          "Rotates all bits right one position. The carry flag is shifted into bit 7 and bit 0 is shifted into the carry flag.",
		AssemblyLanguageForm: "ROR",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            RotateRight,
	}
	Rti = MnemonicOperation{
		Name:                 "Return from interrupt",
		Description:          "The status register is pulled with the break and constant bits ignored (constant is forced on and break is forced clear). Then program counter is then pulled from the stack.",
		AssemblyLanguageForm: "RTI",
		AffectedFlags:        Flags{Negative: true, Zero: true, Carry: true, Overflow: true, Interrupt: true, Break: true},
		Bytes:                1,
		Cycles:               6,
		PageBoundaryPenalty:  false,
		Operation:            ReturnFromInterrupt,
	}
	Rts = MnemonicOperation{
		Name:                 "Return from subroutine",
		Description:          "Pulls the top two bytes from the top of the stack and transfers control to that address + 1.",
		AssemblyLanguageForm: "RTS",
		Bytes:                1,
		Cycles:               6,
		PageBoundaryPenalty:  false,
		Operation:            ReturnFromSubroutine,
	}
	Sbc = MnemonicOperation{
		Name:                 "Subtract with borrow",
		Description:          "Subtract memory from accumulator with borrow. Results are dependant on the setting of the decimal flag. In decimal mode, subtraction is carried out on the assumption that the values involved are packed BCD (Binary Coded Decimal). There is no way to subtract without borrow.",
		AssemblyLanguageForm: "SBC",
		AffectedFlags:        Flags{Negative: true, Overflow: true, Zero: true, Carry: true},
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  true,
		Operation:            SubtractWithCarry,
	}
	Sec = MnemonicOperation{
		Name:                 "Set carry flag",
		Description:          "Sets the carry flag (i.e. carry = 1).",
		AssemblyLanguageForm: "SEC",
		AffectedFlags:        Flags{Carry: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            SetCarry,
	}
	Sed = MnemonicOperation{
		Name:                 "Set decimal flag",
		Description:          "Sets the decimal mode flag (i.e. decimal = 1).",
		AssemblyLanguageForm: "SED",
		AffectedFlags:        Flags{Decimal: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            SetDecimal,
	}
	Sei = MnemonicOperation{
		Name:                 "Set interrupt flag",
		Description:          "Sets the interrupt disable flag (i.e. interrupt = 1).",
		AssemblyLanguageForm: "SEI",
		AffectedFlags:        Flags{Interrupt: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            SetInterrupt,
	}
	Sta = MnemonicOperation{
		Name:                 "Store A",
		Description:          "Store the accumulator to memory.",
		AssemblyLanguageForm: "STA",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  false,
		Operation:            StoreA,
	}
	Stx = MnemonicOperation{
		Name:                 "Store X",
		Description:          "Store the X register to memory.",
		AssemblyLanguageForm: "STX",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  false,
		Operation:            StoreX,
	}
	Sty = MnemonicOperation{
		Name:                 "Store Y",
		Description:          "Store the Y register to memory.",
		AssemblyLanguageForm: "STY",
		Bytes:                1,
		Cycles:               1,
		PageBoundaryPenalty:  false,
		Operation:            StoreY,
	}
	Tax = MnemonicOperation{
		Name:                 "Transfer A to X",
		Description:          "Transfer the value in the accumulator to X.",
		AssemblyLanguageForm: "TAX",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            TransferAtoX,
	}
	Tay = MnemonicOperation{
		Name:                 "Transfer A to Y",
		Description:          "Transfer the value in the accumulator to Y.",
		AssemblyLanguageForm: "TAY",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            TransferAtoY,
	}
	Tsx = MnemonicOperation{
		Name:                 "Transfer SP to X",
		Description:          "Transfer the value in the stack pointer to X.",
		AssemblyLanguageForm: "TSX",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            TransferSPtoX,
	}
	Txa = MnemonicOperation{
		Name:                 "Transfer X to A",
		Description:          "Transfer the value in X to the accumulator.",
		AssemblyLanguageForm: "TXA",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            TransferXtoA,
	}
	Txs = MnemonicOperation{
		Name:                 "Transfer X to SP",
		Description:          "Transfer the value in X to the stack pointer.",
		AssemblyLanguageForm: "TXS",
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            TransferXtoSP,
	}
	Tya = MnemonicOperation{
		Name:                 "Transfer Y to A",
		Description:          "Transfer the value in Y to the accumulator.",
		AssemblyLanguageForm: "TYA",
		AffectedFlags:        Flags{Negative: true, Zero: true},
		Bytes:                1,
		Cycles:               2,
		PageBoundaryPenalty:  false,
		Operation:            TransferYtoA,
	}
)

// Data from the following sources:
//   - http://www.6502.org/tutorials/6502opcodes.html
//   - https://www.masswerk.at/6502/6502_instruction_set.html
//   - https://dwheeler.com/6502/oneelkruns/asm1step.html
var opcodes = []Mnemonic{
	/*
	   ADC
	   addressing    assembler      opc  bytes  cycles
	   (indirect,X)  ADC ($FF,X)    61    2      6
	   zeropage      ADC $FF        65    2      3
	   immediate     ADC #$FF       69    2      2
	   absolute      ADC $FFFF      6D    3      4
	   (indirect),Y	 ADC ($FF),Y    71    2      5*
	   zeropage,X    ADC $FF,X      75    2      4
	   absolute,Y    ADC $FFFF,Y    79    3      4*
	   absolute,X    ADC $FFFF,X    7D    3      4*
	*/
	{Opcode: 0x61, Operation: Adc, Addressing: IndX},
	{Opcode: 0x65, Operation: Adc, Addressing: Zpg},
	{Opcode: 0x69, Operation: Adc, Addressing: Imm},
	{Opcode: 0x6D, Operation: Adc, Addressing: Abs},
	{Opcode: 0x71, Operation: Adc, Addressing: IndY},
	{Opcode: 0x75, Operation: Adc, Addressing: ZpgX},
	{Opcode: 0x79, Operation: Adc, Addressing: AbsY},
	{Opcode: 0x7D, Operation: Adc, Addressing: AbsX},
	/*
		AND
		addressing    assembler     opc  bytes  cycles
		(indirect,X)  AND ($FF,X)   21    2      6
		zeropage      AND $FF       25    2      3
		immediate     AND #$FF      29    2      2
		absolute      AND $FFFF     2D    3      4
		(indirect),Y  AND ($FF),Y   31    2      5*
		zeropage,X    AND $FF,X     35    2      4
		absolute,Y    AND $FFFF,Y   39    3      4*
		absolute,X    AND $FFFF,X   3D    3      4*
	*/
	{Opcode: 0x21, Operation: And, Addressing: IndX},
	{Opcode: 0x25, Operation: And, Addressing: Zpg},
	{Opcode: 0x29, Operation: And, Addressing: Imm},
	{Opcode: 0x2D, Operation: And, Addressing: Abs},
	{Opcode: 0x31, Operation: And, Addressing: IndY},
	{Opcode: 0x35, Operation: And, Addressing: ZpgX},
	{Opcode: 0x39, Operation: And, Addressing: AbsY},
	{Opcode: 0x3D, Operation: And, Addressing: AbsX},
	/*
		ASL
		addressing   assembler   opc  bytes  cycles
		zeropage     ASL $FF      06    2      5
		accumulator  ASL A        0A    1      2
		absolute     ASL $FFFF    0E    3      6
		zeropage,X   ASL $FF,X    16    2      6
		absolute,X   ASL $FFFF,X  1E    3      7
	*/
	{Opcode: 0x06, Operation: Asl, Addressing: Zpg, CycleAdjust: 1},
	{Opcode: 0x0A, Operation: Asl, Addressing: Acc},
	{Opcode: 0x0E, Operation: Asl, Addressing: Abs, CycleAdjust: 1},
	{Opcode: 0x16, Operation: Asl, Addressing: ZpgX, CycleAdjust: 1},
	{Opcode: 0x1E, Operation: Asl, Addressing: AbsX, CycleAdjust: 2},
	/* TODO: A branch not taken requires two machine cycles. Add one if the branch is taken and add one more if the branch crosses a page boundary.*/
	/*
		BCC
		addressing  assembler  opc  bytes  cycles
		relative    BCC $FF    90   2      2**
	*/
	{Opcode: 0x90, Operation: Bcc, Addressing: Rel},
	/*
		BCS
		addressing  assembler  opc  bytes  cycles
		relative    BCS $FF    B0   2      2**
	*/
	{Opcode: 0xB0, Operation: Bcs, Addressing: Rel},
	/*
		BEQ
		addressing  assembler  opc  bytes  cycles
		relative    BEQ $FF    F0   2      2**
	*/
	{Opcode: 0xF0, Operation: Beq, Addressing: Rel},
	/*
		BIT
		addressing  assembler  opc  bytes  cycles
		zeropage    BIT $FF    24   2      3
		absolute    BIT $FFFF  2C   3      4
	*/
	{Opcode: 0x24, Operation: Bit, Addressing: Zpg},
	{Opcode: 0x2C, Operation: Bit, Addressing: Abs},
	/*
		BMI
		addressing  assembler  opc  bytes  cycles
		relative    BMI $FF    30   2      2**
	*/
	{Opcode: 0x30, Operation: Bmi, Addressing: Rel},
	/*
		BNE
		addressing  assembler  opc  bytes  cycles
		relative    BNE $FF    D0   2      2**
	*/
	{Opcode: 0xD0, Operation: Bne, Addressing: Rel},
	/*
		BPL
		addressing  assembler  opc  bytes  cycles
		relative    BPL $FF    10   2      2**
	*/
	{Opcode: 0x10, Operation: Bpl, Addressing: Rel},
	/*
		BRK
		addressing  assembler  opc  bytes  cycles
		implied      BRK       00   1      7
	*/
	{Opcode: 0x00, Operation: Brk, Addressing: Imp},
	/*
		BVC
		addressing  assembler  opc  bytes  cycles
		relative    BVC $FF    50   2      2**
	*/
	{Opcode: 0x50, Operation: Bvc, Addressing: Rel},
	/*
		BVS
		addressing  assembler  opc  bytes  cycles
		relative    BVS $FF    70   2      2**
	*/
	{Opcode: 0x70, Operation: Bvs, Addressing: Rel},
	/*
		CLC
		addressing  assembler  opc  bytes  cycles
		implied     CLC        18   1      2
	*/
	{Opcode: 0x18, Operation: Clc, Addressing: Imp},
	/*
		CLD
		addressing  assembler  opc  bytes  cycles
		implied     CLD        D8   1      2
	*/
	{Opcode: 0xD8, Operation: Cld, Addressing: Imp},
	/*
		CLI
		addressing  assembler  opc  bytes  cycles
		implied     CLI        58   1      2
	*/
	{Opcode: 0x58, Operation: Cli, Addressing: Imp},
	/*
		CLV
		addressing  assembler  opc  bytes  cycles
		implied     CLV        B8   1      2
	*/
	{Opcode: 0xB8, Operation: Clv, Addressing: Imp},
	/*
		CMP
		addressing  assembler       opc  bytes  cycles
		(indirect,X)  CMP ($FF,X)   C1   2      6
		zeropage      CMP $FF       C5   2      3
		immediate     CMP #$FF      C9   2      2
		absolute      CMP $FFFF     CD   3      4
		(indirect),Y  CMP ($FF),Y   D1   2      5*
		zeropage,X    CMP $FF,X     D5   2      4
		absolute,Y    CMP $FFFF,Y   D9   3      4*
		absolute,X    CMP $FFFF,X   DD   3      4*
	*/
	{Opcode: 0xC1, Operation: Cmp, Addressing: IndX},
	{Opcode: 0xC5, Operation: Cmp, Addressing: Zpg},
	{Opcode: 0xC9, Operation: Cmp, Addressing: Imm},
	{Opcode: 0xCD, Operation: Cmp, Addressing: Abs},
	{Opcode: 0xD1, Operation: Cmp, Addressing: IndY},
	{Opcode: 0xD5, Operation: Cmp, Addressing: ZpgX},
	{Opcode: 0xD9, Operation: Cmp, Addressing: AbsY},
	{Opcode: 0xDD, Operation: Cmp, Addressing: AbsX},
	/*
		CPX
		addressing  assembler  opc  bytes  cycles
		immediate   CPX #$FF   E0   2      2
		zeropage    CPX $FF    E4   2      3
		absolute    CPX $FFFF  EC   3      4
	*/
	{Opcode: 0xE0, Operation: Cpx, Addressing: Imm},
	{Opcode: 0xE4, Operation: Cpx, Addressing: Zpg},
	{Opcode: 0xEC, Operation: Cpx, Addressing: Abs},
	/*
		CPY
		addressing  assembler  opc  bytes  cycles
		immediate   CPY #$FF   C0   2      2
		zeropage    CPY $FF    C4   2      3
		absolute    CPY $FFFF  CC   3      4
	*/
	{Opcode: 0xC0, Operation: Cpy, Addressing: Imm},
	{Opcode: 0xC4, Operation: Cpy, Addressing: Zpg},
	{Opcode: 0xCC, Operation: Cpy, Addressing: Abs},
	/*
		DEC
		addressing  assembler   opc  bytes  cycles
		zeropage    DEC $FF     C6   2      5
		absolute    DEC $FFFF   CE   3      6
		zeropage,X  DEC $FF,X   D6   2      6
		absolute,X  DEC $FFFF,X DE   3      7
	*/
	{Opcode: 0xC6, Operation: Dec, Addressing: Zpg},
	{Opcode: 0xCE, Operation: Dec, Addressing: Abs, CycleAdjust: 1},
	{Opcode: 0xD6, Operation: Dec, Addressing: ZpgX},
	{Opcode: 0xDE, Operation: Dec, Addressing: AbsX, CycleAdjust: 1},
	/*
		DEX
		addressing  assembler  opc  bytes  cycles
		implied     DEX        CA   1      2
	*/
	{Opcode: 0xCA, Operation: Dex, Addressing: Imp},
	/*
		DEY
		addressing  assembler  opc  bytes  cycles
		implied     DEY        88   1      2
	*/
	{Opcode: 0x88, Operation: Dey, Addressing: Imp},
	/*
		EOR
		addressing    assembler     opc  bytes  cycles
		(indirect,X)  EOR ($FF,X)   41   2      6
		zeropage      EOR $FF       45   2      3
		immediate     EOR #$FF      49   2      2
		absolute      EOR $FFFF     4D   3      4
		(indirect),Y  EOR ($FF),Y   51   2      5*
		zeropage,X    EOR $FF,X     55   2      4
		absolute,Y    EOR $FFFF,Y   59   3      4*
		absolute,X    EOR $FFFF,X   5D   3      4*
	*/
	{Opcode: 0x41, Operation: Eor, Addressing: IndX},
	{Opcode: 0x45, Operation: Eor, Addressing: Zpg},
	{Opcode: 0x49, Operation: Eor, Addressing: Imm},
	{Opcode: 0x4D, Operation: Eor, Addressing: Abs},
	{Opcode: 0x51, Operation: Eor, Addressing: IndY},
	{Opcode: 0x55, Operation: Eor, Addressing: ZpgX},
	{Opcode: 0x59, Operation: Eor, Addressing: AbsY},
	{Opcode: 0x5D, Operation: Eor, Addressing: AbsX},
	/*
		INC
		addressing  assembler   opc  bytes  cycles
		zeropage    INC $FF     E6   2      5
		absolute    INC $FFFF   EE   3      6
		zeropage,X  INC $FF,X   F6   2      6
		absolute,X  INC $FFFF,X FE   3      7
	*/
	{Opcode: 0xE6, Operation: Inc, Addressing: Zpg},
	{Opcode: 0xEE, Operation: Inc, Addressing: Abs, CycleAdjust: 1},
	{Opcode: 0xF6, Operation: Inc, Addressing: ZpgX},
	{Opcode: 0xFE, Operation: Inc, Addressing: AbsX, CycleAdjust: 1},
	/*
		INX
		addressing  assembler  opc  bytes  cycles
		implied     INX        E8   1      2
	*/
	{Opcode: 0xE8, Operation: Inx, Addressing: Imp},
	/*
		INY
		addressing  assembler  opc  bytes  cycles
		implied     INY        C8   1      2
	*/
	{Opcode: 0xC8, Operation: Iny, Addressing: Imp},
	/*
		JMP
		addressing  assembler   opc  bytes  cycles
		absolute    JMP $FFFF   4C   3      3
		indirect    JMP ($FF)   6C   3      5
	*/
	{Opcode: 0x4C, Operation: Jmp, Addressing: Abs, CycleAdjust: -1},
	{Opcode: 0x6C, Operation: Jmp, Addressing: Ind},
	/*
		JSR
		addressing  assembler  opc  bytes  cycles
		absolute    JSR $FFFF   20   3      6
	*/
	{Opcode: 0x20, Operation: Jsr, Addressing: Abs},
	/*
		LDA
		addressing    assembler    opc  bytes  cycles
		(indirect,X)  LDA ($FF,X)  A1   2      6
		zeropage      LDA $FF      A5   2      3
		immediate     LDA #$FF     A9   2      2
		absolute      LDA $FFFF    AD   3      4
		(indirect),Y  LDA ($FF),Y  B1   2      5*
		zeropage,X    LDA $FF,X    B5   2      4
		absolute,Y    LDA $FFFF,Y  B9   3      4*
		absolute,X    LDA $FFFF,X  BD   3      4*
	*/
	{Opcode: 0xA1, Operation: Lda, Addressing: IndX},
	{Opcode: 0xA5, Operation: Lda, Addressing: Zpg},
	{Opcode: 0xA9, Operation: Lda, Addressing: Imm},
	{Opcode: 0xAD, Operation: Lda, Addressing: Abs},
	{Opcode: 0xB1, Operation: Lda, Addressing: IndY},
	{Opcode: 0xB5, Operation: Lda, Addressing: ZpgX},
	{Opcode: 0xB9, Operation: Lda, Addressing: AbsY},
	{Opcode: 0xBD, Operation: Lda, Addressing: AbsX},
	/*
		LDX
		addressing  assembler    opc  bytes  cycles
		immediate   LDX #$FF     A2   2      2
		zeropage    LDX $FF      A6   2      3
		absolute    LDX $FFFF    AE   3      4
		zeropage,Y  LDX $FF,Y    B6   2      4
		absolute,Y  LDX $FFFF,Y  BE   3      4*
	*/
	{Opcode: 0xA2, Operation: Ldx, Addressing: Imm},
	{Opcode: 0xA6, Operation: Ldx, Addressing: Zpg},
	{Opcode: 0xAE, Operation: Ldx, Addressing: Abs},
	{Opcode: 0xB6, Operation: Ldx, Addressing: ZpgY},
	{Opcode: 0xBE, Operation: Ldx, Addressing: AbsY},
	/*
		LDY
		addressing  assembler    opc  bytes  cycles
		immediate   LDY #$FF     A0   2      2
		zeropage    LDY $FF      A4   2      3
		absolute    LDY $FFFF    AC   3      4
		zeropage,X  LDY $FF,X    B4   2      4
		absolute,X  LDY $FFFF,X  BC   3      4*
	*/
	{Opcode: 0xA0, Operation: Ldy, Addressing: Imm},
	{Opcode: 0xA4, Operation: Ldy, Addressing: Zpg},
	{Opcode: 0xAC, Operation: Ldy, Addressing: Abs},
	{Opcode: 0xB4, Operation: Ldy, Addressing: ZpgX},
	{Opcode: 0xBC, Operation: Ldy, Addressing: AbsX},
	/*
		LSR
		addressing  assembler    opc  bytes  cycles
		accumulator LSR A        4A   1      2
		zeropage    LSR $FF      46   2      5
		absolute    LSR $FFFF    4E   3      6
		zeropage,X  LSR $FF,X    56   2      6
		absolute,X  LSR $FFFF,X  5E   3      7
	*/
	{Opcode: 0x46, Operation: Lsr, Addressing: Zpg, CycleAdjust: 1},
	{Opcode: 0x4A, Operation: Lsr, Addressing: Acc},
	{Opcode: 0x4E, Operation: Lsr, Addressing: Abs, CycleAdjust: 1},
	{Opcode: 0x56, Operation: Lsr, Addressing: ZpgX, CycleAdjust: 1},
	{Opcode: 0x5E, Operation: Lsr, Addressing: AbsX, CycleAdjust: 2},
	/*
		NOP
		addressing  assembler  opc  bytes  cycles
		implied     NOP        EA   1      2
	*/
	{Opcode: 0xEA, Operation: Nop, Addressing: Imp},
	/*
		ORA
		addressing    assembler     opc  bytes  cycles
		(indirect,X)  ORA ($FF,X)   01   2      6
		zeropage      ORA $FF       05   2      3
		immediate     ORA #$FF      09   2      2
		absolute      ORA $FFFF     0D   3      4
		(indirect),Y  ORA ($FF),Y   11   2      5*
		zeropage,X    ORA $FF,X     15   2      4
		absolute,Y    ORA $FFFF,Y   19   3      4*
		absolute,X    ORA $FFFF,X   1D   3      4*
	*/
	{Opcode: 0x01, Operation: Ora, Addressing: IndX},
	{Opcode: 0x05, Operation: Ora, Addressing: Zpg},
	{Opcode: 0x09, Operation: Ora, Addressing: Imm},
	{Opcode: 0x0D, Operation: Ora, Addressing: Abs},
	{Opcode: 0x11, Operation: Ora, Addressing: IndY},
	{Opcode: 0x15, Operation: Ora, Addressing: ZpgX},
	{Opcode: 0x19, Operation: Ora, Addressing: AbsY},
	{Opcode: 0x1D, Operation: Ora, Addressing: AbsX},
	/*
		PHA
		addressing  assembler  opc  bytes  cycles
		implied     PHA        48   1      3
	*/
	{Opcode: 0x48, Operation: Pha, Addressing: Imp},
	/*
		PHP
		addressing  assembler  opc  bytes  cycles
		implied     PHP        08   1      3
	*/
	{Opcode: 0x08, Operation: Php, Addressing: Imp},
	/*
		PLA
		addressing  assembler  opc  bytes  cycles
		implied     PLA        68   1      4
	*/
	{Opcode: 0x68, Operation: Pla, Addressing: Imp},
	/*
		PLP
		addressing  assembler  opc  bytes  cycles
		implied     PLP        28   1      4
	*/
	{Opcode: 0x28, Operation: Plp, Addressing: Imp},
	/*
		ROL
		addressing   assembler   opc  bytes  cycles
		zeropage     ROL $FF     26   2      5
		accumulator  ROL A       2A   1      2
		absolute     ROL $FFFF   2E   3      6
		zeropage,X   ROL $FF,X   36   2      6
		absolute,X   ROL $FFFF,X 3E   3      7
	*/
	{Opcode: 0x26, Operation: Rol, Addressing: Zpg, CycleAdjust: 1},
	{Opcode: 0x2A, Operation: Rol, Addressing: Acc},
	{Opcode: 0x2E, Operation: Rol, Addressing: Abs, CycleAdjust: 1},
	{Opcode: 0x36, Operation: Rol, Addressing: ZpgX, CycleAdjust: 1},
	{Opcode: 0x3E, Operation: Rol, Addressing: AbsX, CycleAdjust: 2},
	/*
		ROR
		addressing  assembler    opc  bytes  cycles
		zeropage	ROR $FF      66  2      5
		accumulator	ROR A        6A  1      2
		absolute	ROR $FFFF    6E  3      6
		zeropage,X	ROR $FF,X    76  2      6
		absolute,X	ROR $FFFF,X  7E  3      7
	*/
	{Opcode: 0x66, Operation: Ror, Addressing: Zpg, CycleAdjust: 1},
	{Opcode: 0x6A, Operation: Ror, Addressing: Acc},
	{Opcode: 0x6E, Operation: Ror, Addressing: Abs, CycleAdjust: 1},
	{Opcode: 0x76, Operation: Ror, Addressing: ZpgX, CycleAdjust: 1},
	{Opcode: 0x7E, Operation: Ror, Addressing: AbsX, CycleAdjust: 2},
	/*
		RTI
		addressing  assembler  opc  bytes  cycles
		implied     RTI        40   1      6
	*/
	{Opcode: 0x40, Operation: Rti, Addressing: Imp},
	/*
		RTS
		addressing  assembler  opc  bytes  cycles
		implied     RTS        60   1      6
	*/
	{Opcode: 0x60, Operation: Rts, Addressing: Imp},
	/*
		SBC
		addressing    assembler    opc  bytes  cycles
		(indirect,X)  SBC ($FF,X)  E1   2      6
		zeropage      SBC $FF      E5   2      3
		immediate     SBC #$FF     E9   2      2
		absolute      SBC $FFFF    ED   3      4
		(indirect),Y  SBC ($FF),Y  F1   2      5*
		zeropage,X    SBC $FF,X    F5   2      4
		absolute,Y    SBC $FFFF,Y  F9   3      4*
		absolute,X    SBC $FFFF,X  FD   3      4*
	*/
	{Opcode: 0xE1, Operation: Sbc, Addressing: IndX},
	{Opcode: 0xE5, Operation: Sbc, Addressing: Zpg},
	{Opcode: 0xE9, Operation: Sbc, Addressing: Imm},
	{Opcode: 0xED, Operation: Sbc, Addressing: Abs},
	{Opcode: 0xF1, Operation: Sbc, Addressing: IndY},
	{Opcode: 0xF5, Operation: Sbc, Addressing: ZpgX},
	{Opcode: 0xF9, Operation: Sbc, Addressing: AbsY},
	{Opcode: 0xFD, Operation: Sbc, Addressing: AbsX},
	/*
		SEC
		addressing  assembler  opc  bytes  cycles
		implied     SEC        38   1      2
	*/
	{Opcode: 0x38, Operation: Sec, Addressing: Imp},
	/*
		SED
		addressing  assembler  opc  bytes  cycles
		implied     SED        F8   1      2
	*/
	{Opcode: 0xF8, Operation: Sed, Addressing: Imp},
	/*
		SEI
		addressing  assembler  opc  bytes  cycles
		implied     SEI        78   1      2
	*/
	{Opcode: 0x78, Operation: Sei, Addressing: Imp},
	/*
		STA
		addressing    assembler    opc  bytes  cycles
		(indirect,X)  STA ($FF,X)  81   2      6
		zeropage      STA $FF      85   2      3
		absolute      STA $FFFF    8D   3      4
		(indirect),Y  STA ($FF),Y  91   2      6
		zeropage,X    STA $FF,X    95   2      4
		absolute,Y    STA $FFFF,Y  99   3      5
		absolute,X    STA $FFFF,X  9D   3      5
	*/
	{Opcode: 0x81, Operation: Sta, Addressing: IndX},
	{Opcode: 0x85, Operation: Sta, Addressing: Zpg},
	{Opcode: 0x8D, Operation: Sta, Addressing: Abs},
	{Opcode: 0x91, Operation: Sta, Addressing: IndY, CycleAdjust: 1},
	{Opcode: 0x95, Operation: Sta, Addressing: ZpgX},
	{Opcode: 0x99, Operation: Sta, Addressing: AbsY, CycleAdjust: 1},
	{Opcode: 0x9D, Operation: Sta, Addressing: AbsX, CycleAdjust: 1},
	/*
		STX
		addressing  assembler  opc  bytes  cycles
		zeropage	STX $FF    86   2      3
		absolute	STX $FFFF  8E   3      4
		zeropage,Y	STX $FF,Y  96   2      4
	*/
	{Opcode: 0x86, Operation: Stx, Addressing: Zpg},
	{Opcode: 0x8E, Operation: Stx, Addressing: Abs},
	{Opcode: 0x96, Operation: Stx, Addressing: ZpgY},
	/*
		STY
		addressing  assembler  opc  bytes  cycles
		zeropage    STY $FF    84   2      3
		absolute    STY $FFFF  8C   3      4
		zeropage,X  STY $FF,X  94   2      4
	*/
	{Opcode: 0x84, Operation: Sty, Addressing: Zpg},
	{Opcode: 0x8C, Operation: Sty, Addressing: Abs},
	{Opcode: 0x94, Operation: Sty, Addressing: ZpgX},
	/*
		TAX
		addressing  assembler  opc  bytes  cycles
		implied     TAX        AA   1      2
	*/
	{Opcode: 0xAA, Operation: Tax, Addressing: Imp},
	/*
		TAY
		addressing  assembler  opc  bytes  cycles
		implied     TAY        A8   1      2

	*/
	{Opcode: 0xA8, Operation: Tay, Addressing: Imp},
	/*
		TSX
		addressing  assembler  opc  bytes  cycles
		implied     TSX        BA   1      2
	*/
	{Opcode: 0xBA, Operation: Tsx, Addressing: Imp},
	/*
		TXA
		addressing  assembler  opc  bytes  cycles
		implied     TXA        8A   1      2
	*/
	{Opcode: 0x8A, Operation: Txa, Addressing: Imp},
	/*
		TXS
		addressing  assembler  opc  bytes  cycles
		implied     TXS        9A   1      2
	*/
	{Opcode: 0x9A, Operation: Txs, Addressing: Imp},
	/*
		TYA
		addressing  assembler  opc  bytes  cycles
		implied     TYA        98   1      2
	*/
	{Opcode: 0x98, Operation: Tya, Addressing: Imp},
}
