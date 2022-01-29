package encryption

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	tm "github.com/buger/goterm"
	"github.com/m3dsh/encrypted-file-sender/aes"
)

type data struct{
    fileName  string
    size      int
    time      time.Duration
}

var encryptionData []data
var decryptionData []data

func EncryptFiles(path ,enc, key string) error {
 fileInfo, err := os.Stat(path) 
  aes.InitializeBlock([]byte("a very very very very secret key"))
  if err != nil{
    return err  
  }
  fmt.Print("--------------------------------------\n"+
  "Start Encrypting:\n" +
            "--------------------------------------\n")
  if enc == "AES"{
    if fileInfo.IsDir() {
      files, err := ioutil.ReadDir(path)
      if err != nil {
        return err
      }
      encrypted := 0
      for i, file := range files {
        stat := fmt.Sprintf("Encrypting file -> %s %v/%v \n", file.Name(), i , len(files))
        fmt.Print(stat)
        start :=time.Now()
        if err:= aes.Encrypter(path+string(os.PathSeparator)+ file.Name()); err != nil { 
          fmt.Println(err)
          continue
        }
         elapsed:= time.Since(start)
        encrypted++
        fmt.Println(tm.Color(fmt.Sprintf(" %v", file.Name()),tm.GREEN))
        encryptionData = append(encryptionData,
        data{fileName: file.Name(),
         size: int(file.Size()),
         time: elapsed})

      }
    } else{
      
        stat := fmt.Sprintf("Encrypting file -> %s \n", fileInfo.Name())
        fmt.Print(stat)
        start := time.Now()
        if err:= aes.Encrypter(path); err != nil { 
          fmt.Println(err)
        }
        elapsed := time.Since(start)
        fmt.Println(tm.Color(fmt.Sprintf(" %v", fileInfo.Name()),tm.GREEN))
        encryptionData = append(encryptionData,
        data{fileName: fileInfo.Name(),
         size: int(fileInfo.Size()),
         time: elapsed})
    }
    for _, x := range encryptionData {
      fmt.Printf("%s -> size: %v, time: %v\n",x.fileName,x.size,x.time)
    }
   
  }
  return nil


}

func DecryptFiles(path ,enc, key string) error {
 fileInfo, err := os.Stat(path) 
  aes.InitializeBlock([]byte("a very very very very secret key"))
  if err != nil{
    return err  
  }
  fmt.Print("--------------------------------------\n"+
  "Start Decrypting:\n" +
            "--------------------------------------\n")
  if enc == "AES"{
    if fileInfo.IsDir() {
      files, err := ioutil.ReadDir(path)
      if err != nil {
        return err
      }
      encrypted := 0
      for i, file := range files {
        stat := fmt.Sprintf("Decrypting file -> %s %v/%v \n", file.Name(), i , len(files))
        fmt.Print(stat)
        start := time.Now()
        if err:= aes.Decrypter(path+string(os.PathSeparator)+ file.Name()); err != nil { 
          fmt.Println(err)
          continue
        }
        elapsed := time.Since(start)
        encrypted++
        fmt.Println(tm.Color(fmt.Sprintf(" %v", file.Name()),tm.GREEN))
        decryptionData = append(decryptionData,
        data{fileName: file.Name(),
         size: int(file.Size()),
         time: elapsed})

      }
    } else{
      
        stat := fmt.Sprintf("Decrypting file -> %s \n", fileInfo.Name())
        fmt.Print(stat)
        start := time.Now()
        if err:= aes.Decrypter(path); err != nil { 
          fmt.Println(err)
        }
        elapsed := time.Since(start)
        fmt.Println(tm.Color(fmt.Sprintf(" %v", fileInfo.Name()),tm.GREEN))
        decryptionData = append(decryptionData,
        data{fileName: fileInfo.Name(),
         size: int(fileInfo.Size()),
         time: elapsed})
    }
    for _, x := range decryptionData {
      fmt.Printf("%s -> size: %v, time: %v\n",x.fileName,x.size,x.time)
    }
   
  }
  return nil


}


