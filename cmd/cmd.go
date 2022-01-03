package cmd

import (
	"flag"
)

var (
	IsListen bool
	Port     int
	IP	string	
)

func init() {
	flag.BoolVar(&IsListen, "l", false, "Bind and listen for incomming connection")
	flag.IntVar(&Port, "p", 8888, "Set the port for TCP connection. The Default port is 8888")
	flag.Parse()
	IP =flag.Arg(0)
}

func Run() error{
	return nil	
}
