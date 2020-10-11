//go:generate mapstructure-to-hcl2 -type Config
package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Suite               string `mapstructure:"suite"`
	Mirror              string `mapstructure:"mirror"`
	CacheDir            string `mapstructure:"cache_dir"`
	MachinesDir         string `mapstructure:"machines_dir"`
	Variant             string `mapstructure:"variant"`
	ctx                 interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) error {
	err := config.Decode(c, &config.DecodeOpts{
		Interpolate: true,
	}, raws...)
	if err != nil {
		return err
	}

	if c.Suite == "" {
		c.Suite = "unstable"
	}

	if c.Mirror == "" {
		c.Mirror = "https://deb.debian.org/debian"
	}

	if c.CacheDir == "" {
		c.CacheDir = "/var/cache/apt/archives"
	}

	if c.MachinesDir == "" {
		c.MachinesDir = "/var/lib/machines"
	}

	cache, err := os.Stat(c.CacheDir)
	if err != nil || !cache.IsDir() {
		return fmt.Errorf("Cache directory is not a directory: %s", c.CacheDir)
	}

	return nil
}

func (c *Config) Path() string {
	return filepath.Join(c.MachinesDir, c.PackerBuildName)
}
