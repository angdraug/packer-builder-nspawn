package builder

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepProvision struct{}

func (s *StepProvision) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	command := `echo pts/0 > /etc/securetty && ` +
		`systemctl enable systemd-networkd.service && ` +
		`systemctl enable systemd-resolved.service && ` +
		`apt-get clean && ` +
		`echo 'APT::Install-Recommends "False";' > /etc/apt/apt.conf.d/60no-install-recommends`

	args := []string{
		"/usr/bin/systemd-nspawn",
		"-D", config.Target,
		"-U", "-P",
		"sh", "-c", command,
	}

	if err := Run(ui, args...); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepProvision) Cleanup(state multistep.StateBag) {}
