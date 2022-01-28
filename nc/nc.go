// Package nc for handling networking stuff.
package nc

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tm "github.com/buger/goterm"
)


// serveFileReciever make Http listener for receiving file
func serveFileReciever() {
  receiveFileHandler := http.HandlerFunc(receiveFile)
  http.Handle("/file", receiveFileHandler)
  http.ListenAndServe(":8080", nil)

}

// receiveFile is a handler for recieved file 
func receiveFile(w http.ResponseWriter, request *http.Request) {
  err := request.ParseMultipartForm(200 << 20) //maxMemory 200MB
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  fmt.Printf("key: %v \n",request.Form["key"])
  file, h, err := request.FormFile("file")
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  tmpfile, err := os.Create("./" + h.Filename)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  defer tmpfile.Close()
  _, err = io.Copy(tmpfile, file)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  fmt.Printf("file: %s\nsize: %v\n", 
  h.Filename, h.Size)

  w.WriteHeader(200)
}


func sendFiles(path string, n *NetData, con net.Conn) error{
  fileInfo, err := os.Stat(path) 
  if err != nil{
    return err  
  }
  if fileInfo.IsDir() {
    files, err := ioutil.ReadDir(path)
    if err != nil {
      return err
    }
    sent := 0
    for i, file := range files {
      stat := fmt.Sprintf("Sending file -> %s %v/%v \n", file.Name(), i , len(files))
      fmt.Print(stat)
      WriteStatus(stat, tm.BLUE, &con);

      if err:= sendFileHTTP(path +"/"+file.Name(),fmt.Sprintf("http://%v:8080/file",n.Host), "POST"); err != nil{
        fmt.Println(err)
        continue
      }
      sent++
      fmt.Println(tm.Color(fmt.Sprintf(" %v", file.Name()),tm.GREEN))

    }
    stat := fmt.Sprintf("%v files out of %v files received \n", sent, len(files))
    WriteStatus(stat, tm.GREEN, &con)
    
  } else{
    
      stat := fmt.Sprintf("Sending file -> %s \n", fileInfo.Name())
      fmt.Print(stat)
      WriteStatus(stat, tm.BLUE, &con);
      if err:= sendFileHTTP(path,fmt.Sprintf("http://%v:8080/file",n.Host), "POST"); err != nil{
        fmt.Println(err)
      }
      fmt.Println(tm.Color(fmt.Sprintf(" %v", fileInfo.Name()),tm.GREEN))
  }
   
  return nil
}
// sendFileHTTP make http multipart request for sending a file 
func sendFileHTTP(filePath, urlPath, method string) error{
  client := &http.Client{
        Timeout: time.Second * 10,
    }
    r, w := io.Pipe()
    writer := multipart.NewWriter(w)
    go func() {
      defer w.Close()
      defer writer.Close()
    fw, err := writer.CreateFormField("key")
    if err != nil {
      fmt.Println(err)
    }
    _, err = io.Copy(fw, strings.NewReader("<ThisIsTheKey>"))
    if err != nil {
      fmt.Println(err)
    }

    fw, err = writer.CreateFormFile("file",filePath)
    if err != nil {
      fmt.Println(err)
    }
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    if _, err = io.Copy(fw, file); err != nil{
      return
    }
    }()
    req, err := http.NewRequest(method, urlPath, r)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    rsp, _ := client.Do(req)
    if rsp.StatusCode != http.StatusOK {
        log.Printf("Request failed with response code: %d", rsp.StatusCode)
    }
    return nil
}

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
        fmt.Println(": finished reading string")
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

// PrintStatus Read from a reader and print it every second
func WriteStatus(str string, color int, con *net.Conn) {
      fmt.Fprint(*con, tm.Color(str,color))
}
func ReadStatus(con *net.Conn) {
    buf:= make([]byte, 1024)
    defer func() {
        (*con).Close()
        log.Printf("Connection from %v is closed\n", (*con).RemoteAddr())
    }()
  for {
    nbytes, err := (*con).Read(buf) 
      if err != nil {
            if err != io.EOF {
                fmt.Println("read error:", err)
            }
            break
        }
        fmt.Println(string(buf[0:nbytes]))
    }
}

// NetData all data required by nettool to work.
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

    go func(){
    for {
    listener, err := net.Listen("tcp", n.Port)
    if err != nil {
      log.Fatalln(err)
    }
    fmt.Println("Listening on", n.Port)
      con, err := listener.Accept()
      if err != nil {
        log.Fatalln(err)
      }
      fmt.Println("Connect from", con.RemoteAddr())
      ReadStatus(&con)
    }
  }()
      serveFileReciever()
  }else{ // Client-side handler


    con,err := net.Dial("tcp", n.Host + n.Port)
    if err!= nil {
      log.Fatalln(err)
    }
    fmt.Println(tm.Color(fmt.Sprintln("Connected to", n.Host+n.Port), tm.BLUE))
    fmt.Println("Enter the path for sending Files: ")
    var path string
    fmt.Scanln(&path)
    sendFiles(strings.TrimSpace(path), n, con)

  }

  return nil
}
