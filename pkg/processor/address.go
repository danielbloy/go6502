package processor

import "fmt"

const BaseStack = Address(0x100)

type Address uint16

// MakeAddress converts two 8-bit values into a 16-bit address. The first byte in memory
// is typically the low byte and the second address in memory is the high byte. This is
// the reverse of SplitAddress.
func MakeAddress(l, h uint8) Address {
	return Address(l) | (Address(h) << 8)
}

// SplitAddress converts a 16-bit address into it's low and high 8-bit values.
// The first byte returned is the low byte and the second byte returned is the
// high byte. This is the reverse of MakeAddress.
func SplitAddress(address Address) (uint8, uint8) {
	low := uint8(address & 0xFF)
	high := uint8((address >> 8) & 0xFF)
	return low, high
}

// Addressing instances are generated as a result of executing an AddressingFunc function.
// An Addressing instance
type Addressing struct {
	// How much to change the program counter by as part of addressing. This is essentially
	// how many bytes are read from memory as part of addressing.
	ProgramCounterChange Address

	// The value retrieved by this Addressing mode. It may have come from Memory or the Accumulator.
	// if the Accumulator member is true.
	Value uint8

	// Is this addressing mode using the Accumulator or Memory for writing results.
	Accumulator bool

	// The calculated effective address. If a Value is also returned, this will be the address
	// from Memory that the Value was retrieved from.
	EffectiveAddress Address

	// Does the EffectiveAddress cross a page boundary.
	PageBoundaryCrossed bool

	// The Memory that was used when generating Addressing and where results "may"
	// be written when using Store().
	Memory Memory
}

// Store Value to either the Effective Address location in Memory or the Accumulator based
// on the State of the Addressing instance.
func (as Addressing) Store(state State, value uint8) (State, error) {
	if as.Accumulator {
		state.A = value
	} else {
		if as.Memory == nil {
			return state, MemoryMustBeProvided
		}
		as.Memory.Write(as.EffectiveAddress, value)
	}

	return state, nil
}

// PushByte saves the 8-bit value onto the stack (reducing it by 1 as required) and
// returning the new State.
func (as Addressing) PushByte(state State, byte uint8) (State, error) {

	if as.Memory == nil {
		return state, MemoryMustBeProvided
	}

	address := BaseStack + Address(state.SP)
	state.SP--
	as.Memory.Write(address, byte)

	return state, nil
}

// PushStatus saves the 8-bit status value onto the stack via PushByte(). This forces the constant flag
// to be set.
func (as Addressing) PushStatus(state State, status Status) (State, error) {
	return as.PushByte(state, uint8(status.WithConstantSet()))
}

// PushAddress saves the 16-bit address onto the stack, returning the new State. The high byte
// is pushed first (to SP) and the low byte pushed second (to SP - 1). This is so the low byte
// is fetched first when the stack is popped.
func (as Addressing) PushAddress(state State, address Address) (State, error) {

	lowByte, highByte := SplitAddress(address)

	state, err := as.PushByte(state, highByte)
	if err != nil {
		return state, err
	}

	return as.PushByte(state, lowByte)
}

// PullByte retrieves the 8-bit value from the stack (increasing it by 1 as required) and
// returning the new State. The pulled value is returned.
func (as Addressing) PullByte(state State) (State, uint8, error) {

	if as.Memory == nil {
		return state, 0, MemoryMustBeProvided
	}

	state.SP++
	address := BaseStack + Address(state.SP)
	value := as.Memory.Read(address)

	return state, value, nil
}

// PullStatus pulls the 8-bit status value from the stack via PullByte() and places the
// result in State. This forces the constant flag to be set and the break flag to be
// cleared.
func (as Addressing) PullStatus(state State) (State, error) {

	state, value, err := as.PullByte(state)
	if err != nil {
		return state, err
	}

	status := Status(value)
	status.ClearBreak().SetConstant()
	state.P = status
	return state, nil
}

// PullAddress pulls the 16-bit address from the stack, returning the new State and the
// address. The low byte is pull first (from SP - 1) and the high byte pulled second
// (to SP - 2).
func (as Addressing) PullAddress(state State) (State, Address, error) {

	state, lowByte, err := as.PullByte(state)
	if err != nil {
		return state, Address(0), err
	}

	state, highByte, err := as.PullByte(state)
	if err != nil {
		return state, Address(0), err
	}

	address := MakeAddress(lowByte, highByte)

	return state, address, nil
}

// Converts the Addressing instance into a canonical string form.
func (as Addressing) String() string {
	return fmt.Sprintf(
		"Acc: %v, EA: %04X, V: %02X, PC-delta: %04X, PBC: %v",
		as.Accumulator, as.EffectiveAddress, as.Value, as.ProgramCounterChange, as.PageBoundaryCrossed)
}

// ************************************************************
// ********** AddressingFunc functions
// ************************************************************

// AddressingFunc performs the addressing mode phase of an instructions' execution.
// AddressingFunc is always done before Operation as it will calculate the
// effective address (if relevant) and return the value from that address (if relevant).
type AddressingFunc func(State, Memory) (Addressing, error)

func AbsoluteAddressing(state State, memory Memory, offset Address) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	low := memory.Read(state.PC)
	high := memory.Read(state.PC + 1)
	effectiveAddress := MakeAddress(low, high)

	pageBoundaryCrossed := false

	// Determine if a page boundary has been crossed by the offset
	if offset != 0 {
		startPage := effectiveAddress & 0xFF00
		effectiveAddress += offset
		pageBoundaryCrossed = startPage != (effectiveAddress & 0xFF00)
	}

	result := Addressing{
		EffectiveAddress:     effectiveAddress,
		Value:                memory.Read(effectiveAddress),
		ProgramCounterChange: 2,
		PageBoundaryCrossed:  pageBoundaryCrossed,
		Memory:               memory,
	}

	return result, nil
}

// Absolute addressing uses the two bytes following the Opcode as an address.
func Absolute(state State, memory Memory) (Addressing, error) {
	result, err := AbsoluteAddressing(state, memory, 0)
	result.PageBoundaryCrossed = false
	return result, err
}

// AbsoluteX addressing uses the two bytes following the Opcode as a base address to which
// the X register is added with carry.
func AbsoluteX(state State, memory Memory) (Addressing, error) {
	return AbsoluteAddressing(state, memory, Address(state.X))
}

// AbsoluteY addressing uses the two bytes following the Opcode as a base address to which
// the Y register is added with carry.
func AbsoluteY(state State, memory Memory) (Addressing, error) {
	return AbsoluteAddressing(state, memory, Address(state.Y))
}

// Accumulator addressing does no calculations and returns the Accumulator in Addressing.
func Accumulator(state State, memory Memory) (Addressing, error) {
	return Addressing{
		Accumulator: true,
		Value:       state.A,
		Memory:      memory,
	}, nil
}

// Immediate simply reads the next byte from memory.
func Immediate(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	return Addressing{
		EffectiveAddress:     state.PC,
		Value:                memory.Read(state.PC),
		ProgramCounterChange: 1,
		Memory:               memory,
	}, nil
}

// Implied addressing does no calculations and returns a zero Addressing
func Implied(_ State, memory Memory) (Addressing, error) {
	return Addressing{
		Memory: memory,
	}, nil
}

// Indirect addressing uses the two bytes following the Opcode as an address from
// which to read the low byte of the effective address. If reading across a page boundary
// the address provided wraps around that page boundary.
func Indirect(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	indirectAddressLow := memory.Read(state.PC)
	indirectAddressHigh := memory.Read(state.PC + 1)
	indirectAddress := MakeAddress(indirectAddressLow, indirectAddressHigh)

	// Replicate 6502 page-boundary wraparound.
	indirectAddressPlusOne := (indirectAddress & 0xFF00) | ((indirectAddress + 1) & 0x00FF)

	low := memory.Read(indirectAddress)
	high := memory.Read(indirectAddressPlusOne)
	effectiveAddress := MakeAddress(low, high)

	return Addressing{
		EffectiveAddress:     effectiveAddress,
		Value:                memory.Read(effectiveAddress),
		ProgramCounterChange: 2,
		Memory:               memory,
	}, nil
}

// IndirectX addressing uses the byte following the Opcode added to the X register
// to calculate an address in the zero page. The combination of the byte and X will
// wrap around the zero page. The effective address is then read from this address and
// the following byte.
func IndirectX(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}
	// Calculate table pointer, wrapping around the zero page.
	indirectAddress := Address(memory.Read(state.PC)+state.X) & 0x00FF
	low := memory.Read(indirectAddress)
	high := memory.Read(indirectAddress + 1)
	effectiveAddress := MakeAddress(low, high)

	return Addressing{
		EffectiveAddress:     effectiveAddress,
		Value:                memory.Read(effectiveAddress),
		ProgramCounterChange: 1,
		Memory:               memory,
	}, nil
}

// IndirectY addressing uses the byte following the Opcode as the base location of the
// effective address in the zero page (with wraparound). The Y register is then added
// to calculate the actual effective address with carry (so no wrap around).
func IndirectY(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	indirectAddress := Address(memory.Read(state.PC)) & 0x00FF
	// Replicate 6502 page-boundary wraparound.
	indirectAddressPlusOne := (indirectAddress + 1) & 0x00FF

	low := memory.Read(indirectAddress)
	high := memory.Read(indirectAddressPlusOne)
	effectiveAddress := MakeAddress(low, high)

	// Determine if a page boundary has been crossed by the offset
	startPage := effectiveAddress & 0xFF00
	effectiveAddress += Address(state.Y)
	pageBoundaryCrossed := startPage != (effectiveAddress & 0xFF00)

	return Addressing{
		EffectiveAddress:     effectiveAddress,
		Value:                memory.Read(effectiveAddress),
		ProgramCounterChange: 1,
		PageBoundaryCrossed:  pageBoundaryCrossed,
		Memory:               memory,
	}, nil
}

// Relative addressing uses the byte following the opcode as a signed 8-bit value
// to adjust the program counter by.
func Relative(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	// TODO: If the relative address crosses a page boundary then it incurs a penalty.

	// TODO: Relative addressing looks faulty as the Cpu seems to always calculate the penalty.
	//       However, it should only incur the penalty if the branch is taken.

	// Convert the relative address byte to a 16-bit value that we can add to the
	// program counter. We then need to check to see if it represents a negative
	// address (MSB set) and adjust as required.
	value := memory.Read(state.PC)
	relativeAddress := Address(value)
	if (value & 0x80) != 0 {
		relativeAddress |= 0xFF00
	}
	return Addressing{
		EffectiveAddress:     state.PC + 1 + relativeAddress,
		Value:                value,
		ProgramCounterChange: 1,
		Memory:               memory,
	}, nil
}

// ZeroPage addressing uses the byte following the opcode as the address into the
// first page of Memory (i.e. the high byte of the address is always 0x00).
func ZeroPage(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	effectiveAddress := Address(memory.Read(state.PC))
	return Addressing{
		EffectiveAddress:     effectiveAddress,
		Value:                memory.Read(effectiveAddress),
		ProgramCounterChange: 1,
		Memory:               memory,
	}, nil
}

// ZeroPageX addressing uses the byte following the opcode, incremented by X (without
// carry) as the address into the first page of memory (i.e. the high byte of the
// address is always 0x00).
func ZeroPageX(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	effectiveAddress := Address((memory.Read(state.PC) + state.X) & 0xFF)
	return Addressing{

		EffectiveAddress:     effectiveAddress,
		Value:                memory.Read(effectiveAddress),
		ProgramCounterChange: 1,
		Memory:               memory,
	}, nil
}

// ZeroPageY addressing uses the byte following the opcode, incremented by Y (without
// carry) as the address into the first page of memory (i.e. the high byte of the
// address is always 0x00).
func ZeroPageY(state State, memory Memory) (Addressing, error) {
	if memory == nil {
		return Addressing{}, MemoryMustBeProvided
	}

	effectiveAddress := Address((memory.Read(state.PC) + state.Y) & 0xFF)
	return Addressing{
		EffectiveAddress:     effectiveAddress,
		Value:                memory.Read(effectiveAddress),
		ProgramCounterChange: 1,
		Memory:               memory,
	}, nil
}

// ************************************************************
// ********** End of addressing functions
// ************************************************************
