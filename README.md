# go6502

Please see my website [Code Club Adventures](http://codeclubadventures.com/) for more coding materials.

A virtual machine based around a 6502 processor written in Go.

### Motivation

During lockdown 2020 I experimented with a simple 6502 virtual machine
by placing the existing C fake6502 emulator into my own wrapper. It
worked reasonably well but was missing basic things like comprehensive 
tests and a decent UI (don't judge me too harshly, it was a throw away
project just for fun). 

I was never satisfied with the full quality of the user experience and 
wanted to do more. In particular I wanted an easy to extend 6502 emulator
that allowed me to simulate many different variants of 6502 processors
and even generate my own. I also much prefer Go to C or C++.

This project is the result. Currently, it has a functionally working
processor that passes the Klaus2m5 functional test suite. Most machine
cycle counts will be correct but there is still work to do to make this
fully accurate, particularly around the branching instructions.

There is a decent list of work still to do, in no particular order:
* Accurate instruction counts.
* BCD tests for undocumented behaviour.
* Undocumented and extended instructions.
* 65C02 support.
* DFBP custom CPU support (for experiments).
* Assembler/Disassembler.
* Debugger.
* VM wrapper to allow execution of arbitrary programs.
