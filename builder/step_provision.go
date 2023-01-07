package builder

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepProvision struct{}

func (s *StepProvision) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	hook := state.Get("hook").(packer.Hook)
	ui := state.Get("ui").(packer.Ui)
	machine := state.Get("machine").(*Machine)

	machine.Start()

	comm := &Communicator{machine}

	hookData := commonsteps.PopulateProvisionHookData(state)

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
	command := `rm -f /var/lib/dbus/machine-id && ` +
		`cat /dev/null > /etc/machine-id`
	machine.Chroot("/bin/sh", "-c", command)
}
