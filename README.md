# Packer builder for systemd-nspawn

This plugin can build and provision systemd-nspawn containers. It can import or
clone an existing base image, or generate one from scratch using debootstrap.

## Quick Start

```
sudo apt-get install --no-install-recommends \
 debootstrap golang-go libglib2.0-bin packer systemd-container zstd
git clone https://git.sr.ht/~angdraug/packer-builder-nspawn
cd packer-builder-nspawn
go build
packer build -only='*.base' .
```

## Setup

Prerequisites:
- `debootstrap` to generate a minimal viable chroot image
- `golang-go` to build this plugin from source
- `libglib2.0-bin` to monitor container status with `gdbus`
- `packer` (the Debian package recommends docker, you don't need that)
- `systemd-container` systemd-nspawn and related tools
- (optional) `zstd` for creating and importing .tar.zst images

In most cases, you'll want your container to be able to connect to the network
to provision itself, for that you need to enable systemd-networkd and
systemd-resolved:

```
systemctl enable systemd-networkd.service
systemctl start systemd-networkd.service
systemctl enable systemd-resolved.service
systemctl start systemd-resolved.service
```

The included example `nspawn.pkr.hcl` uses zstd to compress the image tarball,
it is several times faster than gzip at the same or better compression ratio.
You don't need zstd if you use a different method for archiving and delivering
your images.

For compatibility with the Debian package of Packer that is built with a newer
version of [ugorji-go-codec](https://github.com/ugorji/go) than the one pinned
in Packer source, this plugin's `go.mod` includes a replace line to import a
similarly patched version of Packer source. To use this plugin with Packer
built from unpatched upstream source, comment out that replace line.

## Security

In Debian, unprivileged user namespaces are disabled by default and have to be
enabled for nspawn's `-U` option to have any effect:

```
echo kernel.unprivileged_userns_clone=1 > /etc/sysctl.d/nspawn.conf
systemctl restart systemd-sysctl.service
```

See discussion of this kernel feature and its security implications in
[Debian Wiki](https://wiki.debian.org/nspawn#Host_Preparation) and
[Linux Weekly News](https://lwn.net/Articles/673597/).

This builder will work with and without private user namespaces. If you want an
image built on a system with userns enabled to be usable on systems with userns
disabled, use the method offered in
[systemd-nspawn(1)](https://www.freedesktop.org/software/systemd/man/systemd-nspawn.html#-U)
to reset chroot file ownership before archiving the image:

```
systemd-nspawn ... --private-users=0 --private-users-chown
```

## Configuration

All configuration options for this plugin are optional.

- `name` - Standard Packer build name parameter. The default is the builder
  name `nspawn`. This will be used as container name and will be configured as
  the hostname within the container.

- `import` - Import container image from a URL, file, or a directory tree, in a
  format recognized by `import-*` and `pull-*` commands of
  [machinectl(1)](https://www.freedesktop.org/software/systemd/man/machinectl.html).

- `clone` - Name of a local container to clone. When neither `import` nor
  `clone` options are set, a new image will be created with
  [debootstrap(8)](https://manpages.debian.org/unstable/debootstrap/debootstrap.8.en.html).

- `suite` - Distribution release code name as recognized by
  [debootstrap(8)](https://manpages.debian.org/unstable/debootstrap/debootstrap.8.en.html).
  The default is `unstable`.

- `mirror` - URL for the distribution mirror. The default is
  [https://deb.debian.org/debian](https://deb.debian.org/debian).

- `variant` - The bootstrap script variant as recognized by
  [debootstrap(8)](https://manpages.debian.org/unstable/debootstrap/debootstrap.8.en.html).
  The default is to not pass `--variant` to debootstrap, which will install
  required and important packages. The `minbase` variant will only install
  required packages, the plugin explicitly adds several small important
  packages to make sure that even a `minbase` image has the CLI tools likely to
  be used in most provisioning scripts.

- `cache_dir` - Absolute path to a directory where .deb files will be cached.
  The default is the host's APT cache at `/var/cache/apt/archives`.

- `machines_dir` - Absolute path to the directory where systemd-nspawn expects
  to find the container chroots. Unless you know what you're doing, keep the
  default `/var/lib/machines`.

- `timeout` - The timeout in seconds to wait for container startup and
  shutdown. The default is 20 seconds.

See [nspawn.pkr.hcl](nspawn.pkr.hcl) for an example of how to build a minimal
base container, archive it into a tarball, clone it a new container, and import
a container image from the archived tarball.

## Copying

Copyright (c) 2020  Dmitry Borodaenko <angdraug@debian.org>

This program is free software. You can distribute/modify this program under
the terms of the GNU General Public License version 3 or later, or under
the terms of the Mozilla Public License, v. 2.0, at your discretion.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.

This Source Code Form is not "Incompatible With Secondary Licenses",
as defined by the Mozilla Public License, v. 2.0.
