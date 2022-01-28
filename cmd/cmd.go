// Package cmd is command-line interface for encrypted file sender.
package main

import (
	"fmt"
	"os"

	"github.com/m3dsh/encrypted-file-sender/nc"
	"github.com/m3dsh/encrypted-file-sender/encryption"
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

func encryptAction(c *cli.Context) error{
  encryption.EncryptFiles(c.String("file"),c.String("algorithm"),c.String("key"))
  return nil
}
func decryptAction(c *cli.Context) error{
  encryption.DecryptFiles(c.String("file"),c.String("algorithm"),c.String("key"))
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
        &cli.StringFlag{
          Name:     "files",
          Aliases:  []string{"d","f","o"}, 
          Usage: " path to directory or file in order to send ",
          Value: "./",
        },
      },
      Action:   connectAction,
    },
    {
      Name:     "listen",
      Category: "Network",
      Usage:    "listen [-p <PORT>]",
      Flags:    []cli.Flag{
        portFlag,
        &cli.StringFlag{
          Name:     "save",
          Aliases:  []string{"s","f","o"}, 
          Usage: " Directory path for saving files ",
          Value: "./",
        },
      },
      Action:   listenAction,
    },
    {
      Name:     "encrypt",
      Category: "Encryption",
      Usage:    "Encryption -f <filespath> -alg <algorithm> -k <keypath>",
      Flags:    []cli.Flag{
        &cli.StringFlag{
          Name:     "algorithm",
          Aliases:  []string{"alg"}, 
          Usage: " define the algorithm <DES|AES> ",
          Value: "DES",
        },
        &cli.StringFlag{
          Name:     "key",
          Aliases:  []string{"k"}, 
          Usage: "key location <path> ",
          Value: "key",
        },
        &cli.StringFlag{
          Name:     "file",
          Aliases:  []string{"f"}, 
          Usage: "directory or file path for ciphertext",
          Value: "key",
        },
      },
      Action:   encryptAction,
    },
    {
      Name:     "decrypt",
      Category: "Encryption",
      Usage:    "Decryption -f <filespath> -alg <algorithm> -k <keypath>",
      Flags:    []cli.Flag{
        &cli.StringFlag{
          Name:     "algorithm",
          Aliases:  []string{"alg"}, 
          Usage: " define the algorithm <DES|AES> ",
          Value: "DES",
        },
        &cli.StringFlag{
          Name:     "key",
          Aliases:  []string{"k"}, 
          Usage: "key <key> ",
          Value: "key",
        },
        &cli.StringFlag{
          Name:     "file",
          Aliases:  []string{"f"}, 
          Usage: "directory or file path for ciphertext",
          Value: "key",
        },
      },
      Action:   decryptAction,
    },
  }

  // RUN app
  if err:= app.Run(os.Args); err != nil {
    fmt.Println(err)
  }
}
