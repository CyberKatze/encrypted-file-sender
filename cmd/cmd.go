// Package cmd is command-line interface for encrypted file sender.
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// Variables that is being set by CLI
var (
  IsListen  bool
  Port      string
  IP        string
)

// connectAction do actions required for client-side
func connectAction(c *cli.Context) error{
  if c.Args().Len() >0 {
  IsListen = false
  IP = c.Args().Get(0)
  Port = c.String("port")
  fmt.Printf("connect to %s:%s",IP, Port)
  return nil
  }

  return fmt.Errorf("no <IP> is specified, set the <IP> " +
  "that you want to connect to")
}

// listenAction do actions requred for server-side
func listenAction(c *cli.Context) error{
  IsListen = true
  Port = c.String("port")
  fmt.Printf("listen to 0.0.0.0:%s", Port)
  return nil
  
}

func main() {
  // Version of app
  cli.VersionFlag = &cli.BoolFlag{
    Name:  "version",
    Aliases: []string{"v"},
    Usage: "output version information",
  }

// CLI app spec
  app := &cli.App{
    Name:     "encrypted-file-sender",
    Usage:    "send encrypted file over network",
    Version:  "1.0.0",
  }

// Global flags
  app.Flags = []cli.Flag{
    &cli.BoolFlag{
      Name:     "verbose",
      Aliases:  []string{"V"}, 
      Usage: "verbose output",
    },
  }

// share flag between multiple subcommand
  portFlag := &cli.StringFlag{
          Name:     "port",
          Aliases:  []string{"p"}, 
          Usage:    "Define `PORT`",
          Required: true,
        }

// CLI subcommands
  app.Commands = []*cli.Command{
    {
      Name:     "connect",
      Category: "Network",
      Usage:    "Connect [-p <PORT>] <IP>",
      Flags:    []cli.Flag{
        portFlag,
      },
      Action:   connectAction,
    },
    {
      Name:     "listen",
      Category: "Network",
      Usage:    "listen [-p <PORT>]",
      Flags:    []cli.Flag{
        portFlag,
      },
      Action:   listenAction,
    },
  }

  // RUn app
  if err:= app.Run(os.Args); err != nil {
    fmt.Println(err)
  }
}
