package builder

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepDebootstrap struct{}

func (s *StepDebootstrap) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	args := []string{
		"--include=systemd-container",
		fmt.Sprintf("--cache-dir=%s", config.CacheDir),
	}
	if config.Variant != "" {
		args = append(args, fmt.Sprintf("--variant=%s", config.Variant))
	}
	args = append(args, config.Suite, config.Target, config.Mirror)

	cmd := exec.Command("/usr/sbin/debootstrap", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ui.Say(fmt.Sprintf("Running debootstrap: %s", cmd))
	err := cmd.Run()
	log.Printf("stdout: %s", stdout)
	log.Printf("stderr: %s", stderr)
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDebootstrap) Cleanup(state multistep.StateBag) {}
