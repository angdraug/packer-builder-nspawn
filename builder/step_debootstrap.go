package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
)

type StepDebootstrap struct{}

func (s *StepDebootstrap) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	machine := state.Get("machine").(*Machine)

	args := []string{
		"/usr/sbin/debootstrap",
		"--include=apt-utils,iputils-ping,netbase,procps,systemd-container",
		fmt.Sprintf("--cache-dir=%s", config.CacheDir),
	}
	if config.Variant != "" {
		args = append(args, fmt.Sprintf("--variant=%s", config.Variant))
	}
	args = append(args, config.Suite, machine.Path(), config.Mirror)

	if err := machine.RunLocal(args...); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	command := `echo pts/0 > /etc/securetty && ` +
		`rm /etc/hostname && ` +
		`systemctl enable systemd-networkd.service && ` +
		`systemctl enable systemd-resolved.service && ` +
		`echo 'APT::Install-Recommends "False";' ` +
			`> /etc/apt/apt.conf.d/60no-install-recommends && ` +
		`echo 'path-exclude /usr/share/doc/*\n` +
			`path-include /usr/share/doc/*/copyright\n` +
			`path-exclude /usr/share/info/*\n` +
			`path-exclude /usr/share/lintian/*\n` +
			`path-exclude /usr/share/man/*' ` +
			`> /etc/dpkg/dpkg.cfg.d/no-documentation && ` +
		`find /usr \( -name copyright -prune \) -o \( ` +
			`-path '/usr/share/doc/*' -o ` +
			`-path '/usr/share/info/*' -o ` +
			`-path '/usr/share/lintian/*' -o ` +
			`-path '/usr/share/man/*' \) -type f -print | xargs rm -f`

	if err := machine.Chroot("/bin/sh", "-c", command); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDebootstrap) Cleanup(state multistep.StateBag) {
	machine := state.Get("machine").(*Machine)
	machine.Chroot("/usr/bin/apt-get", "clean")
}
