package processor

type Memory interface {
	Read(Address) uint8
	Write(Address, uint8)
}

/*
The 6502 CPU expects interrupt vectors in a fixed place at the end of the memory space:
$FFFA–$FFFB: NMI vector
$FFFC–$FFFD: Reset vector
$FFFE–$FFFF: IRQ/BRK vector
*/

const (
	NmiVectorAddress   = 0xFFFA
	ResetVectorAddress = 0xFFFC
	IrqVectorAddress   = 0xFFFE
)

// WriteVectorToMemory writes the passed in vector to the memory
// locations address and address+1; low byte in address , high byte in address+1.
func WriteVectorToMemory(memory Memory, address Address, vector Address) error {
	if memory == nil {
		return MemoryMustBeProvided
	}

	lowByte, highByte := SplitAddress(vector)
	memory.Write(address, lowByte)
	memory.Write(address+1, highByte)

	return nil
}

// ReadVectorFromMemory returns the vector in the memory
// locations address and address+1; low byte in address , high byte in address+1.
func ReadVectorFromMemory(memory Memory, address Address) (Address, error) {
	if memory == nil {
		return 0, MemoryMustBeProvided
	}

	lowByte := memory.Read(address)
	highByte := memory.Read(address + 1)

	return MakeAddress(lowByte, highByte), nil
}

// WriteNmiVectorToMemory writes the passed in vector to the memory
// locations 0xFFFA and 0xFFFB; low byte in 0xFFFA.
func WriteNmiVectorToMemory(memory Memory, vector Address) error {
	return WriteVectorToMemory(memory, NmiVectorAddress, vector)
}

// ReadNmiVectorFromMemory returns the vector in the memory
// locations 0xFFFA and 0xFFFB; low byte in 0xFFFA.
func ReadNmiVectorFromMemory(memory Memory) (Address, error) {
	return ReadVectorFromMemory(memory, NmiVectorAddress)
}

// WriteResetVectorToMemory writes the passed in vector to the memory
// locations 0xFFFC and 0xFFFD; low byte in 0xFFFC.
func WriteResetVectorToMemory(memory Memory, vector Address) error {
	return WriteVectorToMemory(memory, ResetVectorAddress, vector)
}

// ReadResetVectorFromMemory returns the vector in the memory
// locations 0xFFFC and 0xFFFD; low byte in 0xFFFC.
func ReadResetVectorFromMemory(memory Memory) (Address, error) {
	return ReadVectorFromMemory(memory, ResetVectorAddress)
}

// WriteIrqVectorToMemory writes the passed in vector to the memory
// locations 0xFFFF and 0xFFFE; low byte in 0xFFFE.
func WriteIrqVectorToMemory(memory Memory, vector Address) error {
	return WriteVectorToMemory(memory, IrqVectorAddress, vector)
}

// ReadIrqVectorFromMemory returns the vector in the memory
// locations 0xFFFF and 0xFFFE; low byte in 0xFFFE.
func ReadIrqVectorFromMemory(memory Memory) (Address, error) {
	return ReadVectorFromMemory(memory, IrqVectorAddress)
}

// WriteContiguousDataToMemory copies the contents of data to a contiguous area
// of memory with the first byte specified by address.
func WriteContiguousDataToMemory(memory Memory, start Address, data []uint8) error {
	if memory == nil {
		return MemoryMustBeProvided
	}

	// Push the values into RAM
	for _, v := range data {
		memory.Write(start, v)
		start++
	}

	return nil
}

// WriteDataToMemory can be used to write individual bytes of memory to specific locations
func WriteDataToMemory(memory Memory, data map[Address]uint8) error {
	if memory == nil {
		return MemoryMustBeProvided
	}

	// Push the values into RAM
	for k, v := range data {
		memory.Write(k, v)
	}

	return nil
}

type RepeatingRamSize Address

const (
	EightBytes                    RepeatingRamSize = 0x08
	SixteenBytes                  RepeatingRamSize = 0x10
	ThirtyTwoBytes                RepeatingRamSize = 0x20
	SixtyFourBytes                RepeatingRamSize = 0x40
	OneHundredAndTwentyEightBytes RepeatingRamSize = 0x80
	OneKiloByte                   RepeatingRamSize = 0x400
)

// RepeatingRam allows a smaller RAM size to be repeated across the entire 64K
// address space. Valid sizes are defined by the RepeatingRamSize constants.
type RepeatingRam struct {
	ram  []uint8
	size Address
	mask Address
}

// Read a value from the repeating RAM.
func (r *RepeatingRam) Read(address Address) uint8 {
	if r == nil {
		return 0
	}
	return r.ram[r.mask&address]
}

// Write a value to the repeating RAM.
func (r *RepeatingRam) Write(address Address, value uint8) {
	if r == nil {
		return
	}
	r.ram[r.mask&address] = value
}

// NewRepeatingRam returns a RepeatingRam that will repeat across the full 64K address
// space. Only certain sizes that repeat within the 64K address space are allowed.
func NewRepeatingRam(size RepeatingRamSize) (RepeatingRam, error) {
	switch size {
	case
		EightBytes,
		SixteenBytes,
		ThirtyTwoBytes,
		SixtyFourBytes,
		OneHundredAndTwentyEightBytes,
		OneKiloByte:
		result := RepeatingRam{
			ram:  make([]uint8, size),
			size: Address(size),
			mask: Address(size - 1),
		}
		return result, nil
	}

	return RepeatingRam{}, InvalidMemorySizeProvided
}

// NewPopulatedRam constructs a simple RAM that is pre-populated with the
// specified data starting from zero. This is just a convenience function
// to make it simpler to write the tests. It will simply panic if there
// are any errors.
func NewPopulatedRam(size RepeatingRamSize, data []uint8) RepeatingRam {
	ram, err := NewRepeatingRam(size)
	if err != nil {
		panic(err)
	}
	if data != nil {
		if err = WriteContiguousDataToMemory(&ram, 0, data); err != nil {
			panic(err)
		}
	}
	return ram
}
