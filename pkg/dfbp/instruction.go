package dfbp

import (
	"go6502/pkg/nmos"
	"go6502/pkg/processor"
)

// TODO
func NewDfbpInstructionSet() (processor.InstructionSet, error) {
	instructions, err := nmos.New6502InstructionSet()
	if err != nil {
		return processor.InstructionSet{}, err
	}

	// TODO: Add in DFBP specific instructions.

	return instructions, nil
}
