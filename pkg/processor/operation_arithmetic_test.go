package processor

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_isCarrySet(t *testing.T) {
	tests := []struct {
		name  string
		value uint16
		want  bool
	}{
		{name: "Zero has no carry", value: 0},
		{name: "1 has no carry", value: 1},
		{name: "0xFF has no carry", value: 0xFF},
		{name: "0xFFFF has carry", value: 0xFFFF, want: true},
		{name: "0x01FF has carry", value: 0x01FF, want: true},
		{name: "0x10FF has carry", value: 0x10FF, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCarrySet(tt.value); got != tt.want {
				t.Errorf("isCarrySet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_carrySet(t *testing.T) {
	tests := []struct {
		name      string
		value     uint16
		wantCarry bool
	}{
		{name: "Zero does not set carry", value: 0},
		{name: "1 does not set carry", value: 1},
		{name: "0xFF does not set carry", value: 0xFF},
		{name: "0xFFFF does set carry", value: 0xFFFF, wantCarry: true},
		{name: "0x01FF does set carry", value: 0x01FF, wantCarry: true},
		{name: "0x10FF does set carry", value: 0x10FF, wantCarry: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status Status
			carrySet(&status, tt.value)
			if got := isCarrySet(tt.value); got != tt.wantCarry {
				t.Errorf("carrySet() = %v, want %v", got, tt.wantCarry)
			}
		})
	}
}

func Test_isNegative(t *testing.T) {
	tests := []struct {
		name  string
		value uint16
		want  bool
	}{
		{name: "Zero is not negative", value: 0},
		{name: "1 is not negative", value: 1},
		{name: "0x7F is not negative", value: 0x7F},
		{name: "0x80 is negative", value: 0x80, want: true},
		{name: "0x90 is negative", value: 0x90, want: true},
		{name: "0xFF is negative carry", value: 0xFF, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNegative(tt.value); got != tt.want {
				t.Errorf("isNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_negativeSet(t *testing.T) {
	tests := []struct {
		name         string
		value        uint16
		wantNegative bool
	}{
		{name: "Zero does not set negative", value: 0},
		{name: "1 does not set negative", value: 1},
		{name: "0x7F does not set negative", value: 0x7F},
		{name: "0x80 does set negative", value: 0x80, wantNegative: true},
		{name: "0x90 does se negative", value: 0x90, wantNegative: true},
		{name: "0xFF does se negative carry", value: 0xFF, wantNegative: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status Status
			negativeSet(&status, tt.value)
			if got := isNegative(tt.value); got != tt.wantNegative {
				t.Errorf("negativeSet() = %v, want %v", got, tt.wantNegative)
			}
		})
	}
}

func Test_isZero(t *testing.T) {
	tests := []struct {
		name  string
		value uint16
		want  bool
	}{
		{name: "Zero is zero", value: 0, want: true},
		{name: "0xFF00 is zero", value: 0xFF00, want: true},
		{name: "0x0100 is zero", value: 0x0100, want: true},
		{name: "1 is not zero", value: 1},
		{name: "0x7F is not zero", value: 0x7F},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isZero(tt.value); got != tt.want {
				t.Errorf("isZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zeroSet(t *testing.T) {
	tests := []struct {
		name     string
		value    uint16
		wantZero bool
	}{
		{name: "Zero does set zero", value: 0, wantZero: true},
		{name: "0xFF00 does set zero", value: 0xFF00, wantZero: true},
		{name: "0x0100 does set zero", value: 0x0100, wantZero: true},
		{name: "1 does not set zero", value: 1},
		{name: "0x7F does not set zero", value: 0x7F}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status Status
			zeroSet(&status, tt.value)
			if got := isZero(tt.value); got != tt.wantZero {
				t.Errorf("zeroSet() = %v, want %v", got, tt.wantZero)
			}
		})
	}
}

func Test_isOverflow(t *testing.T) {
	tests := []struct {
		name    string
		decimal bool
		start   uint16
		value   uint16
		result  uint16
		want    bool
	}{
		{name: "All zeros does not overflow"},
		/* These test cases are from http://www.6502.org/tutorials/vflag.html
		CLC      ; 1 + 1 = 2, returns V = 0
		LDA #$01
		ADC #$01

		CLC      ; 1 + -1 = 0, returns V = 0
		LDA #$01
		ADC #$FF

		CLC      ; 127 + 1 = 128, returns V = 1
		LDA #$7F
		ADC #$01

		CLC      ; -128 + -1 = -129, returns V = 1
		LDA #$80
		ADC #$FF

		SEC      ; 0 - 1 = -1, returns V = 0
		LDA #$00
		SBC #$01

		SEC      ; -128 - 1 = -129, returns V = 1
		LDA #$80
		SBC #$01

		SEC      ; 127 - -1 = 128, returns V = 1
		LDA #$7F
		SBC #$FF

		SEC      ; Note: SEC, not CLC
		LDA #$3F ; 63 + 64 + 1 = 128, returns V = 1
		ADC #$40

		CLC      ; Note: CLC, not SEC
		LDA #$C0 ; -64 - 64 - 1 = -129, returns V = 1
		SBC #$40
		*/
		// Binary mode tests
		{name: "Binary ; 1 + 1 = 2, returns V = 0", start: 1, value: 1, result: 2, want: false},
		{name: "Binary ; 1 + -1 = 0, returns V = 0", start: 1, value: 0xFF, result: 0, want: false},
		{name: "Binary ; 127 + 1 = 128, returns V = 1", start: 0x7F, value: 1, result: 0x80, want: true},
		{name: "Binary ; -128 + -1 = -129, returns V = 1", start: 0x80, value: 0xFF, result: 0x17F, want: true},
		{name: "Binary ; 0 - 1 = -1, returns V = 0", start: 0, value: (1 ^ 0xFF) + 1, result: 0xFF, want: false},
		{name: "Binary ; -128 - 1 = -129, returns V = 1", start: 0x80, value: (1 ^ 0xFF) + 1, result: 0x17F, want: true},
		{name: "Binary ; 127 - -1 = 128, returns V = 1", start: 0x7F, value: (0xFF ^ 0xFF) + 1, result: 0x80, want: true},
		{name: "Binary ; 1 + 1 = 2, returns V = 0", start: 1, value: 1, result: 2, want: false},
		{name: "Binary ; 63 + 64 + 1 = 128, returns V = 1", start: 0x3F, value: 0x40, result: 0x80, want: true},
		{name: "Binary ; -64 - 64 - 1 = -129, returns V = 1", start: 0xC0, value: (0x40 ^ 0xFF) + 1, result: 0x17F, want: true},
		// Decimal mode tests.
		{name: "Decimal ; 56 + 47 = 105, returns V = 1", decimal: true, start: 0x56, value: 0x47, result: 0x105, want: true},
		{name: "Decimal ; 12 + 34 = 46, returns V = 0", decimal: true, start: 0x12, value: 0x34, result: 0x46, want: false},
		{name: "Decimal ; 15 + 26 = 41, returns V = 0", decimal: true, start: 0x15, value: 0x26, result: 0x41, want: false},
		{name: "Decimal ; 81 + 92 = 173, returns V = 1", decimal: true, start: 0x81, value: 0x92, result: 0x173, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOverflow(tt.decimal, tt.start, tt.value, tt.result); got != tt.want {
				t.Errorf("isOverflow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_overflowSet(t *testing.T) {
	tests := []struct {
		name         string
		decimal      bool
		start        uint16
		value        uint16
		result       uint16
		wantOverflow bool
	}{
		// Binary mode tests.
		{name: "Binary ; 1 + 1 = 2, returns V = 0", start: 1, value: 1, result: 2, wantOverflow: false},
		{name: "Binary ; 1 + -1 = 0, returns V = 0", start: 1, value: 0xFF, result: 0, wantOverflow: false},
		{name: "Binary ; 127 + 1 = 128, returns V = 1", start: 0x7F, value: 1, result: 0x80, wantOverflow: true},
		{name: "Binary ; -128 + -1 = -129, returns V = 1", start: 0x80, value: 0xFF, result: 0x17F, wantOverflow: true},
		// Decimal mode tests.
		{name: "Decimal ; 56 + 47 = 105, returns V = 1", decimal: true, start: 0x56, value: 0x47, result: 0x105, wantOverflow: true},
		{name: "Decimal ; 12 + 34 = 46, returns V = 0", decimal: true, start: 0x12, value: 0x34, result: 0x46, wantOverflow: false},
		{name: "Decimal ; 15 + 26 = 41, returns V = 0", decimal: true, start: 0x15, value: 0x26, result: 0x41, wantOverflow: false},
		{name: "Decimal ; 81 + 92 = 173, returns V = 1", decimal: true, start: 0x81, value: 0x92, result: 0x173, wantOverflow: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status Status
			overflowSet(&status, tt.start, tt.value, tt.result)
			if got := isOverflow(tt.decimal, tt.start, tt.value, tt.result); got != tt.wantOverflow {
				t.Errorf("overflowSet() = %v, want %v", got, tt.wantOverflow)
			}
		})
	}
}

func TestAddWithCarry(t *testing.T) {
	tests := []testOperationConfig{
		// These test cases are from http://www.6502.org/tutorials/vflag.html
		{ /*
				CLC      ; 1 + 1 = 2, returns V = 0
				LDA #$01
				ADC #$01
			*/
			name:       "1 + 1 = 2, returns V = 0",
			startState: State{A: 0x01},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x02},
		},
		{ /*
				CLC      ; 1 + -1 = 0, returns V = 0
				LDA #$01
				ADC #$FF
			*/
			name:       "1 + -1 = 0, returns V = 0, Z = 1, C = 1",
			startState: State{A: 0x01},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x00, P: FlagZero | FlagCarry},
		},
		{ /*
				CLC      ; 127 + 1 = 128, returns V = 1
				LDA #$7F
				ADC #$01
			*/
			name:       "127 + 1 = 128, returns V = 1, N = 1",
			startState: State{A: 0x7F},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x80, P: FlagOverflow | FlagNegative},
		},
		{ /*
				CLC      ; -128 + -1 = -129, returns V = 1
				LDA #$80
				ADC #$FF
			*/
			name:       "-128 + -1 = -129, returns V = 1, C = 1",
			startState: State{A: 0x80},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x7F, P: FlagOverflow | FlagCarry},
		},
		{ /*
				SEC      ; Note: SEC, not CLC
				LDA #$3F ; 63 + 64 + 1 = 128, returns V = 1
				ADC #$40
			*/
			name:       "63 + 64 + 1 (carry) = 128, returns V = 1, N = 1",
			startState: State{A: 0x3F, P: FlagCarry},
			addressing: Addressing{Value: 0x40},
			wantState:  State{A: 0x80, P: FlagOverflow | FlagNegative},
		},
		{
			name:       "Zero added to Zero",
			startState: State{A: 0x00, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x00, P: FlagZero},
		},
		{
			name:       "Zero added to Zero with carry set",
			startState: State{A: 0x00, P: FlagCarry},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x01, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers no carry) that does not set carry - 1",
			startState: State{A: 0x00, P: FlagNone},
			addressing: Addressing{Value: 0x0F},
			wantState:  State{A: 0x0F, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers no carry) that does not set carry - 2",
			startState: State{A: 0x0F, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0F, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers no carry) that does not set carry - 3",
			startState: State{A: 0x0F, P: FlagNone},
			addressing: Addressing{Value: 0x40},
			wantState:  State{A: 0x4F, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers with carry) that does not set carry - A",
			startState: State{A: 0x00, P: FlagCarry},
			addressing: Addressing{Value: 0x0E},
			wantState:  State{A: 0x0F, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers with carry) that does not set carry - B",
			startState: State{A: 0x0E, P: FlagCarry},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0F, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers with carry) that does not set carry - C",
			startState: State{A: 0x0E, P: FlagCarry},
			addressing: Addressing{Value: 0x40},
			wantState:  State{A: 0x4F, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers no carry) that does not set carry - 1",
			startState: State{A: 0x00, P: FlagNone},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x80, P: FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers no carry) that does not set carry - 2",
			startState: State{A: 0x80, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x80, P: FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers no carry) that does not set carry - 3",
			startState: State{A: 0x0F, P: FlagNone},
			addressing: Addressing{Value: 0xE0},
			wantState:  State{A: 0xEF, P: FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers with carry) that does not set carry - A",
			startState: State{A: 0x00, P: FlagCarry},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x81, P: FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers with carry) that does not set carry - B",
			startState: State{A: 0x80, P: FlagCarry},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x81, P: FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers with carry) that does not set carry - C",
			startState: State{A: 0x0E, P: FlagCarry},
			addressing: Addressing{Value: 0xE0},
			wantState:  State{A: 0xEF, P: FlagNegative},
		},
		{
			name:       "Unsigned arithmetic that sets carry (result in accumulator is zero) - 1",
			startState: State{A: 0x01, P: FlagNone},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x00, P: FlagCarry | FlagZero},
		},
		{
			name:       "Unsigned arithmetic (with carry) that sets carry (result in accumulator is zero) - A",
			startState: State{A: 0x01, P: FlagCarry},
			addressing: Addressing{Value: 0xFE},
			wantState:  State{A: 0x00, P: FlagCarry | FlagZero},
		},
		{
			name:       "Unsigned arithmetic that sets carry (result in accumulator is unsigned) - 1",
			startState: State{A: 0xFF, P: FlagNone},
			addressing: Addressing{Value: 0x02},
			wantState:  State{A: 0x01, P: FlagCarry},
		},
		{
			name:       "Unsigned arithmetic (with carry) that sets carry (result in accumulator is unsigned) - A",
			startState: State{A: 0xFF, P: FlagCarry},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x01, P: FlagCarry},
		},
		{
			name:       "Unsigned arithmetic that sets carry (result in accumulator is signed)",
			startState: State{A: 0xFF, P: FlagNone},
			addressing: Addressing{Value: 0x81},
			wantState:  State{A: 0x80, P: FlagCarry | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (with carry) that sets carry (result in accumulator is signed)",
			startState: State{A: 0xFF, P: FlagCarry},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x80, P: FlagCarry | FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers no carry) that does not set carry - 1",
			startState: State{A: 0x00, P: FlagNone},
			addressing: Addressing{Value: 0x8F},
			wantState:  State{A: 0x8F, P: FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers no carry) that does not set carry - 2",
			startState: State{A: 0x8F, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x8F, P: FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers no carry) that does not set carry - 3",
			startState: State{A: 0x0F, P: FlagNone},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x8F, P: FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers with carry) that does not set carry - A",
			startState: State{A: 0x00, P: FlagCarry},
			addressing: Addressing{Value: 0x8E},
			wantState:  State{A: 0x8F, P: FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers with carry) that does not set carry - B",
			startState: State{A: 0x8E, P: FlagCarry},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x8F, P: FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers with carry) that does not set carry - C",
			startState: State{A: 0x0E, P: FlagCarry},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x8F, P: FlagNegative},
		},
		{
			name:       "Signed arithmetic that sets carry (result in accumulator is zero) - 1",
			startState: State{A: 0x80, P: FlagNone},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x00, P: FlagCarry | FlagZero | FlagOverflow},
		},
		{
			name:       "Signed arithmetic (with carry) that sets carry (result in accumulator is zero) - A",
			startState: State{A: 0x80, P: FlagCarry},
			addressing: Addressing{Value: 0x7F},
			wantState:  State{A: 0x00, P: FlagCarry | FlagZero},
		},
		{
			name:       "Signed arithmetic that sets carry (result in accumulator is unsigned) - 1",
			startState: State{A: 0x80, P: FlagNone},
			addressing: Addressing{Value: 0x81},
			wantState:  State{A: 0x01, P: FlagCarry | FlagOverflow},
		},
		{
			name:       "Signed arithmetic (with carry) that sets carry (result in accumulator is unsigned) - A",
			startState: State{A: 0x80, P: FlagCarry},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x01, P: FlagCarry | FlagOverflow},
		},
		{
			name:       "Signed arithmetic that sets carry (result in accumulator is signed)",
			startState: State{A: 0xFF, P: FlagNone},
			addressing: Addressing{Value: 0x81},
			wantState:  State{A: 0x80, P: FlagCarry | FlagNegative},
		},
		{
			name:       "Signed arithmetic (with carry) that sets carry (result in accumulator is signed)",
			startState: State{A: 0xFF, P: FlagCarry},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x80, P: FlagCarry | FlagNegative},
		},
		// Tests from: http://www.6502.org/tutorials/decimal_mode.html#3.2.1
		{
			/*
				SED      ; Decimal mode (BCD addition: 58 + 46 + 1 = 105)
				SEC      ; Note: carry is set, not clear!
				LDA #$58
				ADC #$46 ; After this instruction, C = 1, A = $05
			*/
			name:       "Decimal mode (BCD addition: 58 + 46 + 1 = 105)",
			startState: State{A: 0x58, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x46},
			wantState:  State{A: 0x05, P: FlagDecimal | FlagCarry | FlagOverflow},
		},
		{
			/*
				SED      ; Decimal mode (BCD addition: 12 + 34 = 46)
				CLC
				LDA #$12
				ADC #$34 ; After this instruction, C = 0, A = $46
			*/
			name:       "Decimal mode (BCD addition: 12 + 34 = 46)",
			startState: State{A: 0x12, P: FlagDecimal},
			addressing: Addressing{Value: 0x34},
			wantState:  State{A: 0x46, P: FlagDecimal},
		},
		{
			/*
				SED      ; Decimal mode (BCD addition: 15 + 26 = 41)
				CLC
				LDA #$15
				ADC #$26 ; After this instruction, C = 0, A = $41

			*/
			name:       "Decimal mode (BCD addition: 15 + 26 = 41)",
			startState: State{A: 0x15, P: FlagDecimal},
			addressing: Addressing{Value: 0x26},
			wantState:  State{A: 0x41, P: FlagDecimal},
		},
		{
			/*
				SED      ; Decimal mode (BCD addition: 81 + 92 = 173)
				CLC
				LDA #$81
				ADC #$92 ; After this instruction, C = 1, A = $73
			*/
			name:       "Decimal mode (BCD addition: 81 + 92 = 173)",
			startState: State{A: 0x81, P: FlagDecimal},
			addressing: Addressing{Value: 0x92},
			wantState:  State{A: 0x73, P: FlagDecimal | FlagCarry | FlagOverflow},
		},
		// Other BCD tests.
		{
			name:       "BCD zero add zero is zero",
			startState: State{A: 0x00, P: FlagDecimal},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagZero},
		},
		{
			name:       "BCD zero add 1 is 1",
			startState: State{A: 0x00, P: FlagDecimal},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x01, P: FlagDecimal},
		},
		{
			name:       "BCD 1 add zero is 1",
			startState: State{A: 0x01, P: FlagDecimal},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x01, P: FlagDecimal},
		},
		{
			name:       "BCD zero add $0A is $10",
			startState: State{A: 0x00, P: FlagDecimal},
			addressing: Addressing{Value: 0x0A},
			wantState:  State{A: 0x10, P: FlagDecimal},
		},
		{
			name:       "BCD $0B add zero is $11",
			startState: State{A: 0x0B, P: FlagDecimal},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x11, P: FlagDecimal},
		},
		{
			name:       "BCD zero add $0C is $12",
			startState: State{A: 0x00, P: FlagDecimal},
			addressing: Addressing{Value: 0x0C},
			wantState:  State{A: 0x12, P: FlagDecimal},
		},
		{
			name:       "BCD $0D add zero is $13",
			startState: State{A: 0x0D, P: FlagDecimal},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x13, P: FlagDecimal},
		},
		{
			name:       "BCD zero add $0E is $14",
			startState: State{A: 0x00, P: FlagDecimal},
			addressing: Addressing{Value: 0x0E},
			wantState:  State{A: 0x14, P: FlagDecimal},
		},
		{
			name:       "BCD $0F add zero is $15",
			startState: State{A: 0x0F, P: FlagDecimal},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x15, P: FlagDecimal},
		},
		{
			name:       "BCD $01 add $0C is $13",
			startState: State{A: 0x01, P: FlagDecimal},
			addressing: Addressing{Value: 0x0C},
			wantState:  State{A: 0x13, P: FlagDecimal},
		},
		{
			name:       "BCD $05 add $05 is $10",
			startState: State{A: 0x05, P: FlagDecimal},
			addressing: Addressing{Value: 0x05},
			wantState:  State{A: 0x10, P: FlagDecimal},
		},
		{
			name:       "BCD $08 add $08 is $16",
			startState: State{A: 0x08, P: FlagDecimal},
			addressing: Addressing{Value: 0x08},
			wantState:  State{A: 0x16, P: FlagDecimal},
		},
		{
			name:       "BCD $10 add $20 is $30",
			startState: State{A: 0x10, P: FlagDecimal},
			addressing: Addressing{Value: 0x20},
			wantState:  State{A: 0x30, P: FlagDecimal},
		},
		{
			name:       "BCD $16 add $25 is $41",
			startState: State{A: 0x16, P: FlagDecimal},
			addressing: Addressing{Value: 0x25},
			wantState:  State{A: 0x41, P: FlagDecimal},
		},
		{
			name:       "BCD $99 add $01 is $00",
			startState: State{A: 0x99, P: FlagDecimal},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagCarry | FlagZero},
		},
		// Test that the carry flag works with BCD.
		{
			name:       "BCD $00 add $00 with Carry is $01",
			startState: State{A: 0x00, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x01, P: FlagDecimal},
		},
		{
			name:       "BCD $02 add $03 with Carry is $06",
			startState: State{A: 0x02, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x03},
			wantState:  State{A: 0x06, P: FlagDecimal},
		},
		{
			name:       "BCD $05 add $04 with Carry is $10",
			startState: State{A: 0x05, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x04},
			wantState:  State{A: 0x10, P: FlagDecimal},
		},
		{
			name:       "BCD $20 add $40 with Carry is $61",
			startState: State{A: 0x20, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x40},
			wantState:  State{A: 0x61, P: FlagDecimal},
		},
		{
			name:       "BCD $99 add $00 with Carry is $00",
			startState: State{A: 0x99, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagCarry | FlagZero},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, AddWithCarry)
		})
	}
}

func TestSubtractWithCarry(t *testing.T) {
	tests := []testOperationConfig{
		// These test cases are from http://www.6502.org/tutorials/vflag.html
		{ /*
			  SEC      ; 0 - 1 = -1, returns V = 0
			  LDA #$00
			  SBC #$01
			*/
			name:       "0 - 1 = -1, returns V = 0",
			startState: State{A: 0x00, P: FlagCarry},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0xFF, P: FlagNegative},
		},
		{ /*
			  SEC      ; -128 - 1 = -129, returns V = 1
			  LDA #$80
			  SBC #$01
			*/
			name:       "-128 - 1 = -129, returns V = 1",
			startState: State{A: 0x80, P: FlagCarry},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x7F, P: FlagCarry | FlagOverflow},
		},
		{ /*
				SEC      ; 127 - -1 = 128, returns V = 1
				LDA #$7F
				SBC #$FF
			*/
			name:       "127 - -1 = 128, returns V = 1",
			startState: State{A: 0x7F, P: FlagCarry},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x80, P: FlagNegative | FlagOverflow},
		},
		{ /*
			  CLC      ; Note: CLC, not SEC
			  LDA #$C0 ; -64 - 64 - 1 = -129, returns V = 1
			  SBC #$40
			*/
			name:       "-64 - 64 - 1 = -129, returns V = 1",
			startState: State{A: 0xC0},
			addressing: Addressing{Value: 0x40},
			wantState:  State{A: 0x7F, P: FlagCarry | FlagOverflow},
		},
		{
			name:       "Zero subtracted from Zero (no borrow)",
			startState: State{A: 0x00, P: FlagCarry}, // No borrow
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x00, P: FlagZero | FlagCarry}, // Did not borrow
		},
		{
			name:       "Zero subtracted from Zero with borrow",
			startState: State{A: 0x00, P: FlagNone}, // Borrow set
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0xFF, P: FlagNone | FlagNegative}, // Did borrow
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers no borrw) that does not borrow - 1",
			startState: State{A: 0x0F, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0F, P: FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers no borrow) that does not borrow - 2",
			startState: State{A: 0x0F, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0F, P: FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers no borrow) that does not set carry - 3",
			startState: State{A: 0x40, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x0F},
			wantState:  State{A: 0x31, P: FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers with borrow) that does not set borrow - A",
			startState: State{A: 0x0F, P: FlagNone},
			addressing: Addressing{Value: 0x0E},
			wantState:  State{A: 0x00, P: FlagZero | FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers with borrow) that does not set borrow - B",
			startState: State{A: 0x0F, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0E, P: FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic (unsigned numbers with borrow) that does not set borrow - C",
			startState: State{A: 0x4F, P: FlagNone},
			addressing: Addressing{Value: 0x40},
			wantState:  State{A: 0x0E, P: FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers no borrow) that does not set borrow - 1",
			startState: State{A: 0x80, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x80, P: FlagNoBorrow | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers no borrow) that does not set borrow - 2",
			startState: State{A: 0x90, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x10},
			wantState:  State{A: 0x80, P: FlagNoBorrow | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers no borrow) that does not set borrow - 3",
			startState: State{A: 0xEF, P: FlagNoBorrow},
			addressing: Addressing{Value: 0xE0},
			wantState:  State{A: 0x0F, P: FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers with borrow) that does not set borrow - A",
			startState: State{A: 0x81, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x80, P: FlagNoBorrow | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers with borrow) that does not set borrow - B",
			startState: State{A: 0x81, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x80, P: FlagNoBorrow | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (possible signed numbers with borrow) that does not set borrow - C",
			startState: State{A: 0xEF, P: FlagNone},
			addressing: Addressing{Value: 0xE0},
			wantState:  State{A: 0x0E, P: FlagNoBorrow},
		},
		{
			name:       "Unsigned arithmetic that sets borrow (result in accumulator is FF) - 1",
			startState: State{A: 0x01, P: FlagNone},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0xFF, P: FlagNone | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (with borrow) that sets borrow (result in accumulator is FF) - A",
			startState: State{A: 0xFF, P: FlagNone},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0xFF, P: FlagNone | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic that sets borrow (result in accumulator is unsigned) - 1",
			startState: State{A: 0x02, P: FlagNoBorrow},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x03, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic (with borrow) that sets borrow (result in accumulator is unsigned) - A",
			startState: State{A: 0x01, P: FlagNone},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x01, P: FlagNone},
		},
		{
			name:       "Unsigned arithmetic that sets borrow (result in accumulator is signed)",
			startState: State{A: 0x80, P: FlagNoBorrow},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x81, P: FlagNone | FlagNegative},
		},
		{
			name:       "Unsigned arithmetic (with borrow) that sets borrow (result in accumulator is signed)",
			startState: State{A: 0x80, P: FlagNone},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x80, P: FlagNone | FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers no borrow) that does not borrow - 1",
			startState: State{A: 0x8F, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x8F, P: FlagNoBorrow | FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers no borrow) that does not borrow - 2",
			startState: State{A: 0x8F, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x0F},
			wantState:  State{A: 0x80, P: FlagNoBorrow | FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers no borrow) that does not borrow - 3",
			startState: State{A: 0x8F, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x0F, P: FlagNoBorrow},
		},
		{
			name:       "Signed arithmetic (signed numbers with borrow) that does not borrow - A",
			startState: State{A: 0x8F, P: FlagNone},
			addressing: Addressing{Value: 0x8E},
			wantState:  State{A: 0x00, P: FlagNoBorrow | FlagZero},
		},
		{
			name:       "Signed arithmetic (signed numbers with borrow) that does not borrow - B",
			startState: State{A: 0x8F, P: FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x8E, P: FlagNoBorrow | FlagNegative},
		},
		{
			name:       "Signed arithmetic (signed numbers with borrow) that does not borrow - C",
			startState: State{A: 0x8F, P: FlagNone},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0x0E, P: FlagNoBorrow},
		},
		{
			name:       "Signed arithmetic that borrows (result in accumulator is 0xFF) - 1",
			startState: State{A: 0x80, P: FlagNoBorrow},
			addressing: Addressing{Value: 0x81},
			wantState:  State{A: 0xFF, P: FlagNone | FlagNegative},
		},
		{
			name:       "Signed arithmetic (with borrow) that borrows (result in accumulator is 0xFF) - A",
			startState: State{A: 0x80, P: FlagNone},
			addressing: Addressing{Value: 0x80},
			wantState:  State{A: 0xFF, P: FlagNone | FlagNegative},
		},
		{
			name:       "Signed arithmetic that that borrows (result in accumulator is unsigned) - 1",
			startState: State{A: 0x40, P: FlagNoBorrow},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x41, P: FlagNone},
		},
		{
			name:       "Signed arithmetic (with borrow) that borrows (result in accumulator is unsigned) - A",
			startState: State{A: 0x01, P: FlagNone},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x01, P: FlagNone},
		},
		{
			name:       "Signed arithmetic that borrows (result in accumulator is signed)",
			startState: State{A: 0x80, P: FlagNoBorrow},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x81, P: FlagNone | FlagNegative},
		},
		{
			name:       "Signed arithmetic (with borrow) that borrows (result in accumulator is signed)",
			startState: State{A: 0x80, P: FlagNone},
			addressing: Addressing{Value: 0xFF},
			wantState:  State{A: 0x80, P: FlagNone | FlagNegative},
		},
		// Tests from: http://www.6502.org/tutorials/decimal_mode.html#3.2.1
		{
			/*
				SED      ; Decimal mode (BCD subtraction: 46 - 12 = 34)
				SEC
				LDA #$46
				SBC #$12 ; After this instruction, C = 1, A = $34)
			*/
			name:       "Decimal mode (BCD subtraction: 46 - 12 = 34)",
			startState: State{A: 0x46, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x12},
			wantState:  State{A: 0x34, P: FlagDecimal | FlagCarry},
		},
		{
			/*
				SED      ; Decimal mode (BCD subtraction: 40 - 13 = 27)
				SEC
				LDA #$40
				SBC #$13 ; After this instruction, C = 1, A = $27)
			*/
			name:       "Decimal mode (BCD subtraction: 40 - 13 = 27)",
			startState: State{A: 0x40, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x13},
			wantState:  State{A: 0x27, P: FlagDecimal | FlagCarry},
		},
		{
			/*
				SED      ; Decimal mode (BCD subtraction: 32 - 2 - 1 = 29)
				CLC      ; Note: carry is clear, not set!
				LDA #$32
				SBC #$02 ; After this instruction, C = 1, A = $29)
			*/
			name:       "Decimal mode (BCD subtraction: 32 - 2 - 1 = 29)",
			startState: State{A: 0x32, P: FlagDecimal},
			addressing: Addressing{Value: 0x02},
			wantState:  State{A: 0x29, P: FlagDecimal | FlagCarry},
		},
		{
			/*
				SED      ; Decimal mode (BCD subtraction: 12 - 21)
				SEC
				LDA #$12
				SBC #$21 ; After this instruction, C = 0, A = $91)
			*/
			name:       "Decimal mode (BCD subtraction: 12 - 21)",
			startState: State{A: 0x12, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x21},
			wantState:  State{A: 0x91, P: FlagDecimal | FlagNegative},
		},
		{
			/*
				SED      ; Decimal mode (BCD subtraction: 21 - 34)
				SEC
				LDA #$21
				SBC #$34 ; After this instruction, C = 0, A = $87)
			*/
			name:       "Decimal mode (BCD subtraction: 21 - 34)",
			startState: State{A: 0x21, P: FlagDecimal | FlagCarry},
			addressing: Addressing{Value: 0x34},
			wantState:  State{A: 0x87, P: FlagDecimal | FlagNegative},
		},
		// Other BCD tests.
		{
			name:       "BCD zero subtract zero is zero",
			startState: State{A: 0x00, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagNoBorrow | FlagZero},
		},
		{
			name:       "BCD 1 subtract zero is 1",
			startState: State{A: 0x01, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x01, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD 1 subtract 1 is zero",
			startState: State{A: 0x01, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagNoBorrow | FlagZero},
		},
		{
			name:       "BCD $0A subtract zero is $0A",
			startState: State{A: 0x0A, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0A, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $0A subtract 1 is $09",
			startState: State{A: 0x0A, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x09, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $0B subtract zero is $0B",
			startState: State{A: 0x0B, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0B, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $0B subtract 1 is $0A",
			startState: State{A: 0x0B, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x0A, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $1A subtract 1 is $19",
			startState: State{A: 0x1A, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x19, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $1B subtract 1 is $1A",
			startState: State{A: 0x1B, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x1A, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $20 subtract $0A is $10",
			startState: State{A: 0x20, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x0A},
			wantState:  State{A: 0x10, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $20 subtract $0B is $0F",
			startState: State{A: 0x20, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x0B},
			wantState:  State{A: 0x0F, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD zero subtract $0C is $8E (and borrows)",
			startState: State{A: 0x00, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x0C},
			wantState:  State{A: 0x8E, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD $0D subtract zero is $0D",
			startState: State{A: 0x0D, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0D, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD zero subtract $0E is $0E (and borrows)",
			startState: State{A: 0x00, P: FlagDecimal},
			addressing: Addressing{Value: 0x0E},
			wantState:  State{A: 0x8B, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD $0F subtract zero (with borrow) is $0F",
			startState: State{A: 0x0F, P: FlagDecimal},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x0E, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $01 subtract $0C (with borrow) is $13",
			startState: State{A: 0x01, P: FlagDecimal},
			addressing: Addressing{Value: 0x0C},
			wantState:  State{A: 0x8E, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD $10 subtract $05 is $05",
			startState: State{A: 0x10, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x05},
			wantState:  State{A: 0x05, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $16 subtract $08 (with borrow) is $07",
			startState: State{A: 0x16, P: FlagDecimal},
			addressing: Addressing{Value: 0x08},
			wantState:  State{A: 0x07, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $30 subtract $20 is $10",
			startState: State{A: 0x30, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x20},
			wantState:  State{A: 0x10, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $01 subtract $01 is $00",
			startState: State{A: 0x01, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagNoBorrow | FlagZero},
		},
		{
			name:       "BCD $01 subtract $00 with borrow is $00",
			startState: State{A: 0x01, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagNoBorrow | FlagZero},
		},
		{
			name:       "BCD $01 subtract $01 with borrow is $99 with borrow",
			startState: State{A: 0x01, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x99, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD $01 subtract $10 with borrow is $90",
			startState: State{A: 0x01, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x10},
			wantState:  State{A: 0x90, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD $11 subtract $13 is $98 with borrow",
			startState: State{A: 0x11, P: FlagDecimal | FlagNoBorrow},
			addressing: Addressing{Value: 0x13},
			wantState:  State{A: 0x98, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD $11 subtract $13 with borrow is $97 with borrow",
			startState: State{A: 0x11, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x13},
			wantState:  State{A: 0x97, P: FlagDecimal | FlagNone | FlagNegative},
		},
		// Test that the borrow flag works with BCD.
		{
			name:       "BCD zero subtract zero with borrow is 99 with borrow",
			startState: State{A: 0x00, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x99, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD 1 subtract zero with borrow is zero",
			startState: State{A: 0x01, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x00, P: FlagDecimal | FlagNoBorrow | FlagZero},
		},
		{
			name:       "BCD 1 subtract 1 with borrow is 99 with borrow",
			startState: State{A: 0x01, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x01},
			wantState:  State{A: 0x99, P: FlagDecimal | FlagNone | FlagNegative},
		},
		{
			name:       "BCD $06 subtract $02 with borrow is $03",
			startState: State{A: 0x06, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x02},
			wantState:  State{A: 0x03, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $10 subtract $05 with borrow is $04",
			startState: State{A: 0x10, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x05},
			wantState:  State{A: 0x04, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $40 subtract $20 with borrow is $19",
			startState: State{A: 0x40, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x20},
			wantState:  State{A: 0x19, P: FlagDecimal | FlagNoBorrow},
		},
		{
			name:       "BCD $99 subtract $00 with borrow is $98",
			startState: State{A: 0x99, P: FlagDecimal | FlagNone},
			addressing: Addressing{Value: 0x00},
			wantState:  State{A: 0x98, P: FlagDecimal | FlagNoBorrow | FlagNegative},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOperation(t, tt, SubtractWithCarry)
		})
	}
}

// Tests for all possible combinations of ADC and SBC:
//
// Explanation of the overflow flag: http://www.6502.org/tutorials/vflag.html
// Explanation of flags for BCD: http://www.6502.org/tutorials/decimal_mode.html
func TestAllAddsAndSubtracts(t *testing.T) {

	// Converts an 8-bit signed integer into a Go signed int for easy
	// manipulation.
	toSigned := func(value uint8) int {
		// Get the unsigned portion
		result := int(value & 0x7F)

		if (value & 0x80) != 0 {
			result -= 128
		}
		return result
	}

	calcZeroNegativeOverflowFlags := func(unsigned, signed int) (result Status) {

		if uint8(unsigned) == 0 {
			result.SetZero()
		}

		if uint8(unsigned&0x80) != 0 {
			result.SetNegative()
		}

		if signed < -128 || signed > 127 {
			result.SetOverflow()
		}

		return
	}

	validateAddBinary := func(t *testing.T, numOne, numTwo int, carry bool) {
		unsigned := numOne + numTwo
		if carry {
			unsigned++
		}

		signedOne := toSigned(uint8(numOne))
		signedTwo := toSigned(uint8(numTwo))
		signed := signedOne + signedTwo
		if carry {
			signed++
		}

		expectedFlags := calcZeroNegativeOverflowFlags(unsigned, signed)

		if unsigned > 0xFF {
			expectedFlags.SetCarry()
		}

		p := Status(FlagNone)
		if carry {
			p.SetCarry()
		}

		startState := State{A: uint8(numOne), P: p}
		addressing := Addressing{Value: uint8(numTwo)}

		result, err := AddWithCarry(startState, addressing)
		if err != nil {
			panic(err)
		}

		report := func() {
			fmt.Printf("ADD Binary Report\n")
			fmt.Printf("  Unsigned .. : U1: % 4d, U2: % 4d, Carry: %v\n", numOne, numTwo, carry)
			fmt.Printf("  Signed .... : S1: % 4d, S2: % 4d, Carry: %v\n", signedOne, signedTwo, carry)
			fmt.Printf("  Expected .. : U: % 4d, S: % 4d, Carry: %v, Flags: %v\n", unsigned, signed, carry, expectedFlags.ToFlags())
			fmt.Printf("  Actual .... : U: % 4d, S: % 4d, Carry: %v, Flags: %v\n", result.A, toSigned(result.A), carry, result.P.ToFlags())
		}

		if !reflect.DeepEqual(result.P, expectedFlags) {
			t.Errorf("ADD Binary Flags got  = %v, want = %v", result.P.ToFlags(), expectedFlags.ToFlags())
			report()
		}
		if !reflect.DeepEqual(result.A, uint8(unsigned)) {
			t.Errorf("ADD Binary Accumulator got = %v, want = %v", result.A, unsigned)
			report()
		}
	}

	validateSubBinary := func(t *testing.T, numOne, numTwo int, borrow bool) {
		unsigned := numOne - numTwo
		if borrow {
			unsigned--
		}

		signedOne := toSigned(uint8(numOne))
		signedTwo := toSigned(uint8(numTwo))
		signed := signedOne - signedTwo
		if borrow {
			signed--
		}

		expectedFlags := calcZeroNegativeOverflowFlags(unsigned, signed)

		if unsigned >= 0x00 { // i.e. no borrow
			expectedFlags.SetCarry()
		}

		p := Status(FlagCarry)
		if borrow {
			p.ClearCarry()
		}

		startState := State{A: uint8(numOne), P: p}
		addressing := Addressing{Value: uint8(numTwo)}

		result, err := SubtractWithCarry(startState, addressing)
		if err != nil {
			panic(err)
		}

		report := func() {
			fmt.Printf("SUB Binary Report\n")
			fmt.Printf("  Unsigned .. : U1: % 4d, U2: % 4d, Borrow: %v\n", numOne, numTwo, borrow)
			fmt.Printf("  Signed .... : S1: % 4d, S2: % 4d, Borrow: %v\n", signedOne, signedTwo, borrow)
			fmt.Printf("  Expected .. : U: % 4d, S: % 4d, Borrow: %v, Flags: %v\n", unsigned, signed, borrow, expectedFlags.ToFlags())
			fmt.Printf("  Actual .... : U: % 4d, S: % 4d, Borrow: %v, Flags: %v\n", result.A, toSigned(result.A), borrow, result.P.ToFlags())
		}

		if !reflect.DeepEqual(result.P, expectedFlags) {
			t.Errorf("SUB Binary Flags got  = %v, want = %v", result.P.ToFlags(), expectedFlags.ToFlags())
			report()
		}
		if !reflect.DeepEqual(result.A, uint8(unsigned)) {
			t.Errorf("SUB Binary Accumulator got = %v, want = %v", result.A, unsigned)
			report()
		}
	}

	validateAddDecimal := func(t *testing.T, numOne, numTwo int, carry bool) {

		decimalOne := DecodeBcdValue(uint8(numOne))
		decimalTwo := DecodeBcdValue(uint8(numTwo))

		unsigned := int(decimalOne) + int(decimalTwo)
		if carry {
			unsigned++
		}
		expectedBcd := EncodeBcdValue(uint8(unsigned & 0xFF))

		expectedFlags := calcZeroNegativeOverflowFlags(int(expectedBcd), 0)
		if unsigned > 99 {
			expectedFlags.SetCarry()
		}
		expectedFlags.SetDecimal()

		p := Status(FlagDecimal)
		if carry {
			p.SetCarry()
		}

		startState := State{A: uint8(numOne), P: p}
		addressing := Addressing{Value: uint8(numTwo)}

		result, err := AddWithCarry(startState, addressing)
		if err != nil {
			panic(err)
		}

		report := func() {
			fmt.Printf("ADD Decimal Report\n")
			fmt.Printf("  Unsigned .. : B1: 0x%02X, B2: 0x%02X, Carry: %v\n", numOne, numTwo, carry)
			fmt.Printf("  Decimal ... : D1: % 4d, D2: % 4d, Carry: %v\n", decimalOne, decimalTwo, carry)
			fmt.Printf("  Expected .. : B: 0x%02X, Carry: %v, Flags: %v, D: % 4d\n", expectedBcd, carry, expectedFlags.ToFlags(), unsigned)
			fmt.Printf("  Actual .... : B: 0x%02X, Carry: %v, Flags: %v\n", result.A, carry, result.P.ToFlags())
		}

		// On the 6502, only the C flag is valid in decimal mode. On the 65C02 and 65816, the Z flag is valid.
		// In our simulator, the C, N amd Z flags are valid but V is not.
		expectedFlags.ClearOverflow()
		result.P.ClearOverflow()

		if !reflect.DeepEqual(result.P, expectedFlags) {
			t.Errorf("ADD Decimal Flags got  = %v, want = %v", result.P.ToFlags(), expectedFlags.ToFlags())
			report()
		}
		if !reflect.DeepEqual(result.A, expectedBcd) {
			t.Errorf("ADD Decimal Accumulator got = %v, want = %v", result.A, expectedBcd)
			report()
		}
	}

	validateSubDecimal := func(t *testing.T, numOne, numTwo int, borrow bool) {

		decimalOne := DecodeBcdValue(uint8(numOne))
		decimalTwo := DecodeBcdValue(uint8(numTwo))

		unsigned := int(decimalOne) - int(decimalTwo)
		if borrow {
			unsigned--
		}

		// if we've borrowed (i.e. less than zero), we need to add the borrowed 100 before encoding.
		modified := unsigned
		if unsigned < 0 {
			modified += 100
		}
		expectedBcd := EncodeBcdValue(uint8(modified & 0xFF))

		expectedFlags := calcZeroNegativeOverflowFlags(int(expectedBcd), 0)
		if unsigned >= 0 { // i.e. no borrow
			expectedFlags.SetCarry()
		}
		expectedFlags.SetDecimal()

		p := Status(FlagDecimal | FlagCarry)
		if borrow {
			p.ClearCarry()
		}

		startState := State{A: uint8(numOne), P: p}
		addressing := Addressing{Value: uint8(numTwo)}

		result, err := SubtractWithCarry(startState, addressing)
		if err != nil {
			panic(err)
		}

		report := func() {
			fmt.Printf("SUB Decimal Report\n")
			fmt.Printf("  Unsigned .. : B1: 0x%02X, B2: 0x%02X, Borrow: %v\n", numOne, numTwo, borrow)
			fmt.Printf("  Decimal ... : D1: % 4d, D2: % 4d, Borrow: %v\n", decimalOne, decimalTwo, borrow)
			fmt.Printf("  Expected .. : B: 0x%02X, Borrow: %v, Flags: %v, D: % 4d\n", expectedBcd, borrow, expectedFlags.ToFlags(), unsigned)
			fmt.Printf("  Actual .... : B: 0x%02X, Borrow: %v, Flags: %v\n", result.A, borrow, result.P.ToFlags())
		}

		// On the 6502, only the C flag is valid in decimal mode. On the 65C02 and 65816, the Z flag is valid.
		// In our simulator, the C, N amd Z flags are valid but V is not.
		expectedFlags.ClearOverflow()
		result.P.ClearOverflow()

		if !reflect.DeepEqual(result.P, expectedFlags) {
			t.Errorf("SUB Decimal Flags got  = %v, want = %v", result.P.ToFlags(), expectedFlags.ToFlags())
			report()
		}
		if !reflect.DeepEqual(result.A, expectedBcd) {
			t.Errorf("SUB Decimal Accumulator got = %v, want = %v", result.A, expectedBcd)
			report()
		}
	}

	for _, carry := range [...]bool{true, false} {
		for numOne := range 256 {
			for numTwo := range 256 {

				borrow := !carry

				validateAddBinary(t, numOne, numTwo, carry)
				validateSubBinary(t, numOne, numTwo, borrow)

				// Convert the input numbers to their decimal equivalents and only consider valid BCD values.
				decimalOne := DecodeBcdValue(uint8(numOne))
				decimalTwo := DecodeBcdValue(uint8(numTwo))

				if decimalOne > 100 || decimalTwo > 100 {
					continue
				}

				validateAddDecimal(t, numOne, numTwo, carry)
				validateSubDecimal(t, numOne, numTwo, borrow)
			}
		}
	}
}
