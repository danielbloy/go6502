package processor

//
// Resources and references used to build the 6502 emulator:
//
//   * mass:werk by Norbert Landsteiner. This contains an excellent explanation of
//     the processor and appears to be correct throughout. This is my main reference
//     for any ambiguities. There is also a simulator (though its not as trivial
//     to use  as Easy 6502):
//       https://www.masswerk.at/6502/6502_instruction_set.html - reference
//       https://www.masswerk.at/6502/ - simulator
//
//   * Easy 6502 by Nick Morgan - an easy to use online 6502 emulator
//       https://skilldrick.github.io/easy6502/
//
//   * NMOS 6502 opcodes by John Pickens, Updated by Bruce Clark and by Ed Spittles
//       http://www.6502.org/tutorials/6502opcodes.html
//
//   * 6502 instruction reference on NesDev
//       https://www.nesdev.org/wiki/6502_instructions
//       https://www.nesdev.org/obelisk-6502-guide/
//       https://www.nesdev.org/obelisk-6502-guide/reference.html
//
//   * BCD mode information by Bruce Clark. This contains a very detailed explanation
//     and sample code to verify documented and undocumented operation.
//       http://www.6502.org/tutorials/decimal_mode.html
//
//   * Interrupts explanation by Garth Wilson
//       http://6502.org/tutorials/interrupts.html
//
//   * Validating your emulator. There are a range of resources for this, including:
//       GitHub - Klaus2m5/6502_65C02_functional_tests: Tests for all valid opcodes of the 6502 and 65C02 processor
//       https://github.com/Klaus2m5/6502_65C02_functional_tests
//
//       6502.org • View topic - Klaus' Test Suite Tutorial for Noobs?
//       http://forum.6502.org/viewtopic.php?f=8&t=5298
//
//       https://github.com/tom-seddon/6502-tests
//
//       https://github.com/tom-seddon/b2/tree/master/etc/testsuite-2.15
//
//
// My original inspiration for the emulator came from using fake6502 that I used previously
// in a C based emulator (fake6502.h/.c). I spotted some issues with the emulator when
// building my own so also investigated the skilldrick implementation that can be found here:
//    https://github.com/skilldrick/easy6502/blob/gh-pages/simulator/assembler.js

// TODO: Write instructions on how to use.
