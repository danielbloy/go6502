package processor

import "testing"

func TestStatus_WithFlagsSetOrCleared(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		fn     func(s *Status) Status
		want   string
	}{
		{
			name:   "Set constant flag",
			status: 0,
			fn:     (*Status).WithConstantSet,
			want:   "n v - b d i z c",
		},
		{
			name:   "Clear constant flag",
			status: 0,
			fn:     (*Status).WithConstantCleared,
			want:   "n v - b d i z c",
		},
		{
			name:   "Set carry flag",
			status: 0,
			fn:     (*Status).WithCarrySet,
			want:   "n v - b d i z C",
		},
		{
			name:   "Clear carry flag",
			status: 0xFF,
			fn:     (*Status).WithCarryCleared,
			want:   "N V - B D I Z c",
		},
		{
			name:   "Set zero flag",
			status: 0,
			fn:     (*Status).WithZeroSet,
			want:   "n v - b d i Z c",
		},
		{
			name:   "Clear zero flag",
			status: 0xFF,
			fn:     (*Status).WithZeroCleared,
			want:   "N V - B D I z C",
		},
		{
			name:   "Set interrupt flag",
			status: 0,
			fn:     (*Status).WithInterruptSet,
			want:   "n v - b d I z c",
		},
		{
			name:   "Clear interrupt flag",
			status: 0xFF,
			fn:     (*Status).WithInterruptCleared,
			want:   "N V - B D i Z C",
		},
		{
			name:   "Set decimal flag",
			status: 0,
			fn:     (*Status).WithDecimalSet,
			want:   "n v - b D i z c",
		},
		{
			name:   "Clear decimal flag",
			status: 0xFF,
			fn:     (*Status).WithDecimalCleared,
			want:   "N V - B d I Z C",
		},
		{
			name:   "Set break flag",
			status: 0,
			fn:     (*Status).WithBreakSet,
			want:   "n v - B d i z c",
		},
		{
			name:   "Clear break flag",
			status: 0xFF,
			fn:     (*Status).WithBreakCleared,
			want:   "N V - b D I Z C",
		},
		{
			name:   "Set overflow flag",
			status: 0,
			fn:     (*Status).WithOverflowSet,
			want:   "n V - b d i z c",
		},
		{
			name:   "Clear overflow flag",
			status: 0xFF,
			fn:     (*Status).WithOverflowCleared,
			want:   "N v - B D I Z C",
		},
		{
			name:   "Set negative flag",
			status: 0,
			fn:     (*Status).WithNegativeSet,
			want:   "N v - b d i z c",
		},
		{
			name:   "Clear negative flag",
			status: 0xFF,
			fn:     (*Status).WithNegativeCleared,
			want:   "n V - B D I Z C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startValue := tt.status
			modifiedValue := tt.fn(&startValue)

			// Make sure the original was NOT changed but the returned value was
			if got := modifiedValue.String(); got != tt.want {
				t.Errorf("SetClearFlags(); got = %v, want %v", got, tt.want)
			}
			if startValue != tt.status {
				t.Errorf("SetClearFlags() changed the original value; got = %v, want %v", startValue, tt.status)
			}

			// Test with nil
			if tt.fn(nil) != 0 {
				t.Errorf("SetClearFlags() did not return zero when called on nil")
			}
		})
	}
}

func TestStatus_SetClearFlags(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		fn     func(s *Status) *Status
		want   string
	}{
		{
			name:   "Set carry flag",
			status: 0,
			fn:     (*Status).SetCarry,
			want:   "n v - b d i z C",
		},
		{
			name:   "Clear carry flag",
			status: 0xFF,
			fn:     (*Status).ClearCarry,
			want:   "N V - B D I Z c",
		},
		{
			name:   "Set zero flag",
			status: 0,
			fn:     (*Status).SetZero,
			want:   "n v - b d i Z c",
		},
		{
			name:   "Clear zero flag",
			status: 0xFF,
			fn:     (*Status).ClearZero,
			want:   "N V - B D I z C",
		},
		{
			name:   "Set interrupt flag",
			status: 0,
			fn:     (*Status).SetInterrupt,
			want:   "n v - b d I z c",
		},
		{
			name:   "Clear interrupt flag",
			status: 0xFF,
			fn:     (*Status).ClearInterrupt,
			want:   "N V - B D i Z C",
		},
		{
			name:   "Set decimal flag",
			status: 0,
			fn:     (*Status).SetDecimal,
			want:   "n v - b D i z c",
		},
		{
			name:   "Clear decimal flag",
			status: 0xFF,
			fn:     (*Status).ClearDecimal,
			want:   "N V - B d I Z C",
		},
		{
			name:   "Set break flag",
			status: 0,
			fn:     (*Status).SetBreak,
			want:   "n v - B d i z c",
		},
		{
			name:   "Clear break flag",
			status: 0xFF,
			fn:     (*Status).ClearBreak,
			want:   "N V - b D I Z C",
		},
		{
			name:   "Set overflow flag",
			status: 0,
			fn:     (*Status).SetOverflow,
			want:   "n V - b d i z c",
		},
		{
			name:   "Clear overflow flag",
			status: 0xFF,
			fn:     (*Status).ClearOverflow,
			want:   "N v - B D I Z C",
		},
		{
			name:   "Set negative flag",
			status: 0,
			fn:     (*Status).SetNegative,
			want:   "N v - b d i z c",
		},
		{
			name:   "Clear negative flag",
			status: 0xFF,
			fn:     (*Status).ClearNegative,
			want:   "n V - B D I Z C",
		},
		{
			name:   "Set constant flag",
			status: 0,
			fn:     (*Status).SetConstant,
			want:   "n v - b d i z c",
		},
		{
			name:   "Clear constant flag",
			status: 0xFF,
			fn:     (*Status).ClearConstant,
			want:   "N V - B D I Z C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startValue := tt.status
			modifiedValue := tt.fn(&startValue)

			// Make sure the original was changed and the modified value points to the start value
			if got := startValue.String(); got != tt.want {
				t.Errorf("SetClearFlags(); got = %v, want %v", got, tt.want)
			}
			if &startValue != modifiedValue {
				t.Errorf("SetClearFlags() did not return the correct address; got = %v, want %v", &startValue, modifiedValue)
			}

			// Test with nil
			if tt.fn(nil) != nil {
				t.Errorf("SetClearFlags() did not return nil when called on nil")
			}
		})
	}
}

func TestStatus_ToFlags(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   Flags
	}{
		{
			name:   "Status with no flags set",
			status: 0,
			want:   Flags{},
		},
		{
			name:   "Status with carry flag set",
			status: FlagCarry,
			want:   Flags{Carry: true},
		},
		{
			name:   "Status with zero flag set",
			status: FlagZero,
			want:   Flags{Zero: true},
		},
		{
			name:   "Status with interrupt flag set",
			status: FlagInterrupt,
			want:   Flags{Interrupt: true},
		},
		{
			name:   "Status with decimal flag set",
			status: FlagDecimal,
			want:   Flags{Decimal: true},
		},
		{
			name:   "Status with Break flag set",
			status: FlagBreak,
			want:   Flags{Break: true},
		},
		{
			name:   "Status with Overflow flag set",
			status: FlagOverflow,
			want:   Flags{Overflow: true},
		},
		{
			name:   "Status with negative flag set",
			status: FlagNegative,
			want:   Flags{Negative: true},
		},
		{
			name:   "Status with carry, interrupt, break and negative flags set",
			status: FlagCarry | FlagInterrupt | FlagBreak | FlagNegative,
			want:   Flags{Carry: true, Interrupt: true, Break: true, Negative: true},
		},
		{
			name:   "Status with all flags set",
			status: 0xFF,
			want: Flags{
				Carry:     true,
				Zero:      true,
				Interrupt: true,
				Decimal:   true,
				Break:     true,
				Overflow:  true,
				Negative:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.ToFlags(); got != tt.want {
				t.Errorf("ToFlags(); got = %v, want %v", got, tt.want)
			}
		})
	}

	// Test with nil
	{
		fn := (*Status).ToFlags
		zero := Flags{}
		if got := fn(nil); got != zero {
			t.Errorf("ToFlags() did not return zero Flags when called with nil")
		}
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   string
	}{
		{
			name:   "Status with no flags set",
			status: 0,
			want:   "n v - b d i z c",
		},
		{
			name:   "Status with carry flag set",
			status: FlagCarry,
			want:   "n v - b d i z C",
		},
		{
			name:   "Status with zero flag set",
			status: FlagZero,
			want:   "n v - b d i Z c",
		},
		{
			name:   "Status with interrupt flag set",
			status: FlagInterrupt,
			want:   "n v - b d I z c",
		},
		{
			name:   "Status with decimal flag set",
			status: FlagDecimal,
			want:   "n v - b D i z c",
		},
		{
			name:   "Status with Break flag set",
			status: FlagBreak,
			want:   "n v - B d i z c",
		},
		{
			name:   "Status with Overflow flag set",
			status: FlagOverflow,
			want:   "n V - b d i z c",
		},
		{
			name:   "Status with negative flag set",
			status: FlagNegative,
			want:   "N v - b d i z c",
		},
		{
			name:   "Status with carry, interrupt, break and negative flags set",
			status: FlagCarry | FlagInterrupt | FlagBreak | FlagNegative,
			want:   "N v - B d I z C",
		},
		{
			name:   "Status with all flags set",
			status: 0xFF,
			want:   "N V - B D I Z C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("String(); got = %v, want %v", got, tt.want)
			}
		})
	}

	// Test with nil
	{
		fn := (*Status).String
		if got := fn(nil); got != "" {
			t.Errorf("String() did not return empty string when called with nil")
		}
	}
}

func TestFlags_ToStatus(t *testing.T) {
	tests := []struct {
		name  string
		flags Flags
		want  Status
	}{
		{
			name:  "Status with no flags set",
			flags: Flags{},
			want:  FlagConstant,
		},
		{
			name:  "Status with carry flag set",
			flags: Flags{Carry: true},
			want:  FlagConstant | FlagCarry,
		},
		{
			name:  "Status with zero flag set",
			flags: Flags{Zero: true},
			want:  FlagConstant | FlagZero,
		},
		{
			name:  "Status with interrupt flag set",
			flags: Flags{Interrupt: true},
			want:  FlagConstant | FlagInterrupt,
		},
		{
			name:  "Status with decimal flag set",
			flags: Flags{Decimal: true},
			want:  FlagConstant | FlagDecimal,
		},
		{
			name:  "Status with Break flag set",
			flags: Flags{Break: true},
			want:  FlagConstant | FlagBreak,
		},
		{
			name:  "Status with Overflow flag set",
			flags: Flags{Overflow: true},
			want:  FlagConstant | FlagOverflow,
		},
		{
			name:  "Status with negative flag set",
			flags: Flags{Negative: true},
			want:  FlagConstant | FlagNegative,
		},
		{
			name:  "Status with carry, interrupt, break and negative flags set",
			flags: Flags{Carry: true, Interrupt: true, Break: true, Negative: true},
			want:  FlagConstant | FlagCarry | FlagInterrupt | FlagBreak | FlagNegative,
		},
		{
			name: "Status with all flags set",
			flags: Flags{
				Carry:     true,
				Zero:      true,
				Interrupt: true,
				Decimal:   true,
				Break:     true,
				Overflow:  true,
				Negative:  true,
			},
			want: 0xFF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.ToStatus(); got != tt.want {
				t.Errorf("ToStatus(); got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlags_String(t *testing.T) {
	tests := []struct {
		name  string
		flags Flags
		want  string
	}{
		{
			name:  "Status with no flags set",
			flags: Flags{},
			want:  "n v - b d i z c",
		},
		{
			name:  "Status with carry flag set",
			flags: Flags{Carry: true},
			want:  "n v - b d i z C",
		},
		{
			name:  "Status with zero flag set",
			flags: Flags{Zero: true},
			want:  "n v - b d i Z c",
		},
		{
			name:  "Status with interrupt flag set",
			flags: Flags{Interrupt: true},
			want:  "n v - b d I z c",
		},
		{
			name:  "Status with decimal flag set",
			flags: Flags{Decimal: true},
			want:  "n v - b D i z c",
		},
		{
			name:  "Status with Break flag set",
			flags: Flags{Break: true},
			want:  "n v - B d i z c",
		},
		{
			name:  "Status with Overflow flag set",
			flags: Flags{Overflow: true},
			want:  "n V - b d i z c",
		},
		{
			name:  "Status with negative flag set",
			flags: Flags{Negative: true},
			want:  "N v - b d i z c",
		},
		{
			name:  "Status with carry, interrupt, break and negative flags set",
			flags: Flags{Carry: true, Interrupt: true, Break: true, Negative: true},
			want:  "N v - B d I z C",
		},
		{
			name: "Status with all flags set",
			flags: Flags{
				Carry:     true,
				Zero:      true,
				Interrupt: true,
				Decimal:   true,
				Break:     true,
				Overflow:  true,
				Negative:  true,
			},
			want: "N V - B D I Z C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.String(); got != tt.want {
				t.Errorf("String(); got = %v, want %v", got, tt.want)
			}
		})
	}
}
