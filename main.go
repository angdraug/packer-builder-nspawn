package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/angdraug/packer-builder-nspawn-debootstrap/builder"
)

func main() {
	server, err := plugin.Server()
	if err != nil { panic(err) }
	server.RegisterBuilder(new(builder.Builder))
	server.Serve()
}
