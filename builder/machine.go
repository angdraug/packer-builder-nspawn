package builder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
)

type Machine struct {
	name         string
	exec         ExecWrapper
	machines_dir string
}

func (m *Machine) Path() string {
	return filepath.Join(m.machines_dir, m.name)
}

func (m *Machine) Exists() bool {
	_, err := os.Stat(m.Path())
	return (err == nil)
}

func (m *Machine) Remove() error {
	args := []string{"/usr/bin/machinectl", "remove", m.name}
	return m.exec.Run(args...)
}

func (m *Machine) Chroot(args ...string) error {
	args = append([]string{"/usr/bin/systemd-nspawn", "-M", m.name, "-U"}, args...)
	return m.exec.Run(args...)
}

func (m *Machine) monitor() []string {
	return []string{
		"/usr/bin/gdbus", "monitor", "--system", "--dest", "org.freedesktop.systemd1",
		"--object-path", fmt.Sprintf(
			"/org/freedesktop/systemd1/unit/systemd_2dnspawn_40%s_2eservice",
			dbus.PathBusEscape(m.name)),
	}
}

func (m *Machine) RunAndWait(command string, marker string) error {
	finished, err := m.exec.WaitFor("Startup finished", m.monitor()...)
	if err != nil {
		return err
	}

	if err := m.exec.Run("/usr/bin/machinectl", command, m.name); err != nil {
		return err
	}

	if !<-finished {
		return fmt.Errorf("Startup timed out after %s", m.exec.timeout)
	}

	return nil
}

func (m *Machine) Start() error {
	return m.RunAndWait("start", "Startup finished")
}

func (m *Machine) Stop() error {
	return m.RunAndWait("stop", "dead")
}

func (m *Machine) RunLocal(args ...string) error {
	return m.exec.Run(args...)
}

func (m *Machine) run(args ...string) []string {
	return append([]string{
		"/usr/bin/systemd-run", "-M", m.name, "-P", "--wait", "-q",
	}, args...)
}

func (m *Machine) Run(args ...string) error {
	return m.exec.Run(m.run(args...)...)
}

func (m *Machine) Read(src string, w io.Writer) error {
	args := m.run("/bin/sh", "-c", fmt.Sprintf("cat < '%s'", src))
	return m.exec.Read(w, args...)
}

func (m *Machine) Write(dst string, r io.Reader) error {
	args := m.run("/bin/sh", "-c", fmt.Sprintf("cat > '%s'", dst))
	return m.exec.Write(r, args...)
}

func (m *Machine) CopyTo(src string, dst string) error {
	args := []string{"/usr/bin/machinectl", "copy-to", m.name, src, dst}
	return m.exec.Run(args...)
}

func (m *Machine) CopyFrom(src string, dst string) error {
	args := []string{"/usr/bin/machinectl", "copy-from", m.name, src, dst}
	return m.exec.Run(args...)
}

func (m *Machine) Clone(base string) error {
	args := []string{"/usr/bin/machinectl", "clone", base, m.name}
	return m.exec.Run(args...)
}

func isUrl(s string) bool {
	return strings.Contains(s, "://")
}

func isTar(s string) bool {
	fields := strings.FieldsFunc(s, func(c rune) bool { return (c == '.') })
	return fields[len(fields)-1] == "tar" || fields[len(fields)-2] == "tar"
}

func isTarZst(s string) bool {
	fields := strings.FieldsFunc(s, func(c rune) bool { return (c == '.') })
	return fields[len(fields)-2] == "tar" && fields[len(fields)-1] == "zst"
}

func isFile(s string) bool {
	fi, err := os.Stat(s)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

func isDir(s string) bool {
	fi, err := os.Stat(s)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func (m *Machine) Import(image string) error {
	var command string
	switch {
	case isUrl(image) && isTar(image):
		command = "pull-tar"
	case isUrl(image):
		command = "pull-raw"
	case isFile(image) && isTarZst(image):
		// machinectl doesn't expect .tar.zst
		args := []string{
			"/bin/sh", "-c",
			fmt.Sprintf("zstdcat '%s' | machinectl import-tar - '%s'", image, m.name),
		}
		return m.exec.Run(args...)
	case isFile(image) && isTar(image):
		command = "import-tar"
	case isFile(image):
		command = "import-raw"
	case isDir(image):
		command = "import-fs"
	default:
		return fmt.Errorf("Image %s isn't a URL, a file, or a directory", image)
	}
	args := []string{"/usr/bin/machinectl", command, image, m.name}
	return m.exec.Run(args...)
}
