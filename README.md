# Packer builder for systemd-nspawn/debootstrap

This plugin uses debootstrap to create an image suitable for use with systemd-nspawn.

## Quick Start

```
apt install debootstrap packer golang-go
git clone https://git.sr.ht/~angdraug/packer-builder-nspawn-debootstrap
cd packer-builder-nspawn-debootstrap
go build
packer build unstable-minbase.json
```

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
