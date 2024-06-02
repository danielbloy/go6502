package processor

import (
	"reflect"
	"testing"
)

func TestWriteVectorToMemory(t *testing.T) {
	tests := []struct {
		name    string
		address Address
		vector  Address
		want    []uint8
		wantErr bool
	}{
		{
			name:    "Write 0x1234 to 0x0000",
			address: 0,
			vector:  0x1234,
			want:    []uint8{0x34, 0x12, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 0xABCD to 0x3000",
			address: 3,
			vector:  0xABCD,
			want:    []uint8{0, 0, 0, 0xCD, 0xAB, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, nil)
			if err := WriteVectorToMemory(&memory, tt.address, tt.vector); (err != nil) != tt.wantErr {
				t.Errorf("WriteVectorToMemory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(memory.ram, tt.want) {
				t.Errorf("WriteVectorToMemory() got = %v, want = %v", memory.ram, tt.want)
			}
		})
	}

	// Check support for nil
	if err := WriteVectorToMemory(nil, 0, 0); err == nil {
		t.Errorf("WriteVectorToMemory() did not error with nil memory")
	}
}

func TestReadVectorFromMemory(t *testing.T) {
	tests := []struct {
		name    string
		ram     []uint8
		address Address
		want    Address
		wantErr bool
	}{
		{
			name:    "Read 0x1234 from 0x0001",
			ram:     []uint8{0, 0x34, 0x12, 0, 0, 0, 0, 0},
			address: 1,
			want:    0x1234,
		},
		{
			name:    "Read 0xABCD from 0x0005",
			ram:     []uint8{0, 0, 0, 0, 0, 0xCD, 0xAB, 0, 0},
			address: 5,
			want:    0xABCD,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, tt.ram)
			got, err := ReadVectorFromMemory(&memory, tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVectorFromMemory() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadVectorFromMemory() got = %v, want = %v", got, tt.want)
			}
		})
	}
	// Check support for nil
	if _, err := ReadVectorFromMemory(nil, 0); err == nil {
		t.Errorf("ReadVectorFromMemory() did not error with nil memory")
	}
}

func TestWriteNmiVectorToMemory(t *testing.T) {
	tests := []struct {
		name    string
		vector  Address
		want    []uint8
		wantErr bool
	}{
		{
			name:   "Write 0x1234",
			vector: 0x1234,
			want:   []uint8{0, 0, 0x34, 0x12, 0, 0, 0, 0},
		},
		{
			name:   "Write 0xABCD",
			vector: 0xABCD,
			want:   []uint8{0, 0, 0xCD, 0xAB, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, nil)
			if err := WriteNmiVectorToMemory(&memory, tt.vector); (err != nil) != tt.wantErr {
				t.Errorf("WriteNmiVectorToMemory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(memory.ram, tt.want) {
				t.Errorf("WriteNmiVectorToMemory() got = %v, want = %v", memory.ram, tt.want)
			}
		})
	}

	// Check support for nil
	if err := WriteNmiVectorToMemory(nil, 0); err == nil {
		t.Errorf("WriteNmiVectorToMemory() did not error with nil memory")
	}
}

func TestReadNmiVectorFromMemory(t *testing.T) {
	tests := []struct {
		name    string
		ram     []uint8
		want    Address
		wantErr bool
	}{
		{
			name: "Read 0x1234",
			ram:  []uint8{0, 0, 0x34, 0x12, 0, 0, 0, 0},
			want: 0x1234,
		},
		{
			name: "Read 0xABCD",
			ram:  []uint8{0, 0, 0xCD, 0xAB, 0, 0, 0, 0},
			want: 0xABCD,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, tt.ram)
			got, err := ReadNmiVectorFromMemory(&memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNmiVectorFromMemory() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadNmiVectorFromMemory() got = %v, want = %v", got, tt.want)
			}
		})
	}
	// Check support for nil
	if _, err := ReadNmiVectorFromMemory(nil); err == nil {
		t.Errorf("ReadNmiVectorFromMemory() did not error with nil memory")
	}
}

func TestWriteResetVectorToMemory(t *testing.T) {
	tests := []struct {
		name    string
		vector  Address
		want    []uint8
		wantErr bool
	}{
		{
			name:   "Write 0x1234",
			vector: 0x1234,
			want:   []uint8{0, 0, 0, 0, 0x34, 0x12, 0, 0},
		},
		{
			name:   "Write 0xABCD",
			vector: 0xABCD,
			want:   []uint8{0, 0, 0, 0, 0xCD, 0xAB, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, nil)
			if err := WriteResetVectorToMemory(&memory, tt.vector); (err != nil) != tt.wantErr {
				t.Errorf("WriteResetVectorToMemory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(memory.ram, tt.want) {
				t.Errorf("WriteResetVectorToMemory() got = %v, want = %v", memory.ram, tt.want)
			}
		})
	}

	// Check support for nil
	if err := WriteResetVectorToMemory(nil, 0); err == nil {
		t.Errorf("WriteResetVectorToMemory() did not error with nil memory")
	}
}

func TestReadResetVectorFromMemory(t *testing.T) {
	tests := []struct {
		name    string
		ram     []uint8
		want    Address
		wantErr bool
	}{
		{
			name: "Read 0x1234",
			ram:  []uint8{0, 0, 0, 0, 0x34, 0x12, 0, 0},
			want: 0x1234,
		},
		{
			name: "Read 0xABCD",
			ram:  []uint8{0, 0, 0, 0, 0xCD, 0xAB, 0, 0},
			want: 0xABCD,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, tt.ram)
			got, err := ReadResetVectorFromMemory(&memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadResetVectorFromMemory() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadResetVectorFromMemory() got = %v, want = %v", got, tt.want)
			}
		})
	}
	// Check support for nil
	if _, err := ReadResetVectorFromMemory(nil); err == nil {
		t.Errorf("ReadResetVectorFromMemory() did not error with nil memory")
	}
}

func TestWriteIrqVectorToMemory(t *testing.T) {
	tests := []struct {
		name    string
		vector  Address
		want    []uint8
		wantErr bool
	}{
		{
			name:   "Write 0x1234",
			vector: 0x1234,
			want:   []uint8{0, 0, 0, 0, 0, 0, 0x34, 0x12},
		},
		{
			name:   "Write 0xABCD",
			vector: 0xABCD,
			want:   []uint8{0, 0, 0, 0, 0, 0, 0xCD, 0xAB},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, nil)
			if err := WriteIrqVectorToMemory(&memory, tt.vector); (err != nil) != tt.wantErr {
				t.Errorf("WriteIrqVectorToMemory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(memory.ram, tt.want) {
				t.Errorf("WriteIrqVectorToMemory() got = %v, want = %v", memory.ram, tt.want)
			}
		})
	}

	// Check support for nil
	if err := WriteIrqVectorToMemory(nil, 0); err == nil {
		t.Errorf("WriteIrqVectorToMemory() did not error with nil memory")
	}
}

func TestReadIrqVectorFromMemory(t *testing.T) {
	tests := []struct {
		name    string
		ram     []uint8
		want    Address
		wantErr bool
	}{
		{
			name: "Read 0x1234",
			ram:  []uint8{0, 0, 0, 0, 0, 0, 0x34, 0x12},
			want: 0x1234,
		},
		{
			name: "Read 0xABCD",
			ram:  []uint8{0, 0, 0, 0, 0, 0, 0xCD, 0xAB},
			want: 0xABCD,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, tt.ram)
			got, err := ReadIrqVectorFromMemory(&memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNmiVectorFromMemory() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadIrqVectorFromMemory() got = %v, want = %v", got, tt.want)
			}
		})
	}
	// Check support for nil
	if _, err := ReadIrqVectorFromMemory(nil); err == nil {
		t.Errorf("ReadIrqVectorFromMemory() did not error with nil memory")
	}
}

func TestWriteContiguousDataToMemory(t *testing.T) {
	tests := []struct {
		name    string
		start   Address
		data    []uint8
		want    []uint8
		wantErr bool
	}{
		{
			name: "Write no data bytes at zero",
			want: []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "Write empty data bytes at zero",
			data: []uint8{},
			want: []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:  "Write one byte at 0x0000",
			start: 0,
			data:  []uint8{0xA},
			want:  []uint8{0xA, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:  "Write two bytes at 0x0000",
			start: 0,
			data:  []uint8{0xA, 0xB},
			want:  []uint8{0xA, 0xB, 0, 0, 0, 0, 0, 0},
		},
		{
			name:  "Write one byte at 0x0001",
			start: 0x0001,
			data:  []uint8{0xF},
			want:  []uint8{0, 0xF, 0, 0, 0, 0, 0, 0},
		},
		{
			name:  "Write two bytes at 0x0003",
			start: 0x0003,
			data:  []uint8{0xA, 0xB},
			want:  []uint8{0, 0, 0, 0xA, 0xB, 0, 0, 0},
		},
		{
			name:  "Write five bytes at 0x0005",
			start: 0x0005,
			data:  []uint8{5, 4, 3, 2, 1},
			want:  []uint8{2, 1, 0, 0, 0, 5, 4, 3},
		},
		{
			name:  "Write two bytes at 0x0006",
			start: 0x0006,
			data:  []uint8{1, 2},
			want:  []uint8{0, 0, 0, 0, 0, 0, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewPopulatedRam(EightBytes, nil)
			if err := WriteContiguousDataToMemory(&memory, tt.start, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("WriteContiguousDataToMemory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(memory.ram, tt.want) {
				t.Errorf("WriteContiguousDataToMemory() got = %v, want = %v", memory.ram, tt.want)
			}
		})
	}
	// Check support for nil
	if err := WriteContiguousDataToMemory(nil, 0, nil); err == nil {
		t.Errorf("WriteContiguousDataToMemory() did not error with nil memory")
	}

	ram, err := NewRepeatingRam(EightBytes)
	if err != nil {
		panic(err)
	}
	if err = WriteContiguousDataToMemory(&ram, 0, nil); err != nil {
		t.Errorf("WriteContiguousDataToMemory() errored with nil data")
	}
}

func TestWriteDataToMemory(t *testing.T) {
	tests := []struct {
		name    string
		data    map[Address]uint8
		want    []uint8
		wantErr bool
	}{
		{
			name: "Write no data bytes at zero",
			want: []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "Write empty data bytes at zero",
			data: map[Address]uint8{},
			want: []uint8{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "Write oon byte at 0x0000",
			data: map[Address]uint8{0x0: 0xA},
			want: []uint8{0xA, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "Write two bytes at different locations",
			data: map[Address]uint8{0x0: 0xA, 0x01: 0xB},
			want: []uint8{0xA, 0xB, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "Write four bytes at different locations",
			data: map[Address]uint8{0x5: 0xA, 0x01: 0xB, 0xFFFF: 0x3, 0x03: 0xFF},
			want: []uint8{0, 0xB, 0, 0xFF, 0, 0xA, 0, 0x3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ram, err := NewRepeatingRam(EightBytes)
			if err != nil {
				panic(err)
			}
			if err = WriteDataToMemory(&ram, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("WriteDataToMemory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(ram.ram, tt.want) {
				t.Errorf("WriteDataToMemory() got = %v, want = %v", ram.ram, tt.want)
			}
		})
	}
	// Check support for nil
	if err := WriteDataToMemory(nil, nil); err == nil {
		t.Errorf("WriteDataToMemory() did not error with nil memory")
	}

	ram, err := NewRepeatingRam(EightBytes)
	if err != nil {
		panic(err)
	}
	if err = WriteDataToMemory(&ram, nil); err != nil {
		t.Errorf("WriteDataToMemory() errored with nil data")
	}
}

func TestRepeatingRam_Read(t *testing.T) {
	emptyRam := make([]uint8, EightBytes)
	populatedRam := []uint8{1, 2, 3, 4, 5, 6, 7, 8}

	tests := []struct {
		name    string
		ram     []uint8
		address Address
		want    uint8
	}{
		{
			name:    "Empty ram returns zero at address 0x00",
			ram:     emptyRam,
			address: 0x00,
		},
		{
			name:    "Empty ram returns zero at address 0xFFFF",
			ram:     emptyRam,
			address: 0xFFFF,
		},
		{
			name:    "Populated ram returns 1 at address 0x00",
			ram:     populatedRam,
			address: 0x00,
			want:    1,
		},
		{
			name:    "Populated ram returns 2 at address 0x01",
			ram:     populatedRam,
			address: 0x01,
			want:    2,
		},
		{
			name:    "Populated ram returns 3 at address 0x02",
			ram:     populatedRam,
			address: 0x02,
			want:    3,
		},
		{
			name:    "Populated ram returns 4 at address 0x03",
			ram:     populatedRam,
			address: 0x03,
			want:    4,
		},
		{
			name:    "Populated ram returns 5 at address 0x04",
			ram:     populatedRam,
			address: 0x04,
			want:    5,
		},
		{
			name:    "Populated ram returns 6 at address 0x05",
			ram:     populatedRam,
			address: 0x05,
			want:    6,
		},
		{
			name:    "Populated ram returns 7 at address 0x06",
			ram:     populatedRam,
			address: 0x06,
			want:    7,
		},
		{
			name:    "Populated ram returns 8 at address 0x07",
			ram:     populatedRam,
			address: 0x07,
			want:    8,
		},
		// Wrap around tests
		{
			name:    "Populated ram returns 1 at address 0x08",
			ram:     populatedRam,
			address: 0x08,
			want:    1,
		},
		{
			name:    "Populated ram returns 2 at address 0x09",
			ram:     populatedRam,
			address: 0x09,
			want:    2,
		},
		{
			name:    "Populated ram returns 3 at address 0x0A",
			ram:     populatedRam,
			address: 0x0A,
			want:    3,
		},
		{
			name:    "Populated ram returns 4 at address 0x0B",
			ram:     populatedRam,
			address: 0x0B,
			want:    4,
		},
		{
			name:    "Populated ram returns 5 at address 0x0C",
			ram:     populatedRam,
			address: 0x0C,
			want:    5,
		},
		{
			name:    "Populated ram returns 6 at address 0x0D",
			ram:     populatedRam,
			address: 0x0D,
			want:    6,
		},
		{
			name:    "Populated ram returns 7 at address 0x0E",
			ram:     populatedRam,
			address: 0x0E,
			want:    7,
		},
		{
			name:    "Populated ram returns 8 at address 0x0F",
			ram:     populatedRam,
			address: 0x0F,
			want:    8,
		},
		// Test address space limits
		{
			name:    "Populated ram returns 1 at address 0xFFF8",
			ram:     populatedRam,
			address: 0xFFF8,
			want:    1,
		},
		{
			name:    "Populated ram returns 2 at address 0xFFF9",
			ram:     populatedRam,
			address: 0xFFF9,
			want:    2,
		},
		{
			name:    "Populated ram returns 3 at address 0xFFFA",
			ram:     populatedRam,
			address: 0xFFFA,
			want:    3,
		},
		{
			name:    "Populated ram returns 14 at address 0xFFFB",
			ram:     populatedRam,
			address: 0xFFFB,
			want:    4,
		},
		{
			name:    "Populated ram returns 5 at address 0xFFFC",
			ram:     populatedRam,
			address: 0xFFFC,
			want:    5,
		},
		{
			name:    "Populated ram returns 6 at address 0xFFFD",
			ram:     populatedRam,
			address: 0xFFFD,
			want:    6,
		},
		{
			name:    "Populated ram returns 7 at address 0xFFFE",
			ram:     populatedRam,
			address: 0xFFFE,
			want:    7,
		},
		{
			name:    "Populated ram returns 8 at address 0xFFFF",
			ram:     populatedRam,
			address: 0xFFFF,
			want:    8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ram, err := NewRepeatingRam(EightBytes)
			if err != nil {
				panic(err)
			}
			ram.ram = tt.ram

			if got := ram.Read(tt.address); got != tt.want {
				t.Errorf("Read() = %v, want = %v", got, tt.want)
			}
		})
	}

	// Check support for nil
	if got := (*RepeatingRam)(nil).Read(0); got != 0 {
		t.Errorf("Nil Read() = %v, want = 0", got)
	}
}

func TestRepeatingRam_Write(t *testing.T) {
	tests := []struct {
		name    string
		address Address
		value   uint8
		want    []uint8
	}{
		{
			name:    "Write 1 at address 0x00",
			address: 0x00,
			value:   1,
			want:    []uint8{1, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 2 at address 0x01",
			address: 0x01,
			value:   2,
			want:    []uint8{0, 2, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 3 at address 0x02",
			address: 0x02,
			value:   3,
			want:    []uint8{0, 0, 3, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 4 at address 0x03",
			address: 0x03,
			value:   4,
			want:    []uint8{0, 0, 0, 4, 0, 0, 0, 0},
		},
		{
			name:    "Write 0x15 at address 0x04",
			address: 0x04,
			value:   0x15,
			want:    []uint8{0, 0, 0, 0, 0x15, 0, 0, 0},
		},
		{
			name:    "Write 0xFF at address 0x05",
			address: 0x05,
			value:   0xFF,
			want:    []uint8{0, 0, 0, 0, 0, 0xFF, 0, 0},
		},
		{
			name:    "Write 7 at address 0x06",
			address: 0x06,
			value:   7,
			want:    []uint8{0, 0, 0, 0, 0, 0, 7, 0},
		},
		{
			name:    "Write 0x0A at address 0x07",
			address: 0x07,
			value:   0x0A,
			want:    []uint8{0, 0, 0, 0, 0, 0, 0, 0x0A},
		},
		// Wrap around tests
		{
			name:    "Write 1 at address 0x08",
			address: 0x08,
			value:   1,
			want:    []uint8{1, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 2 at address 0x09",
			address: 0x09,
			value:   2,
			want:    []uint8{0, 2, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 3 at address 0x0A",
			address: 0x0A,
			value:   3,
			want:    []uint8{0, 0, 3, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 4 at address 0x0B",
			address: 0x0B,
			value:   4,
			want:    []uint8{0, 0, 0, 4, 0, 0, 0, 0},
		},
		{
			name:    "Write 5 at address 0x0C",
			address: 0x0C,
			value:   5,
			want:    []uint8{0, 0, 0, 0, 5, 0, 0, 0},
		},
		{
			name:    "Write 6 at address 0x0D",
			address: 0x0D,
			value:   6,
			want:    []uint8{0, 0, 0, 0, 0, 6, 0, 0},
		},
		{
			name:    "Write 7 at address 0x0E",
			address: 0x0E,
			value:   7,
			want:    []uint8{0, 0, 0, 0, 0, 0, 7, 0},
		},
		{
			name:    "Write 8 at address 0x0F",
			address: 0x0F,
			value:   8,
			want:    []uint8{0, 0, 0, 0, 0, 0, 0, 8},
		},
		// Test address space limits
		{
			name:    "Write 1 at address 0xFFF8",
			address: 0xFFF8,
			value:   1,
			want:    []uint8{1, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 2 at address 0xFFF9",
			address: 0xFFF9,
			value:   2,
			want:    []uint8{0, 2, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 3 at address 0xFFFA",
			value:   3,
			address: 0xFFFA,
			want:    []uint8{0, 0, 3, 0, 0, 0, 0, 0},
		},
		{
			name:    "Write 14 at address 0xFFFB",
			address: 0xFFFB,
			value:   4,
			want:    []uint8{0, 0, 0, 4, 0, 0, 0, 0},
		},
		{
			name:    "Write 5 at address 0xFFFC",
			address: 0xFFFC,
			value:   5,
			want:    []uint8{0, 0, 0, 0, 5, 0, 0, 0},
		},
		{
			name:    "Write 6 at address 0xFFFD",
			address: 0xFFFD,
			value:   6,
			want:    []uint8{0, 0, 0, 0, 0, 6, 0, 0},
		},
		{
			name:    "Write 7 at address 0xFFFE",
			address: 0xFFFE,
			value:   7,
			want:    []uint8{0, 0, 0, 0, 0, 0, 7, 0},
		},
		{
			name:    "Write 8 at address 0xFFFF",
			address: 0xFFFF,
			value:   8,
			want:    []uint8{0, 0, 0, 0, 0, 0, 0, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ram, err := NewRepeatingRam(EightBytes)
			if err != nil {
				panic(err)
			}
			ram.Write(tt.address, tt.value)

			if !reflect.DeepEqual(ram.ram, tt.want) {
				t.Errorf("Write() got = %v, want = %v", ram.ram, tt.want)
			}
		})
	}

	// Check support for nil
	(*RepeatingRam)(nil).Write(0, 0)
}

func TestNewRepeatingRam(t *testing.T) {
	tests := []struct {
		name    string
		size    RepeatingRamSize
		want    RepeatingRam
		wantErr bool
	}{
		{
			name:    "Invalid size 0 will error",
			wantErr: true,
		},
		{
			name:    "Invalid size 13 will error",
			size:    13,
			wantErr: true,
		},
		{
			name: "Valid size of 8 bytes is fine",
			size: EightBytes,
			want: RepeatingRam{
				ram:  make([]uint8, EightBytes),
				size: Address(EightBytes),
				mask: Address(EightBytes - 1),
			},
		},
		{
			name: "Valid size of 16 bytes is fine",
			size: SixteenBytes,
			want: RepeatingRam{
				ram:  make([]uint8, SixteenBytes),
				size: Address(SixteenBytes),
				mask: Address(SixteenBytes - 1),
			},
		},
		{
			name: "Valid size of 32 bytes is fine",
			size: ThirtyTwoBytes,
			want: RepeatingRam{
				ram:  make([]uint8, ThirtyTwoBytes),
				size: Address(ThirtyTwoBytes),
				mask: Address(ThirtyTwoBytes - 1),
			},
		},
		{
			name: "Valid size of 64 bytes is fine",
			size: SixtyFourBytes,
			want: RepeatingRam{
				ram:  make([]uint8, SixtyFourBytes),
				size: Address(SixtyFourBytes),
				mask: Address(SixtyFourBytes - 1),
			},
		},
		{
			name: "Valid size of 128 bytes is fine",
			size: OneHundredAndTwentyEightBytes,
			want: RepeatingRam{
				ram:  make([]uint8, OneHundredAndTwentyEightBytes),
				size: Address(OneHundredAndTwentyEightBytes),
				mask: Address(OneHundredAndTwentyEightBytes - 1),
			},
		},
		{
			name: "Valid size of 1024 bytes is fine",
			size: OneKiloByte,
			want: RepeatingRam{
				ram:  make([]uint8, OneKiloByte),
				size: Address(OneKiloByte),
				mask: Address(OneKiloByte - 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRepeatingRam(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRepeatingRam() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepeatingRam() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
