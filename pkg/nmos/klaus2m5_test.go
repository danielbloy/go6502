package nmos

import (
	"go6502/pkg/processor"
	"os"
	"testing"
)

type Ram struct {
	ram []uint8
}

func (r *Ram) Read(address processor.Address) uint8 {
	if r == nil {
		return 0
	}
	return r.ram[address]
}

// Write a value to the repeating RAM.
func (r *Ram) Write(address processor.Address, value uint8) {
	if r == nil {
		return
	}
	r.ram[address] = value
}

//   - Validating your emulator. There are a range of resources for this, including:
//     GitHub - Klaus2m5/6502_65C02_functional_tests: Tests for all valid opcodes of the 6502 and 65C02 processor
//     https://github.com/Klaus2m5/6502_65C02_functional_tests
//
//     6502.org • View topic - Klaus' Test Suite Tutorial for Noobs?
//     http://forum.6502.org/viewtopic.php?f=8&t=5298
//
//     https://github.com/tom-seddon/6502-tests
//
//     https://github.com/tom-seddon/b2/tree/master/etc/testsuite-2.15
func TestUsingKlaus2m5FunctionalTest(t *testing.T) {

	ram := Ram{
		ram: make([]byte, 0x10000),
	}

	f, err := os.Open("6502_functional_test.bin")
	if err != nil {
		t.Fatal(err)
	}
	count, err := f.Read(ram.ram)
	if count != 0x10000 {
		t.Errorf("wrong number of bytes read, expected 0x10000, got 0x%x", count)
	}

	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Override reset vector to correct starting location!
	err = processor.WriteResetVectorToMemory(&ram, 0x400)
	if err != nil {
		t.Fatal(err)
	}

	cpu, err := New6502Cpu(&ram)
	if err != nil {
		t.Fatal(err)
	}

	err = cpu.Reset()
	if err != nil {
		t.Fatal(err)
	}

	const CYCLES = 100_000_000
	cycles := uint(0)

	for cycles < CYCLES {
		c, err := cpu.Step()
		cycles += c
		if err != nil {
			t.Fatal(err)
		}

		// Early exist if we get to the correct success location.
		if cpu.State.PC == 0x3469 {
			break
		}
	}

	t.Logf("State: %v, Cycles %v", cpu.State, cycles)
	if cpu.State.PC == 0x3469 {
		t.Logf("SUCCESS")
	} else {
		t.Fatal("FAIL")

	}
}

// TODO: Add the 65C02 extended opcode test from Klaus.
