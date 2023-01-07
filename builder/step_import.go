package builder

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepImport struct{}

func (s *StepImport) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	machine := state.Get("machine").(*Machine)

	if err := machine.Import(config.Import); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepImport) Cleanup(state multistep.StateBag) {}
