//go:generate mapstructure-to-hcl2 -type Config
package builder

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

const BuilderId = "angdraug.nspawn-debootstrap"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Suite               string `mapstructure:"suite"`
	Target              string `mapstructure:"target"`
	Mirror              string `mapstructure:"mirror"`
	CacheDir            string `mapstructure:"cache_dir"`
	Variant             string `mapstructure:"variant"`
	ctx                 interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate: true,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	if b.config.Suite == "" {
		b.config.Suite = "unstable"
	}

	if b.config.Mirror == "" {
		b.config.Mirror = "http://deb.debian.org/debian"
	}

	if b.config.CacheDir == "" {
		b.config.CacheDir = "/var/cache/apt/archives"
	}

	cache, err := os.Stat(b.config.CacheDir)
	if err != nil || !cache.IsDir() {
		return nil, nil, fmt.Errorf("Cache directory is not a directory: %s", b.config.CacheDir)
	}

	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	steps := []multistep.Step{
		&StepPrepareTarget{},
		&StepDebootstrap{},
		&StepProvision{},
	}

	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("ui", ui)

	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	return &Artifact{dir: b.config.Target}, nil
}

type Artifact struct {
	dir string
}

func (*Artifact) BuilderId() string { return BuilderId }

func (a *Artifact) Files() []string { return []string{a.dir} }

func (a *Artifact) Id() string { return a.dir }

func (a *Artifact) String() string { return a.dir }

func (a *Artifact) State(name string) interface{} { return nil }

func (a *Artifact) Destroy() error { return os.RemoveAll(a.dir) }
