package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepDebootstrap struct{}

func (s *StepDebootstrap) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	args := []string{
		"/usr/sbin/debootstrap",
		"--include=systemd-container",
		fmt.Sprintf("--cache-dir=%s", config.CacheDir),
	}
	if config.Variant != "" {
		args = append(args, fmt.Sprintf("--variant=%s", config.Variant))
	}
	args = append(args, config.Suite, config.Target, config.Mirror)

	if err := Run(ui, args...); err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDebootstrap) Cleanup(state multistep.StateBag) {}
