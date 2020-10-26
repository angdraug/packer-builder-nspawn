package main

import (
	"git.sr.ht/~angdraug/packer-builder-nspawn/builder"
	"github.com/hashicorp/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(builder.Builder))
	server.Serve()
}
