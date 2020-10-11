package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepProvision struct{}

func (s *StepProvision) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	hook := state.Get("hook").(packer.Hook)
	ui := state.Get("ui").(packer.Ui)
	exec := state.Get("exec").(ExecWrapper)
	machine := config.PackerBuildName

	command := `echo pts/0 > /etc/securetty && ` +
		fmt.Sprintf(`echo %s > /etc/hostname &&`, machine) +
		`systemctl enable systemd-networkd.service && ` +
		`systemctl enable systemd-resolved.service && ` +
		`echo 'APT::Install-Recommends "False";' > /etc/apt/apt.conf.d/60no-install-recommends`

	args := []string{
		"/usr/bin/systemd-nspawn", "-M", machine, "-U",
		"sh", "-c", command,
	}

	if err := exec.Run(args...); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	if err := exec.Run("/usr/bin/machinectl", "start", machine); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	comm := &Communicator{machine, exec}

	hookData := common.PopulateProvisionHookData(state)

	ui.Say("Running the provision hook")
	if err := hook.Run(ctx, packer.HookProvision, ui, comm, hookData); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepProvision) Cleanup(state multistep.StateBag) {
	config := state.Get("config").(*Config)
	exec := state.Get("exec").(ExecWrapper)

	exec.Run(
		"/usr/bin/systemd-run", "-M", config.PackerBuildName, "-P", "--wait", "-q",
		"/usr/bin/apt-get", "clean",
	)

	exec.Run("/usr/bin/machinectl", "stop", config.PackerBuildName)
}
