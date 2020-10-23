package builder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/packer/packer"
)

type ExecWrapper struct {
	ui      packer.Ui
	timeout time.Duration
}

func (e *ExecWrapper) Run(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	e.ui.Say(fmt.Sprintf("Running: %s", cmd))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go e.scan(stdout, e.ui.Say)
	go e.scan(stderr, e.ui.Error)

	return cmd.Wait()
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

func (e *ExecWrapper) WaitFor(match string, args ...string) (chan bool, error) {
	cmd := exec.Command(args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	foundChan := make(chan bool)

	go e.scan(stdout, func(line string) {
		if strings.Contains(line, match) {
			foundChan <- true
			cmd.Process.Kill()
		}
	})

	go func() {
		<-time.After(e.timeout)
		foundChan <- false
		cmd.Process.Kill()
	}()

	go func() { cmd.Wait() }()

	return foundChan, nil
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

func (e *ExecWrapper) scan(r io.Reader, f func(string)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		f(scanner.Text())
	}
}
