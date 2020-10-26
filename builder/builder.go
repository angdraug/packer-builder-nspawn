package builder

import (
	"context"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

const BuilderId = "angdraug.nspawn"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	err := b.config.Prepare(raws...)
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	machine := Machine{
		name:         b.config.PackerBuildName,
		exec:         ExecWrapper{ui, b.config.Timeout},
		machines_dir: b.config.MachinesDir,
	}

	steps := []multistep.Step{new(StepPrepareTarget)}

	switch {
	case len(b.config.Import) != 0:
		steps = append(steps, new(StepImport))
	case len(b.config.Clone) != 0:
		steps = append(steps, new(StepClone))
	default:
		steps = append(steps, new(StepDebootstrap))
	}

	steps = append(steps, new(StepProvision))

	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("machine", &machine)

	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	return &Artifact{machine}, nil
}
