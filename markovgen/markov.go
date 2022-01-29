package markov

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	markov "github.com/MagnusFrater/markov-chain-text-generator"
	tm "github.com/buger/goterm"
)

func MarkovGen(minW, maxW, count, preLen, sufLen int,
dataPath,savePath string) error{

		buf, err := ioutil.ReadFile(dataPath)
		if err != nil {
			log.Fatal(err)
		}

    corpus := string(buf)
	generator := markov.New(preLen, sufLen)
	generator.Add(corpus)
  step := (maxW - minW) /count
  words := minW
  fmt.Print("--------------------------------------\n"+
  "Start Generating files:\n" +
            "--------------------------------------\n")
  _ = os.Mkdir(savePath, os.ModePerm)
  for i :=0 ; i < 10 ; i++ {
     // Use os.Create to create a file for writing.
    f, err := os.Create(savePath + string(os.PathSeparator) +
    "text_" + strconv.Itoa(i)+ ".txt")
    if err != nil{
      fmt.Println(err)
    }
        stat := fmt.Sprintf("Generate file with %v words -> %s %v/%v \n",words,f.Name() , i ,count)
        fmt.Print(stat)
    
    // Create a new writer.
    w := bufio.NewWriter(f)
    
    // Write a string to the file.
    w.WriteString(generator.Generate(words))
    
    // Flush.
    w.Flush()
        fmt.Println(tm.Color(fmt.Sprintf(" %v", f.Name()),tm.GREEN))
    f.Close()
   words += step 
  } 
  return nil
}

