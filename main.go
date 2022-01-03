package main

import (
	"github.com/m3dsh/encrypted-file-sender/client"
	"github.com/m3dsh/encrypted-file-sender/cmd"
	"github.com/m3dsh/encrypted-file-sender/server"
)

var app string

func main() {
	//get app value from cli
	//it can be either server or client
	_ = cmd.Run()
	if cmd.IsListen {
		server.Start(cmd.Port)
	} else {
		client.Start(cmd.IP, cmd.Port)
	}

}
