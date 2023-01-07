package builder

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepClone struct{}

func (s *StepClone) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	machine := state.Get("machine").(*Machine)

	if err := machine.Clone(config.Clone); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepClone) Cleanup(state multistep.StateBag) {}
