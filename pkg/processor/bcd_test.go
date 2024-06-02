package processor

import "testing"

func TestExtractBcdValues(t *testing.T) {

	tests := []struct {
		name  string
		value uint8
		wanth uint8
		wantl uint8
	}{
		{name: "Zero returns zero", value: 0x00, wanth: 0x0, wantl: 0x0},
		{name: "First digit only - 1", value: 0x01, wanth: 0x0, wantl: 0x1},
		{name: "First digit only - 3", value: 0x03, wanth: 0x0, wantl: 0x3},
		{name: "First digit only - 9", value: 0x09, wanth: 0x0, wantl: 0x9},
		{name: "First digit only - A", value: 0x0A, wanth: 0x0, wantl: 0xA},
		{name: "First digit only - F", value: 0x0F, wanth: 0x0, wantl: 0xF},
		{name: "Second digit only - 1", value: 0x10, wanth: 0x1, wantl: 0x0},
		{name: "Second digit only - 3", value: 0x30, wanth: 0x3, wantl: 0x0},
		{name: "Second digit only - 9", value: 0x90, wanth: 0x9, wantl: 0x0},
		{name: "Second digit only - A", value: 0xA0, wanth: 0xA, wantl: 0x0},
		{name: "Second digit only - F", value: 0xF0, wanth: 0xF, wantl: 0x0},
		{name: "Both digits - 1F", value: 0x1F, wanth: 0x1, wantl: 0xF},
		{name: "Both digits - 3A", value: 0x3A, wanth: 0x3, wantl: 0xA},
		{name: "Both digits - 97", value: 0x97, wanth: 0x9, wantl: 0x7},
		{name: "Both digits - A3", value: 0xA3, wanth: 0xA, wantl: 0x3},
		{name: "Both digits - F1", value: 0xF1, wanth: 0xF, wantl: 0x1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goth, gotl := ExtractBcdValues(tt.value)
			if goth != tt.wanth {
				t.Errorf("ExtractBcdValues() got = %v, want %v", goth, tt.wanth)
			}
			if gotl != tt.wantl {
				t.Errorf("ExtractBcdValues() got1 = %v, want %v", gotl, tt.wantl)
			}
		})
	}
}

func TestDecodeBcdValue(t *testing.T) {

	tests := []struct {
		name  string
		value uint8
		want  uint8
	}{
		{name: "Zero returns zero", value: 0x00, want: 0},
		{name: "First digit only - 1", value: 0x01, want: 1},
		{name: "First digit only - 3", value: 0x03, want: 3},
		{name: "First digit only - 9", value: 0x09, want: 9},
		{name: "First digit only - A", value: 0x0A, want: 0xFF},
		{name: "First digit only - F", value: 0x0F, want: 0xFF},
		{name: "Second digit only - 1", value: 0x10, want: 10},
		{name: "Second digit only - 3", value: 0x30, want: 30},
		{name: "Second digit only - 9", value: 0x90, want: 90},
		{name: "Second digit only - A", value: 0xA0, want: 0xFF},
		{name: "Second digit only - F", value: 0xF0, want: 0xFF},
		{name: "Both digits - 1F", value: 0x1F, want: 0xFF},
		{name: "Both digits - 3A", value: 0x3A, want: 0xFF},
		{name: "Both digits - 13", value: 0x13, want: 13},
		{name: "Both digits - 47", value: 0x47, want: 47},
		{name: "Both digits - 99", value: 0x99, want: 99},
		{name: "Both digits - A3", value: 0xA3, want: 0xFF},
		{name: "Both digits - F1", value: 0xF1, want: 0xFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeBcdValue(tt.value); got != tt.want {
				t.Errorf("DecodeBcdValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBcdValue(t *testing.T) {
	tests := []struct {
		name  string
		value uint8
		want  uint8
	}{
		{name: "Zero returns zero", value: 0, want: 0x00},
		{name: "First digit only - 1", value: 1, want: 0x01},
		{name: "First digit only - 3", value: 3, want: 0x03},
		{name: "First digit only - 9", value: 9, want: 0x09},
		{name: "Second digit only - 1", value: 10, want: 0x10},
		{name: "Second digit only - 3", value: 30, want: 0x30},
		{name: "Second digit only - 9", value: 90, want: 0x90},
		{name: "Both digits - 13", value: 13, want: 0x13},
		{name: "Both digits - 47", value: 47, want: 0x47},
		{name: "Both digits - 99", value: 99, want: 0x99},
		{name: "Out of range - 100", value: 100, want: 0x00},
		{name: "Out of range - 102", value: 102, want: 0x02},
		{name: "Out of range - 123", value: 123, want: 0x23},
		{name: "Out of range - 200", value: 200, want: 0x00},
		{name: "Out of range - 201", value: 201, want: 0x01},
		{name: "Out of range - 231", value: 231, want: 0x31},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeBcdValue(tt.value); got != tt.want {
				t.Errorf("EncodeBcdValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
