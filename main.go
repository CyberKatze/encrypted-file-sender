package main

import (
	"log"

	"github.com/m3dsh/encrypted-file-sender/client"
	"github.com/m3dsh/encrypted-file-sender/cmd"
	"github.com/m3dsh/encrypted-file-sender/server"
)

var app string

func main() {
	//get app value from cli
	//it can be either server or client
	app, err := cmd.GetApp()
	if err != nil {
		log.Fatal(err)
	}
	if app == "server" {
		server.Start()
	} else {
		client.Start()
	}

}
