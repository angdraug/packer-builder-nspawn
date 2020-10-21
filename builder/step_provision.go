package builder

import (
	"context"
	"fmt"

	"github.com/coreos/go-systemd/v22/dbus"
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
		`rm /etc/hostname && ` +
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

	finished, err := exec.WaitFor("Startup finished", s.monitor(machine)...)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	if err := exec.Run("/usr/bin/machinectl", "start", machine); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	if !<-finished {
		state.Put("error", fmt.Errorf("Startup timed out after %s", config.Timeout))
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
	machine := config.PackerBuildName

	exec.Run(
		"/usr/bin/systemd-run", "-M", machine, "-P", "--wait", "-q",
		"/usr/bin/apt-get", "clean",
	)

	deadChan, err := exec.WaitFor("dead", s.monitor(machine)...)
	exec.Run("/usr/bin/machinectl", "stop", machine)
	if err == nil {
		<-deadChan
	}
}

func (s *StepProvision) monitor(machine string) []string {
	return []string{
		"/usr/bin/gdbus", "monitor", "--system", "--dest", "org.freedesktop.systemd1",
		"--object-path", fmt.Sprintf(
			"/org/freedesktop/systemd1/unit/systemd_2dnspawn_40%s_2eservice",
			dbus.PathBusEscape(machine)),
	}
}
