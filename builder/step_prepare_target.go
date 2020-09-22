package builder

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepPrepareTarget struct{}

func (s *StepPrepareTarget) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	if len(config.Target) > 0 {
		if _, err := os.Stat(config.Target); err == nil && config.PackerForce {
			ui.Say("Deleting previous target directory...")
			if err := os.RemoveAll(config.Target); err != nil {
				state.Put("error", err)
				return multistep.ActionHalt
			}
		}
	} else {
		tempDir, err := ioutil.TempDir("/var/lib/machines", "nspawn-debootstrap-")
		if err != nil {
			state.Put("error", err)
			return multistep.ActionHalt
		}
		config.Target = tempDir
		state.Put("tempdir", tempDir)
		ui.Say(fmt.Sprintf("Created temporary target directory: %s", tempDir))
	}

	return multistep.ActionContinue
}

func (s *StepPrepareTarget) Cleanup(state multistep.StateBag) {
	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if !(cancelled || halted) { return }

	tempdir, ok := state.GetOk("tempdir")
	if ok { os.RemoveAll(tempdir.(string)) }
}
