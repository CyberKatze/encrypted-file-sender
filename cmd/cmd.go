// Package cmd is command-line interface for encrypted file sender.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m3dsh/encrypted-file-sender/encryption"
	markov "github.com/m3dsh/encrypted-file-sender/markovgen"
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
     Encryption: c.Bool("encryption"),
     Path: c.String("file"),
   }

   err:= nc.Run(n)
  if err != nil {
    log.Println(err)
  }
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
    Encryption: c.Bool("encryption"),
    Path: c.String("save"),
  }
  err := nc.Run(n)
  if err != nil {
    log.Println(err)
  }
  return nil
  
}

func encryptAction(c *cli.Context) error{
  err := encryption.EncryptFiles(c.String("file"),c.String("algorithm"),c.String("key"))
  if err != nil {
    log.Println(err)
  }
  return nil
}
func decryptAction(c *cli.Context) error{
  err := encryption.DecryptFiles(c.String("file"),c.String("algorithm"),c.String("key"))
  if err != nil {
    log.Println(err)
  }
  return nil
}

func filegenAction(c *cli.Context) error{
//func MarkovGen(minW, maxW, count, preLen, sufLen int,
//dataPath,savePath string) error{
  err := markov.MarkovGen(c.Int("minwords"),c.Int("maxwords"),c.Int("count"),
  c.Int("prefixlength"), c.Int("suffixlength"), c.String("data"),
  c.String("save"))
  if err != nil {
    log.Println(err)
  }
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
        &cli.BoolFlag{
          Name:     "encryption",
          Aliases:  []string{"e","enc"}, 
          Usage: "encrypt data before sending them out",
          Value: false,
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
        &cli.BoolFlag{
          Name:     "encryption",
          Aliases:  []string{"e","enc"}, 
          Usage: "decrypt after saving them",
          Value: false,
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
          Usage: "directory or file path for plaintext",
          Value: "files",
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
          Value: "files",
        },
      },
      Action:   decryptAction,
    },
    {
      Name:     "filegen",
      Category: "Markov Chain Generator",
      Usage:    "filegen -minwords <int> -maxwords <int> -count <int>",
      Flags:    []cli.Flag{
        &cli.IntFlag{
          Name:     "minwords",
          Aliases:  []string{"min"}, 
          Usage: "minumum words count for text files",
          Value: 1000000,
        },
        &cli.IntFlag{
          Name:     "maxwords",
          Aliases:  []string{"max"}, 
          Usage: "maximum words count for text files",
          Value: 20000000,
        },
        &cli.IntFlag{
          Name:     "count",
          Aliases:  []string{"c"}, 
          Usage: "the count of text files",
          Value: 10,
        },
        &cli.IntFlag{
          Name:     "prefixlength",
          Aliases:  []string{"pl"}, 
          Usage: "length of the Markov Chain's prefixes",
          Value: 2,
        },
        &cli.IntFlag{
          Name:     "suffixlength",
          Aliases:  []string{"sl"}, 
          Usage: "length of the Markov Chain's suffixlength",
          Value: 2,
        },
        &cli.StringFlag{
          Name:     "data",
          Aliases:  []string{"d"}, 
          Usage: "path to data file for generator",
          Value: "sample_data.txt",
        },
        &cli.StringFlag{
          Name:     "save",
          Aliases:  []string{"o", "s"}, 
          Usage: "path for saving generated files",
          Value: "files",
        },
      },
      Action:   filegenAction,
    },
  }

  // RUN app
  if err:= app.Run(os.Args); err != nil {
    fmt.Println(err)
  }
}
