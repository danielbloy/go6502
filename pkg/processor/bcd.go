package processor

// ExtractBcdValues returns both digits from a single byte. The digit from the high
// nibble is returned first. Both returned values will be in the range of 0x0 to 0xF.
// Naturally, any value > 9 can be ignored.
func ExtractBcdValues(value uint8) (uint8, uint8) {
	vl := value & 0x0F
	vh := (value & 0xF0) >> 4
	return vh, vl
}

// DecodeBcdValue Converts the byte containing two BCD digits to its decimal form.
// If either of the digits exceed the 0-9 range then 0xFF is returned.
func DecodeBcdValue(value uint8) uint8 {
	h, l := ExtractBcdValues(value)

	if h > 9 || l > 9 {
		return 0xFF
	}
	return (h * 10) + l
}

// EncodeBcdValue converts a number from 0 to 99 to it's BCD form. Any value provided
// that is greater than 99 is reduced to the range of 0 to 99.
func EncodeBcdValue(value uint8) uint8 {
	for value > 99 {
		value -= 100
	}

	h := uint8(0)
	for value > 9 {
		h++
		value -= 10
	}

	return (h << 4) + value
}
