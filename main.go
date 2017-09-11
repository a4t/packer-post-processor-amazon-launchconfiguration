package main

import (
	"log"

	"github.com/a4t/packer-post-processor-amazon-launchconfiguration/packer-post-processor-amazon-launchconfiguration"
	"github.com/hashicorp/packer/packer/plugin"
)

var (
	Version  string
	Revision string
)

func main() {
	log.Println("version: " + Version)
	log.Println("revision: " + Revision)

	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterPostProcessor(&amazonlaunchconfiguration.PostProcessor{})
	server.Serve()
}
