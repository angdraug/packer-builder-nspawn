package builder

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/hashicorp/packer/packer"
)

type ExecWrapper struct {
	ui packer.Ui
}

func (e *ExecWrapper) Run(args ...string) error {
	return e.wrap(func(cmd *exec.Cmd) error {
		e.ui.Say(fmt.Sprintf("Running: %s", cmd))
		return cmd.Run()
	}, args...)
}

func (e *ExecWrapper) Write(r io.Reader, args ...string) error {
	return e.wrap(func(cmd *exec.Cmd) error {
		w, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		go func() {
			defer w.Close()
			io.Copy(w, r)
		}()
		e.ui.Say(fmt.Sprintf("Writing to: %s", cmd))
		return cmd.Run()
	}, args...)
}

func (e *ExecWrapper) Read(w io.Writer, args ...string) error {
	return e.wrap(func(cmd *exec.Cmd) error {
		r, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		go func() {
			defer r.Close()
			io.Copy(w, r)
		}()
		e.ui.Say(fmt.Sprintf("Reading from: %s", cmd))
		return cmd.Run()
	}, args...)
}

func (e *ExecWrapper) wrap(f func(*exec.Cmd) error, args ...string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := f(cmd)

	if len(stdout.String()) > 0 {
		e.ui.Message(strings.TrimSpace(stdout.String()))
	}
	if len(stderr.String()) > 0 {
		e.ui.Error(strings.TrimSpace(stderr.String()))
	}

	return err
}
