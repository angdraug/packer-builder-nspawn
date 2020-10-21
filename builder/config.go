//go:generate mapstructure-to-hcl2 -type Config
package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// Distribution release code name as recognized by debootstrap(8).
	// The default is `unstable`.
	Suite string `mapstructure:"suite"`
	// URL for the distribution mirror.
	// The default is https://deb.debian.org/debian.
	Mirror string `mapstructure:"mirror"`
	// Absolute path to a directory where .deb files will be cached.
	// The default is the host's APT cache at `/var/cache/apt/archives`.
	CacheDir string `mapstructure:"cache_dir"`
	// Absolute path to the directory where systemd-nspawn expects to find
	// the container chroots. The default is `/var/lib/machines`.
	MachinesDir string `mapstructure:"machines_dir"`
	// The bootstrap script variant as recognized by debootstrap(8).
	Variant string `mapstructure:"variant"`
	// The timeout in seconds to wait for the container to start.
	// The default is 20 seconds.
	Timeout time.Duration `mapstructure:"timeout"`

	ctx interpolate.Context
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

	if c.Timeout == 0 {
		c.Timeout = 20 * time.Second
	}

	return nil
}

func (c *Config) Path() string {
	return filepath.Join(c.MachinesDir, c.PackerBuildName)
}
