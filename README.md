# Packer builder for systemd-nspawn/debootstrap

This plugin uses debootstrap to create an image suitable for use with
systemd-nspawn.

## Quick Start

```
sudo apt-get install --no-install-recommends debootstrap packer golang-go
git clone https://git.sr.ht/~angdraug/packer-builder-nspawn-debootstrap
cd packer-builder-nspawn-debootstrap
go build
sudo packer build unstable-minbase.json
```

## Setup

Prerequisites:
- `golang-go` to build this plugin from source
- `packer` (the Debian package recommends docker, you don't need that)
- `debootstrap` to generate a minimal viable chroot image
- `systemd-container` to install nspawn and related tools
- (optional) `zstd` for faster image tarball compression

In most cases, you'll want your container to be able to connect to the network
to provision itself, for that you need to enable systemd-networkd and
systemd-resolved:

```
systemctl enable systemd-networkd.service
systemctl start systemd-networkd.service
systemctl enable systemd-resolved.service
systemctl start systemd-resolved.service
```

The included example `unstable-minbase.json` uses zstd to compress the image
tarball, because it is several times faster than gzip at the same or better
compression ratio. You don't need zstd if you use a different method for
archiving and delivering your images.

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
  name `nspawn-debootstrap`. This will be used as container name and will be
  configured as the hostname within the container.

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

See [unstable-minbase.json](/unstable-minbase.json) for an example of how to
build a minimal base image with a unique name, install additional software in
it during provisioning, and archive it into a tarball.

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
