package builder

import (
	"fmt"
	"os"
)

type Artifact struct {
	machine string
	path string
}

func (*Artifact) BuilderId() string { return BuilderId }

func (a *Artifact) Files() []string { return []string{a.path} }

func (a *Artifact) Id() string { return "Machine" }

func (a *Artifact) String() string { return fmt.Sprintf("nspawn container: %s", a.machine) }

func (a *Artifact) State(name string) interface{} { return nil }

func (a *Artifact) Destroy() error { return os.RemoveAll(a.path) }
