package builder

import (
	"context"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepProvision struct{}

func (s *StepProvision) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	hook := state.Get("hook").(packer.Hook)
	ui := state.Get("ui").(packer.Ui)
	machine := state.Get("machine").(*Machine)

	machine.Start()

	comm := &Communicator{machine}

	hookData := common.PopulateProvisionHookData(state)

	ui.Say("Running the provision hook")
	if err := hook.Run(ctx, packer.HookProvision, ui, comm, hookData); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepProvision) Cleanup(state multistep.StateBag) {
	machine := state.Get("machine").(*Machine)
	machine.Stop()
}
