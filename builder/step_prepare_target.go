package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepPrepareTarget struct{}

func (s *StepPrepareTarget) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)
	machine := state.Get("machine").(*Machine)

	if machine.Exists() {
		if config.PackerForce {
			ui.Say(fmt.Sprintf("Container %s already exists, removing", machine.name))
			if err := machine.Remove(); err != nil {
				state.Put("error", err)
				return multistep.ActionHalt
			}
		} else {
			state.Put("error", fmt.Errorf("Container %s already exists", machine.name))
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

	machine := state.Get("machine").(*Machine)
	machine.Remove()
}
