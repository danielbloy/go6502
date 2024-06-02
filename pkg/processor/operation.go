package processor

// ************************************************************
// ********** Logic operation functions
// ************************************************************

// Operation performs the actual logic operation of the instruction using the
// starting State and the State returned by addressing. The new CPU State is
// returned. The PC of the input State will be pointing to the next instruction
// to execute.
type Operation func(State, Addressing) (State, error)

// AddWithCarry (ADC). ADC results are dependent on the setting of the decimal flag.
// In decimal mode, addition is carried out on the assumption that the values involved
// are packed BCD (Binary Coded Decimal). There is no way to add without carry.
//
// Implementation borrowed from testADC() in https://github.com/skilldrick/easy6502/blob/gh-pages/simulator/assembler.js
//
//	See: APPENDIX A: WHAT ABOUT INVALID BCD VALUES AND INVALID FLAGS?
//	At: http://www.6502.org/tutorials/decimal_mode.html#A
func AddWithCarry(state State, addressing Addressing) (State, error) {
	accum := uint16(state.A)
	value := uint16(addressing.Value)
	carry := uint16(0)
	result := uint16(0)

	if state.P.ToFlags().Carry {
		carry = 1
	}

	if state.P.ToFlags().Decimal {
		// Do the lower nibble first.
		result = (accum & 0x0F) + (value & 0x0F) + carry

		if result >= 10 {
			result = 0x10 | ((result + 0x06) & 0x0F)
		}

		// Do the upper nibble.
		result += (accum & 0xF0) + (value & 0xF0)

		// Adjust for carry in decimal mode.
		if result >= 0xA0 {
			result += 0x60
		}

	} else {
		result = accum + value + carry
	}

	overflowSet(&state.P, accum, value, result)
	carrySet(&state.P, result)
	negativeSet(&state.P, result)
	zeroSet(&state.P, result)

	state.A = uint8(result)
	return state, nil
}

// AndWithA (AND). Bitwise AND memory with accumulator register A.
func AndWithA(state State, addressing Addressing) (State, error) {

	state.A = addressing.Value & state.A

	negativeSet(&state.P, uint16(state.A))
	zeroSet(&state.P, uint16(state.A))

	return state, nil
}

// ArithmeticShiftLeft (ASL).Shift Left One Bit (Memory or Accumulator).
func ArithmeticShiftLeft(state State, addressing Addressing) (State, error) {

	value := uint16(addressing.Value) << 1

	carrySet(&state.P, value)
	negativeSet(&state.P, value)
	zeroSet(&state.P, value)

	return addressing.Store(state, uint8(value&0x00FF))
}

// TODO: A branch not taken requires two machine cycles. Add one if the branch is taken and add one more if the branch crosses a page boundary.

// BranchOnCarryClear (BCC). Branch on Carry flag not set.
func BranchOnCarryClear(state State, addressing Addressing) (State, error) {

	if !state.P.ToFlags().Carry {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// BranchOnCarrySet (BCS). Branch on Carry flag set.
func BranchOnCarrySet(state State, addressing Addressing) (State, error) {

	if state.P.ToFlags().Carry {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// BranchOnEqual (BEQ). Branch on Zero flag not set.
func BranchOnEqual(state State, addressing Addressing) (State, error) {

	if state.P.ToFlags().Zero {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// BranchOnMinus (BMI). Branch on Negative flag set.
func BranchOnMinus(state State, addressing Addressing) (State, error) {

	if state.P.ToFlags().Negative {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// BranchOnNotEqual (BNE). Branch on Zero flag not set.
func BranchOnNotEqual(state State, addressing Addressing) (State, error) {

	if !state.P.ToFlags().Zero {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// BranchOnOverflowClear (BVC). Branch on Carry flag not set.
func BranchOnOverflowClear(state State, addressing Addressing) (State, error) {

	if !state.P.ToFlags().Overflow {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// BranchOnOverflowSet (BVS). Branch on Overflow flag set.
func BranchOnOverflowSet(state State, addressing Addressing) (State, error) {

	if state.P.ToFlags().Overflow {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// BranchOnPlus (BPL). Branch on Plus; Negative flag not set.
func BranchOnPlus(state State, addressing Addressing) (State, error) {

	if !state.P.ToFlags().Negative {
		state.PC = addressing.EffectiveAddress
	}

	return state, nil
}

// Break (BRK). Break initiates a non-maskable software interrupt similar to a
// hardware interrupt. Break increments the program counter by one before pushing
// the return address to the stack. Therefore, an RTI will go to the address of the
// Break instruction +2 so that Break may be used to replace a two-byte instruction
// for debugging and the subsequent RTI will be correct. The return address is
// pushed high byte first, low byte second.
//
// After the return address is pushed, the status register will be pushed to the
// stack with the break flag set to 1. However, when retrieved during RTI or by
// a PLP instruction, the break flag will be ignored. The interrupt disable flag
// is not set automatically.
//
// See also Nmi() and Interrupt()  which are very similar.
func Break(state State, addressing Addressing) (State, error) {

	// Get the interrupt vector from memory. We do this first to avoid
	// the vector being overwritten in tests that use a tiny memory and
	// the vector overlaps with the stack.
	irqVector, err := ReadIrqVectorFromMemory(addressing.Memory)
	if err != nil {
		return state, err
	}

	// Advance program counter by 1 and push
	addressToPush := state.PC + 1
	state, err = addressing.PushAddress(state, addressToPush)
	if err != nil {
		return state, err
	}

	state, err = addressing.PushStatus(state, state.P.WithBreakSet())
	if err != nil {
		return state, err
	}

	// Set the interrupt flag and new program counter.
	state.P.SetInterrupt()
	state.PC = irqVector

	return state, nil
}

// ClearCarry (CLC). Clear carry flag.
func ClearCarry(state State, _ Addressing) (State, error) {
	state.P.ClearCarry()
	return state, nil
}

// ClearDecimal (CLD). Clear decimal flag.
func ClearDecimal(state State, _ Addressing) (State, error) {
	state.P.ClearDecimal()
	return state, nil
}

// ClearInterrupt (CLI). Clear interrupt flag.
func ClearInterrupt(state State, _ Addressing) (State, error) {
	state.P.ClearInterrupt()
	return state, nil
}

// ClearOverflow (CLV). Clear overflow flag.
func ClearOverflow(state State, _ Addressing) (State, error) {
	state.P.ClearOverflow()
	return state, nil
}

// Performs the actual compare for the three CompareWith... functions; setting or clearing the
// flags as appropriate. If register is equal or greater than value, the Carry will be set. The
// equal (Z) flags will be set based on equality or lack thereof and the negative (N) flag will
// be set on the sign of the result (i.e. result >= $80) of the accumulator.
func compare(s *Status, register, value uint8) {

	if register >= value {
		s.SetCarry()
	} else {
		s.ClearCarry()
	}

	result := uint16(register) - uint16(value)

	negativeSet(s, result)
	zeroSet(s, result)
}

// CompareWithA (CMP). Compare memory with accumulator. Compare sets flags as if a subtraction
// had been carried out. If the value in the accumulator is equal or greater than the compared
// value, the Carry will be set. The equal (Z) and negative (N) flags will be set based on
// equality or lack thereof and the sign (i.e. A>=$80) of the accumulator.
func CompareWithA(state State, addressing Addressing) (State, error) {
	compare(&state.P, state.A, addressing.Value)
	return state, nil
}

// CompareWithX (CPX). Compare memory with X. Same rules as CompareWithAccumulator but based on X.
func CompareWithX(state State, addressing Addressing) (State, error) {
	compare(&state.P, state.X, addressing.Value)
	return state, nil
}

// CompareWithY (CPY). Compare memory with Y. Same rules as CompareWithAccumulator but based on Y.
func CompareWithY(state State, addressing Addressing) (State, error) {
	compare(&state.P, state.Y, addressing.Value)
	return state, nil
}

// Decrement (DEC). Decrement Memory by One.
func Decrement(state State, addressing Addressing) (State, error) {

	value := addressing.Value - 1

	negativeSet(&state.P, uint16(value))
	zeroSet(&state.P, uint16(value))

	return addressing.Store(state, value)
}

// DecrementX (DEX). Decrement Index X by One.
func DecrementX(state State, _ Addressing) (State, error) {

	state.X--

	negativeSet(&state.P, uint16(state.X))
	zeroSet(&state.P, uint16(state.X))

	return state, nil
}

// DecrementY (DEY). Decrement Index Y by One.
func DecrementY(state State, _ Addressing) (State, error) {

	state.Y--

	negativeSet(&state.P, uint16(state.Y))
	zeroSet(&state.P, uint16(state.Y))

	return state, nil
}

// ExclusiveOrWithA (EOR). Exclusive or memory with accumulator register A.
func ExclusiveOrWithA(state State, addressing Addressing) (State, error) {

	state.A = addressing.Value ^ state.A

	negativeSet(&state.P, uint16(state.A))
	zeroSet(&state.P, uint16(state.A))

	return state, nil
}

// Increment (INC). Increment memory by One.
func Increment(state State, addressing Addressing) (State, error) {

	value := addressing.Value + 1

	negativeSet(&state.P, uint16(value))
	zeroSet(&state.P, uint16(value))

	return addressing.Store(state, value)
}

// IncrementX (INX). Increment index X by One.
func IncrementX(state State, _ Addressing) (State, error) {

	state.X++

	negativeSet(&state.P, uint16(state.X))
	zeroSet(&state.P, uint16(state.X))

	return state, nil
}

// IncrementY (INY). Increment index Y by One.
func IncrementY(state State, _ Addressing) (State, error) {

	state.Y++

	negativeSet(&state.P, uint16(state.Y))
	zeroSet(&state.P, uint16(state.Y))

	return state, nil
}

// Interrupt performs a hardware interrupt. This is similar in operation to break
// except it does not advance the program counter before pushing and does not set
// the break flag before pushing the status register to the stack. See also Nmi().
func Interrupt(state State, addressing Addressing) (State, error) {

	// Get the interrupt vector from memory. We do this first to avoid
	// the vector being overwritten in tests that use a tiny memory and
	// the vector overlaps with the stack.
	irqVector, err := ReadIrqVectorFromMemory(addressing.Memory)
	if err != nil {
		return state, err
	}

	state, err = addressing.PushAddress(state, state.PC)
	if err != nil {
		return state, err
	}

	state, err = addressing.PushStatus(state, state.P)
	if err != nil {
		return state, err
	}

	// Set the interrupt flag and new program counter.
	state.P.SetInterrupt()
	state.PC = irqVector

	return state, nil
}

// Jump (JMP). Jump to new Location.
func Jump(state State, addressing Addressing) (State, error) {

	state.PC = addressing.EffectiveAddress
	return state, nil
}

// JumpSubRoutine (JSR). Jump to new location saving return address onto the stack. The
// return address that is pushed is the program counter - 1.
func JumpSubRoutine(state State, addressing Addressing) (State, error) {

	addressToPush := state.PC - 1
	state.PC = addressing.EffectiveAddress
	return addressing.PushAddress(state, addressToPush)
}

// LoadA (LDA). Load accumulator with memory.
func LoadA(state State, addressing Addressing) (State, error) {

	state.A = addressing.Value

	negativeSet(&state.P, uint16(state.A))
	zeroSet(&state.P, uint16(state.A))

	return state, nil
}

// LoadX (LDX). Load index X with memory.
func LoadX(state State, addressing Addressing) (State, error) {

	state.X = addressing.Value

	negativeSet(&state.P, uint16(state.X))
	zeroSet(&state.P, uint16(state.X))

	return state, nil
}

// LoadY (LDY). Load index Y with memory.
func LoadY(state State, addressing Addressing) (State, error) {

	state.Y = addressing.Value

	negativeSet(&state.P, uint16(state.Y))
	zeroSet(&state.P, uint16(state.Y))

	return state, nil
}

// LogicalShiftRight (LSR). Shift one bit right (memory or accumulator). Bit 7 is set to zero and bit
// zero is shifted into the carry flag.
func LogicalShiftRight(state State, addressing Addressing) (State, error) {

	value := uint16(addressing.Value) >> 1

	// Set carry by switching on a bit in the higher byte.
	if addressing.Value&0x01 != 0 {
		value |= 0x0100
	}

	carrySet(&state.P, value)
	negativeSet(&state.P, value)
	zeroSet(&state.P, value)

	return addressing.Store(state, uint8(value&0x00FF))
}

// Nmi performs a non-maskable interrupt. This is very similar to Break and
// Interrupt but uses a different vector.
func Nmi(state State, addressing Addressing) (State, error) {

	// Get the nmi vector from memory. We do this first to avoid
	// the vector being overwritten in tests that use a tiny memory and
	// the vector overlaps with the stack.
	nmiVector, err := ReadNmiVectorFromMemory(addressing.Memory)
	if err != nil {
		return state, err
	}

	state, err = addressing.PushAddress(state, state.PC)
	if err != nil {
		return state, err
	}

	state, err = addressing.PushStatus(state, state.P)
	if err != nil {
		return state, err
	}

	// Set the interrupt flag and new program counter.
	state.P.SetInterrupt()
	state.PC = nmiVector

	return state, nil
}

// NoOperation simply returns the State passed in and does not access the
// memory nor use the value.
func NoOperation(state State, _ Addressing) (State, error) {
	return state, nil
}

// OrWithA (ORA) performs a bitwise OR with Accumulator.
func OrWithA(state State, addressing Addressing) (State, error) {

	state.A = addressing.Value | state.A

	negativeSet(&state.P, uint16(state.A))
	zeroSet(&state.P, uint16(state.A))

	return state, nil
}

// PushA (PHA) pushes the accumulator onto the stack.
func PushA(state State, addressing Addressing) (State, error) {

	state, err := addressing.PushByte(state, state.A)
	if err != nil {
		return state, err
	}
	return state, nil
}

// PushP (PHP) pushes the processor status register onto the stack. The status register will
// be pushed with the break flag (and constant flag) set.
func PushP(state State, addressing Addressing) (State, error) {

	state, err := addressing.PushStatus(state, state.P.WithBreakSet())
	if err != nil {
		return state, err
	}
	return state, nil
}

// PullA (PLA) pulls the accumulator from the stack. This will set the sign and zero
// flags based on the result pulled.
func PullA(state State, addressing Addressing) (State, error) {

	state, value, err := addressing.PullByte(state)
	if err != nil {
		return state, err
	}
	state.A = value

	negativeSet(&state.P, uint16(value))
	zeroSet(&state.P, uint16(value))

	return state, nil
}

// PullP (PLP) pulls the processor status from the stack. The status register will be
// pulled with the break flag and constant flag ignored. See these references:
// https://www.masswerk.at/6502/6502_instruction_set.html#:~:text=Since%20there%20is%20no%20actual,there%20is%20no%20internal%20representation.
// The break flag (B) is not an actual flag implemented in a register, and rather
//
//	appears only, when the status register is pushed onto or pulled from the stack.
//
// ...
// Since there is no actual slot for the break flag, it will be always ignored, when
// retrieved (PLP or RTI). The break flag is not accessed by the CPU at anytime and
// there is no internal representation.
func PullP(state State, addressing Addressing) (State, error) {
	state, err := addressing.PullStatus(state)
	if err != nil {
		return state, err
	}
	return state, nil
}

// RotateLeft (ROL) shifts all bits left one position. The Carry is
// shifted into bit 0 and the original bit 7 is shifted into the Carry.
func RotateLeft(state State, addressing Addressing) (State, error) {

	value := uint16(addressing.Value) << 1
	if state.P.ToFlags().Carry {
		value++
	}

	carrySet(&state.P, value)
	negativeSet(&state.P, value)
	zeroSet(&state.P, value)

	return addressing.Store(state, uint8(value&0x00FF))
}

// RotateRight (ROR) shifts all bits right one position. The Carry is
// shifted into bit 7 and the original bit 0 is shifted into the Carry.
func RotateRight(state State, addressing Addressing) (State, error) {

	value := uint16(addressing.Value)

	// Apply the existing carry bit ready for shifting
	if state.P.ToFlags().Carry {
		value = value | 0x100
	}

	// Roll bit zero around ready to be the next carry
	if value&0x01 != 0 {
		value = value | 0x200
	}

	value = value >> 1

	carrySet(&state.P, value)
	negativeSet(&state.P, value)
	zeroSet(&state.P, value)

	return addressing.Store(state, uint8(value&0x00FF))
}

// ReturnFromInterrupt (RTI) retrieves the Processor Status Word (flags)
// and the Program Counter from the stack in that order (interrupts push
// the PC first and then P). Note that unlike RTS, the return address on
// the stack is the actual address rather than the address - 1.
func ReturnFromInterrupt(state State, addressing Addressing) (State, error) {

	state, err := addressing.PullStatus(state)
	if err != nil {
		return state, err
	}

	state, addressFromStack, err := addressing.PullAddress(state)
	if err != nil {
		return state, err
	}

	state.PC = addressFromStack

	return state, nil
}

// ReturnFromSubroutine (RTS) pulls the top two bytes off the stack (low
// byte first) and transfers program control to that address + 1. It is
// used, as expected, to exit a subroutine invoked via JSR which pushed
// the address - 1.
func ReturnFromSubroutine(state State, addressing Addressing) (State, error) {

	state, addressFromStack, err := addressing.PullAddress(state)
	if err != nil {
		return state, err
	}

	state.PC = addressFromStack + 1

	return state, nil
}

// SubtractWithCarry (SBC). This subtracts a value from the accumulator and also subtracts the borrow
// bit. There is no explicit borrow flag, instead the complement of the carry flag is used. If the:
//   - Carry flag is 1, then borrow is 0
//   - Carry flag is 0, then borrow is 1
//
// If the (unsigned) operation results in a borrow (i.e. the result is negative), then the borrow bit
// is set (i.e. the Carry flag is cleared). For more details, see:
//
//	https://www.righto.com/2012/12/the-6502-overflow-flag-explained.html#:~:text=The%206502%20has%20a%20SBC,the%20carry%20flag%20is%20used.
func SubtractWithCarry(state State, addressing Addressing) (State, error) {

	accum := uint16(state.A)
	value := uint16(addressing.Value) ^ 0x00FF
	carry := uint16(0)
	result := uint16(0)

	if state.P.ToFlags().Carry {
		carry = 1
	}

	if state.P.ToFlags().Decimal {

		// Do the lower nibble
		result = 0x0F + (accum & 0x0F) - (uint16(addressing.Value) & 0x0F) + carry
		upperNibble := uint16(0)
		if result < 0x10 {
			result -= 0x06
		} else {
			result -= 0x10
			upperNibble = 0x10
		}

		// Do the upper nibble
		upperNibble += 0xF0 + (accum & 0xF0) - (uint16(addressing.Value) & 0xF0)

		if upperNibble < 0x100 {
			upperNibble -= 0x60
		}
		result += upperNibble

	} else {
		result = accum + value + carry
	}

	overflowSet(&state.P, accum, value, result)
	carrySet(&state.P, result)
	negativeSet(&state.P, result)
	zeroSet(&state.P, result)

	state.A = uint8(result)
	return state, nil
}

// SetCarry (SEC). Sets the carry flag.
func SetCarry(state State, _ Addressing) (State, error) {
	state.P.SetCarry()
	return state, nil
}

// SetDecimal (SED). Sets the decimal flag.
func SetDecimal(state State, _ Addressing) (State, error) {
	state.P.SetDecimal()
	return state, nil
}

// SetInterrupt (SEI). Sets the interrupt flag to present maskable interrupts (aka IRQs).
func SetInterrupt(state State, _ Addressing) (State, error) {
	state.P.SetInterrupt()
	return state, nil
}

// StoreA (STA). Store accumulator to memory.
func StoreA(state State, addressing Addressing) (State, error) {
	return addressing.Store(state, state.A)
}

// StoreX (STX). Store index X to memory.
func StoreX(state State, addressing Addressing) (State, error) {
	return addressing.Store(state, state.X)
}

// StoreY (STY). Store index Y to memory.
func StoreY(state State, addressing Addressing) (State, error) {
	return addressing.Store(state, state.Y)
}

// TestBitsInMemoryWithAccumulator (BIT). Sets the Zero (Z) flag as though the value in the address tested were ANDed
// with the accumulator (but does not change the accumulator). Bits 7 and 6 of the value from memory are copied into
// the Sign (N) flag (bit 7) and Overflow (V) flags (bit 6).
func TestBitsInMemoryWithAccumulator(state State, addressing Addressing) (State, error) {

	// Do zero and negative flags first.
	value := uint16(addressing.Value & state.A)
	zeroSet(&state.P, value)
	negativeSet(&state.P, uint16(addressing.Value))

	state.P.ClearOverflow()
	if addressing.Value&FlagOverflow != 0 {
		state.P.SetOverflow()
	}

	return state, nil
}

// TransferAtoX (TAX). Transfer accumulator to index X.
func TransferAtoX(state State, _ Addressing) (State, error) {

	state.X = state.A

	value := uint16(state.X)
	zeroSet(&state.P, value)
	negativeSet(&state.P, value)

	return state, nil
}

// TransferAtoY (TAY). Transfer accumulator to index Y.
func TransferAtoY(state State, _ Addressing) (State, error) {

	state.Y = state.A

	value := uint16(state.Y)
	zeroSet(&state.P, value)
	negativeSet(&state.P, value)

	return state, nil
}

// TransferSPtoX (TSX). Transfer stack pointer to index X.
func TransferSPtoX(state State, _ Addressing) (State, error) {

	state.X = state.SP

	value := uint16(state.X)
	zeroSet(&state.P, value)
	negativeSet(&state.P, value)

	return state, nil
}

// TransferXtoA (TXA). Transfer index X to accumulator.
func TransferXtoA(state State, _ Addressing) (State, error) {

	state.A = state.X

	value := uint16(state.A)
	zeroSet(&state.P, value)
	negativeSet(&state.P, value)

	return state, nil
}

// TransferXtoSP (TXS). Transfer index X to stack pointer.
func TransferXtoSP(state State, _ Addressing) (State, error) {

	state.SP = state.X

	return state, nil
}

// TransferYtoA (TYA). Transfer index Y to accumulator.
func TransferYtoA(state State, _ Addressing) (State, error) {

	state.A = state.Y

	value := uint16(state.A)
	zeroSet(&state.P, value)
	negativeSet(&state.P, value)

	return state, nil
}

// ************************************************************
// ********** Utility helper functions
// ************************************************************

// These functions help with the implementation of some of the arithmetic and operation.
// They take unsigned 16-bit values as input, representing the 9th bit as the carry flag.

// isCarrySet returns whether any bit in the higher byte is set.
func isCarrySet(value uint16) bool {
	return (value & 0xFF00) != 0
}

// carrySet sets or clears the carry flag if the upper byte of value has
// any bit set.
func carrySet(s *Status, value uint16) {
	if isCarrySet(value) {
		s.SetCarry()
	} else {
		s.ClearCarry()
	}
}

// isNegative returns whether bit 8 is set or not.
func isNegative(value uint16) bool {
	return (value & 0x80) != 0
}

// negativeSet sets or clears the negative flag based on whether the
// sign bit (bit 8) is set or not.
func negativeSet(s *Status, value uint16) {
	if isNegative(value) {
		s.SetNegative()
	} else {
		s.ClearNegative()
	}
}

// isZero returns whether lower byte is zero or not.
func isZero(value uint16) bool {
	return (value & 0x00FF) == 0
}

// zeroSet sets or clears the zero flag based on whether the
// lower byte is zero or not.
func zeroSet(s *Status, value uint16) {
	if isZero(value) {
		s.SetZero()
	} else {
		s.ClearZero()
	}
}

// isOverflow determines if an overflow occurred when value was added to start
// (with carry) which resulted in result. This supports signed addition too.
//
// From https://www.nesdev.org/wiki/Status_flags:
//
//	ADC and SBC will set this flag if the signed result would be invalid[1],
func isOverflow(decimal bool, accumulator, value, result uint16) bool {
	overflow := (accumulator^value)&0x80 == 0

	// Pick the correct comparator based on Decimal or Binary mode.
	comparator := uint16(0x100)
	if decimal {
		comparator = 0xA0
	}

	if result >= comparator {
		if overflow && result >= 0x180 {
			overflow = false
		}
	} else {
		if overflow && result < 0x80 {
			overflow = false
		}
	}

	return overflow
}

// overflowSet determines with an overflow occurred when value was added to
// start (with carry) which resulted in result.
func overflowSet(s *Status, accumulator, value, result uint16) {
	if isOverflow(s.ToFlags().Decimal, accumulator, value, result) {
		s.SetOverflow()
	} else {
		s.ClearOverflow()
	}
}
