source "nspawn" "base" {
    variant = "minbase"
}

build {
  sources = ["source.nspawn.base"]

  post-processors {
    post-processor "shell-local" {
      command = "tar --zstd -C /var/lib/machines/base -cf base.tar.zst ."
    }

    post-processor "artifice" {
      files = ["base.tar.zst"]
    }
  }
}

source "nspawn" "test-clone" {
  clone = "base"
}

build {
  sources = ["source.nspawn.test-clone"]

  provisioner "apt" {
    packages = ["less", "vim-tiny"]
  }
}

source "nspawn" "test-import" {
  import = "base.tar.zst"
}

build {
  sources = ["source.nspawn.test-import"]

  provisioner "apt" {
    sources = ["deb http://security.debian.org/debian-security buster/updates main contrib"]
    keys = ["/etc/apt/trusted.gpg.d/debian-archive-buster-security-automatic.gpg"]
  }
}
