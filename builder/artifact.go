package builder

import (
	"fmt"
)

type Artifact struct {
	machine Machine
}

func (*Artifact) BuilderId() string { return BuilderId }

func (a *Artifact) Files() []string { return []string{} }

func (a *Artifact) Id() string { return "Machine" }

func (a *Artifact) String() string { return fmt.Sprintf("nspawn container: %s", a.machine.name) }

func (a *Artifact) State(name string) interface{} { return nil }

func (a *Artifact) Destroy() error { return a.machine.Remove() }
