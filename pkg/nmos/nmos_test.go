package nmos

import (
	"go6502/pkg/processor"
	"reflect"
	"testing"
)

func TestNew6502Cpu(t *testing.T) {
	memory, err := processor.NewRepeatingRam(processor.SixteenBytes)
	if err != nil {
		panic(err)
	}

	is, err := New6502InstructionSet()
	if err != nil {
		panic(err)
	}

	want, err := processor.NewCpu(is, &memory)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		memory  processor.Memory
		wantErr bool
	}{
		{
			name:    "Nil memory should error",
			wantErr: true,
		},
		{
			name:   "Valid memory should be fine",
			memory: &memory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New6502Cpu(tt.memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("New6502Cpu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got.State, want.State) {
				t.Errorf("New6502Cpu() got State = %v, want State %v", got.State, want.State)
			}

			gotMemory, _ := got.Memory()
			if !reflect.DeepEqual(gotMemory, tt.memory) {
				t.Errorf("New6502Cpu() got memory = %v, want memory %v", gotMemory, tt.memory)
			}

			// We can only really compare opcode as reflect.DeepEqual does not work with function pointers.
			wantOpcodes, _ := want.Opcodes()
			gotOpcodes, _ := got.Opcodes()

			if !reflect.DeepEqual(wantOpcodes, gotOpcodes) {
				t.Errorf("New6502Cpu() got opcodes = %v, want opcodes %v", wantOpcodes, gotOpcodes)
			}
		})
	}
}

func TestNew65C02Cpu(t *testing.T) {
	memory, err := processor.NewRepeatingRam(processor.SixteenBytes)
	if err != nil {
		panic(err)
	}

	is, err := New65C02InstructionSet()
	if err != nil {
		panic(err)
	}

	want, err := processor.NewCpu(is, &memory)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		memory  processor.Memory
		wantErr bool
	}{
		{
			name:    "Nil memory should error",
			wantErr: true,
		},
		{
			name:   "Valid memory should be fine",
			memory: &memory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New65C02Cpu(tt.memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("New65C02Cpu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got.State, want.State) {
				t.Errorf("New65C02Cpu() got State = %v, want State %v", got.State, want.State)
			}

			gotMemory, _ := got.Memory()
			if !reflect.DeepEqual(gotMemory, tt.memory) {
				t.Errorf("New65C02Cpu() got memory = %v, want memory %v", gotMemory, tt.memory)
			}

			// We can only really compare opcode as reflect.DeepEqual does not work with function pointers.
			wantOpcodes, _ := want.Opcodes()
			gotOpcodes, _ := got.Opcodes()

			if !reflect.DeepEqual(wantOpcodes, gotOpcodes) {
				t.Errorf("New65C02Cpu() got opcodes = %v, want opcodes %v", wantOpcodes, gotOpcodes)
			}
		})
	}
}
