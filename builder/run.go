package builder

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/packer/packer"
)

func Run(ui packer.Ui, args ...string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ui.Say(fmt.Sprintf("Running: %s", cmd))
	err := cmd.Run()
	if len(stdout.String()) > 0 {
		ui.Message(strings.TrimSpace(stdout.String()))
	}
	if len(stderr.String()) > 0 {
		ui.Error(strings.TrimSpace(stderr.String()))
	}

	return err
}
