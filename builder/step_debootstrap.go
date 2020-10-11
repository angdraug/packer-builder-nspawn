package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
)

type StepDebootstrap struct{}

func (s *StepDebootstrap) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	exec := state.Get("exec").(ExecWrapper)

	args := []string{
		"/usr/sbin/debootstrap",
		"--include=apt-utils,iputils-ping,netbase,procps,systemd-container",
		fmt.Sprintf("--cache-dir=%s", config.CacheDir),
	}
	if config.Variant != "" {
		args = append(args, fmt.Sprintf("--variant=%s", config.Variant))
	}
	args = append(args, config.Suite, config.Path(), config.Mirror)

	if err := exec.Run(args...); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDebootstrap) Cleanup(state multistep.StateBag) {}
