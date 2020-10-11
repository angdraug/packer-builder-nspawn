package builder

import (
	"context"
	"os"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepPrepareTarget struct{}

func (s *StepPrepareTarget) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	if _, err := os.Stat(config.Path()); err == nil && config.PackerForce {
		ui.Say("Deleting previous target directory")
		if err := os.RemoveAll(config.Path()); err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *StepPrepareTarget) Cleanup(state multistep.StateBag) {
	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if !(cancelled || halted) {
		return
	}

	config := state.Get("config").(*Config)
	os.RemoveAll(config.Path())
}
