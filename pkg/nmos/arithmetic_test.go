package nmos

import (
	"testing"
)

// Tests for documented and undocumented decimal mode:
// From: APPENDIX B: A PROGRAM TO VERIFY DECIMAL MODE OPERATION
// At: http://www.6502.org/tutorials/decimal_mode.html#B
// See also: https://github.com/Klaus2m5/6502_65C02_functional_tests/blob/master/6502_decimal_test.a65
//
// The program works by testing all 262144 = 2 * 2 * 256 * 256 cases:
// * 2 instructions (ADC and SBC),
// * 2 possible values for the carry flag beforehand,
// * 256 possible values for the first number (of the two numbers being added or subtracted),
// * 256 possible values for the second number.
//
// 0 is stored in ERROR if successful, and 1 is stored in ERROR if unsuccessful.
//
// This is done not by executing the program from APPENDIX B, but by generating all of
// the possible combinations as tests.
func TestBCDWithAccumulator(t *testing.T) {
	// TODO: Add tests
}
