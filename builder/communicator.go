package builder

import (
	"context"
	"io"
	"os"

	"github.com/hashicorp/packer/packer"
)

type Communicator struct {
	machine *Machine
}

func (c *Communicator) Start(ctx context.Context, remote *packer.RemoteCmd) error {
	go func() {
		err := c.machine.Run("/bin/sh", "-c", remote.Command)
		if err != nil {
			remote.SetExited(1)
		} else {
			remote.SetExited(0)
		}
	}()

	return nil
}

func (c *Communicator) Upload(dst string, r io.Reader, fi *os.FileInfo) error {
	return c.machine.Write(dst, r)
}

func (c *Communicator) UploadDir(dst string, src string, exclude []string) error {
	return c.machine.CopyTo(src, dst)
}

func (c *Communicator) Download(src string, w io.Writer) error {
	return c.machine.Read(src, w)
}

func (c *Communicator) DownloadDir(src string, dst string, exclude []string) error {
	return c.machine.CopyFrom(src, dst)
}
