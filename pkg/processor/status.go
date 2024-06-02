package processor

import "fmt"

const FlagNone = 0x00
const FlagCarry = 0x01
const FlagNoBorrow = 0x01
const FlagZero = 0x02
const FlagInterrupt = 0x04
const FlagDecimal = 0x08
const FlagBreak = 0x10
const FlagConstant = 0x20
const FlagOverflow = 0x40
const FlagNegative = 0x80
const FlagAll = 0xFF

// Status is the 8-bit want register that makes up the Flags.
type Status uint8

// WithConstantSet returns the status flag with the constant flag set.
func (s *Status) WithConstantSet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagConstant
}

// WithConstantCleared returns the status flag with the constant flag cleared.
func (s *Status) WithConstantCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagConstant)
}

// WithCarrySet returns the status flag with the carry flag set.
func (s *Status) WithCarrySet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagCarry
}

// WithCarryCleared returns the status flag with the carry flag cleared.
func (s *Status) WithCarryCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagCarry)
}

// WithZeroSet returns the status flag with the zero flag set.
func (s *Status) WithZeroSet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagZero
}

// WithZeroCleared returns the status flag with the zero flag cleared.
func (s *Status) WithZeroCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagZero)
}

// WithInterruptSet returns the status flag with the interrupt flag set.
func (s *Status) WithInterruptSet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagInterrupt
}

// WithInterruptCleared returns the status flag with the interrupt flag cleared.
func (s *Status) WithInterruptCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagInterrupt)
}

// WithDecimalSet returns the status flag with the decimal flag set.
func (s *Status) WithDecimalSet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagDecimal
}

// WithDecimalCleared returns the status flag with the decimal flag cleared.
func (s *Status) WithDecimalCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagDecimal)
}

// WithBreakSet returns the status flag with the break flag set.
func (s *Status) WithBreakSet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagBreak
}

// WithBreakCleared returns the status flag with the break flag cleared.
func (s *Status) WithBreakCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagBreak)
}

// WithOverflowSet returns the status flag with the overflow flag set.
func (s *Status) WithOverflowSet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagOverflow
}

// WithOverflowCleared returns the status flag with the overflow flag cleared.
func (s *Status) WithOverflowCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagOverflow)
}

// WithNegativeSet returns the status flag with the sign flag set.
func (s *Status) WithNegativeSet() Status {
	if s == nil {
		return 0
	}
	return *s | FlagNegative
}

// WithNegativeCleared returns the status flag with the sign flag cleared.
func (s *Status) WithNegativeCleared() Status {
	if s == nil {
		return 0
	}
	return *s & (0xFF ^ FlagNegative)
}

// SetCarry modifies and returns the status flag with the carry flag set.
func (s *Status) SetCarry() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithCarrySet()
	return s
}

// ClearCarry modifies and returns the status flag with the carry flag cleared.
func (s *Status) ClearCarry() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithCarryCleared()
	return s
}

// SetZero modifies and returns the status flag with the zero flag set.
func (s *Status) SetZero() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithZeroSet()
	return s
}

// ClearZero modifies and returns the status flag with the zero flag cleared.
func (s *Status) ClearZero() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithZeroCleared()
	return s
}

// SetInterrupt modifies and returns the status flag with the interrupt flag set.
func (s *Status) SetInterrupt() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithInterruptSet()
	return s
}

// ClearInterrupt modifies and returns the status flag with the interrupt flag cleared.
func (s *Status) ClearInterrupt() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithInterruptCleared()
	return s
}

// SetDecimal modifies and returns the status flag with the decimal flag set.
func (s *Status) SetDecimal() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithDecimalSet()
	return s
}

// ClearDecimal modifies and returns the status flag with the decimal flag cleared.
func (s *Status) ClearDecimal() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithDecimalCleared()
	return s
}

// SetBreak modifies and returns the status flag with the break flag set.
func (s *Status) SetBreak() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithBreakSet()
	return s
}

// SetConstant modifies and returns the status flag with the constant flag set.
func (s *Status) SetConstant() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithConstantSet()
	return s
}

// ClearBreak modifies and returns the status flag with the break flag cleared.
func (s *Status) ClearBreak() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithBreakCleared()
	return s
}

// ClearConstant modifies and returns the status flag with the constant flag cleared.
func (s *Status) ClearConstant() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithConstantCleared()
	return s
}

// SetOverflow modifies and returns the status flag with the overflow flag set.
func (s *Status) SetOverflow() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithOverflowSet()
	return s
}

// ClearOverflow modifies and returns the status flag with the overflow flag cleared.
func (s *Status) ClearOverflow() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithOverflowCleared()
	return s
}

// SetNegative modifies and returns the status flag with the sign flag set.
func (s *Status) SetNegative() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithNegativeSet()
	return s
}

// ClearNegative modifies and returns the status flag with the sign flag cleared.
func (s *Status) ClearNegative() *Status {
	if s == nil {
		return nil
	}
	*s = s.WithNegativeCleared()
	return s
}

// ToFlags returns an instance of Flags with want flags set
// according to bits on the Status instance.
func (s *Status) ToFlags() Flags {
	if s == nil {
		return Flags{}
	}
	return Flags{
		Carry:     (*s & FlagCarry) != 0,
		Zero:      (*s & FlagZero) != 0,
		Interrupt: (*s & FlagInterrupt) != 0,
		Decimal:   (*s & FlagDecimal) != 0,
		Break:     (*s & FlagBreak) != 0,
		Overflow:  (*s & FlagOverflow) != 0,
		Negative:  (*s & FlagNegative) != 0,
	}
}

// String returns a string representation of the flags. This
// is the same representation as an equivalent Flags instance.
func (s *Status) String() string {
	if s == nil {
		return ""
	}
	return s.ToFlags().String()
}

// Flags represents the individual flags within a Status register
// as booleans to make it easier to work with.
type Flags struct {
	Carry     bool
	Zero      bool
	Interrupt bool
	Decimal   bool
	Break     bool
	Overflow  bool
	Negative  bool
}

// ToStatus converts the Flags instance to the 8-bit status flag
// representation.
func (f Flags) ToStatus() Status {
	result := Status(FlagConstant)
	if f.Carry {
		result |= FlagCarry
	}
	if f.Zero {
		result |= FlagZero
	}
	if f.Interrupt {
		result |= FlagInterrupt
	}
	if f.Decimal {
		result |= FlagDecimal
	}
	if f.Break {
		result |= FlagBreak
	}
	if f.Overflow {
		result |= FlagOverflow
	}
	if f.Negative {
		result |= FlagNegative
	}
	return result
}

// Converts the Flag instance into a canonical string form.
func (f Flags) String() string {
	carry := 'c'
	zero := 'z'
	interrupt := 'i'
	decimal := 'd'
	brk := 'b'
	overflow := 'v'
	negative := 'n'

	if f.Carry {
		carry = 'C'
	}
	if f.Zero {
		zero = 'Z'
	}
	if f.Interrupt {
		interrupt = 'I'
	}
	if f.Decimal {
		decimal = 'D'
	}
	if f.Break {
		brk = 'B'
	}
	if f.Overflow {
		overflow = 'V'
	}
	if f.Negative {
		negative = 'N'
	}
	return fmt.Sprintf(
		"%c %c %c %c %c %c %c %c",
		negative, overflow, '-', brk, decimal, interrupt, zero, carry)
}
