package builder

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/packer/packer"
)

type Communicator struct {
	machine string
	exec    ExecWrapper
}

func (c *Communicator) Start(ctx context.Context, remote *packer.RemoteCmd) error {
	args := c.run(remote.Command)

	go func() {
		err := c.exec.Run(args...)
		if err != nil {
			remote.SetExited(1)
		} else {
			remote.SetExited(0)
		}
	}()

	return nil
}

func (c *Communicator) Upload(dst string, r io.Reader, fi *os.FileInfo) error {
	args := c.run(fmt.Sprintf("cat > '%s'", dst))
	return c.exec.Write(r, args...)
}

func (c *Communicator) UploadDir(dst string, src string, exclude []string) error {
	args := []string{"/usr/bin/machinectl", "copy-to", c.machine, src, dst}
	return c.exec.Run(args...)
}

func (c *Communicator) Download(src string, w io.Writer) error {
	args := c.run(fmt.Sprintf("cat < '%s'", src))
	return c.exec.Read(w, args...)
}

func (c *Communicator) DownloadDir(src string, dst string, exclude []string) error {
	args := []string{"/usr/bin/machinectl", "copy-from", c.machine, src, dst}
	return c.exec.Run(args...)
}

func (c *Communicator) run(command string) []string {
	return []string{
		"/usr/bin/systemd-run", "-M", c.machine, "-P", "--wait", "-q",
		"/bin/sh", "-c", command,
	}
}
