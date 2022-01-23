// Package cmd is command-line interface for encrypted file sender.
package main

import (
	"fmt"
	"os"

  "github.com/m3dsh/encrypted-file-sender/nc"
	"github.com/urfave/cli/v2"
)

// connectAction do actions required for client-side
func connectAction(c *cli.Context) error{
  if c.Args().Len() >0 {
  n := &nc.NetData{
     Host: c.Args().Get(0),
     Port: c.String("port"),
     IsListen: false,
   }

   nc.Run(n)
  return nil
  }

  return fmt.Errorf("no <IP> is specified, set the <IP> " +
  "that you want to connect to")
}

// listenAction do actions requred for server-side
func listenAction(c *cli.Context) error{
  n := &nc.NetData{
    Host: "0.0.0.0",
    Port:  c.String("port"),
    IsListen: true,
  }
  nc.Run(n)
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

// shared flags between multiple subcommands
  portFlag := &cli.StringFlag{
          Name:     "port",
          Aliases:  []string{"p"}, 
          Usage:    "Define `PORT`",
          Value:    "8888",
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

  // RUN app
  if err:= app.Run(os.Args); err != nil {
    fmt.Println(err)
  }
}
