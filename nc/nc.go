// Package nc for handling networking stuff.
package nc

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

// tcp_con_handle handle TCP connection
// TCP -> Stdout and Stdin -> TCP.
func tcp_con_handle(con net.Conn) {
  chan_to_stdout := stream_copy(con, os.Stdout)
  chan_to_remote := stream_copy(os.Stdin, con)
  select {
  case <-chan_to_stdout:
    log.Println("Remote connection is closed")
  case <-chan_to_remote:
    log.Println("Local program is terminated")
  }
}

// stream_copy copy operation between streams: os and tcp streams.
func stream_copy(src io.Reader, dst io.Writer) <-chan int {
  buf := make([]byte, 1024)
  sync_channel := make(chan int)
  go func() {
    defer func() {
      if con, ok := dst.(net.Conn); ok {
        con.Close()
        log.Printf("Connection from %v is closed\n", con.RemoteAddr())
      }
      sync_channel <- 0 // Notify that processing is finished
    }()
    for {
      var nBytes int
      var err error
      nBytes, err = src.Read(buf)
      if err != nil {
        if err != io.EOF {
          log.Printf("Read error: %s\n", err)
        }
        break
      }
      _, err = dst.Write(buf[0:nBytes])
      if err != nil {
        log.Fatalf("Write error: %s\n", err)
      }
    }
  }()
  return sync_channel
}

// NetData all data required by nettool to work
type NetData struct{
  IsListen bool
  Host string
  Port string
}


func Run(n *NetData) error{

  if  n.Port!= "" {
    if _, err := strconv.Atoi(n.Port); err != nil {
      log.Println("Port shall not be empty and has integer value")
      os.Exit(1)
    }
  }
  n.Port = ":"+ n.Port
  // Server-side handler
  if n.IsListen {
    listener, err := net.Listen("tcp", n.Port)
    if err != nil {
      log.Fatalln(err)
    }
    log.Println("Listening on", n.Port)
      con, err := listener.Accept()
      if err != nil {
        log.Fatalln(err)
      }
      log.Println("Connect from", con.RemoteAddr())
      tcp_con_handle(con)
  }else{ //Client-side handler
    con,err := net.Dial("tcp", n.Host + n.Port)
    if err!= nil {
      log.Fatalln(err)
    }
    log.Println("Connected to", n.Host+n.Port)
    tcp_con_handle(con)

  }

  return nil
}
