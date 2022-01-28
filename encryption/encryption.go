package encryption

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/m3dsh/encrypted-file-sender/aes"
	tm "github.com/buger/goterm"
)


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

        if err:= aes.Encrypter(path+string(os.PathSeparator)+ file.Name()); err != nil { 
          fmt.Println(err)
          continue
        }
        encrypted++
        fmt.Println(tm.Color(fmt.Sprintf(" %v", file.Name()),tm.GREEN))

      }
    } else{
      
        stat := fmt.Sprintf("Encrypting file -> %s \n", fileInfo.Name())
        fmt.Print(stat)
        if err:= aes.Encrypter(path); err != nil { 
          fmt.Println(err)
        }
        fmt.Println(tm.Color(fmt.Sprintf(" %v", fileInfo.Name()),tm.GREEN))
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

        if err:= aes.Decrypter(path+string(os.PathSeparator)+ file.Name()); err != nil { 
          fmt.Println(err)
          continue
        }
        encrypted++
        fmt.Println(tm.Color(fmt.Sprintf(" %v", file.Name()),tm.GREEN))

      }
    } else{
      
        stat := fmt.Sprintf("Decrypting file -> %s \n", fileInfo.Name())
        fmt.Print(stat)
        if err:= aes.Decrypter(path); err != nil { 
          fmt.Println(err)
        }
        fmt.Println(tm.Color(fmt.Sprintf(" %v", fileInfo.Name()),tm.GREEN))
    }
   
  }
  return nil


}


