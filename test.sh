#!/bin/sh -eux
packer build -only='*.base' .
packer build -only='*.test-*' .
