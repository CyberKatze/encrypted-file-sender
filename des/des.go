package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)
const (
  ENCRYPT = iota
  DECRYPT
)

var initialPermutation []int = []int{
  58, 50, 42, 34, 26, 18, 10, 2, 60, 52, 44, 36, 28, 20, 12, 4,
  62, 54, 46, 38, 30, 22, 14, 6, 64, 56, 48, 40, 32, 24, 16, 8,
  57, 49, 41, 33, 25, 17, 9, 1, 59, 51, 43, 35, 27, 19, 11, 3,
  61, 53, 45, 37, 29, 21, 13, 5, 63, 55, 47, 39, 31, 23, 15, 7,
}

var finalPermutation []int = []int{
  40, 8, 48, 16, 56, 24, 64, 32, 39, 7, 47, 15, 55, 23, 63, 31,
  38, 6, 46, 14, 54, 22, 62, 30, 37, 5, 45, 13, 53, 21, 61, 29,
  36, 4, 44, 12, 52, 20, 60, 28, 35, 3, 43, 11, 51, 19, 59, 27,
  34, 2, 42, 10, 50, 18, 58, 26, 33, 1, 41, 9, 49, 17, 57, 25,
}

var keyPermutation []int = []int{
  57, 49, 41, 33, 25, 17, 9, 1, 58, 50, 42, 34, 26, 18,
  10, 2, 59, 51, 43, 35, 27, 19, 11, 3, 60, 52, 44, 36,
  63, 55, 47, 39, 31, 23, 15, 7, 62, 54, 46, 38, 30, 22,
  14, 6, 61, 53, 45, 37, 29, 21, 13, 5, 28, 20, 12, 4,
}

var keyShiftsPerRound []int = []int{
  1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1,
}

var compressionPermutation []int = []int{
  14, 17, 11, 24, 1, 5, 3, 28, 15, 6, 21, 10,
  23, 19, 12, 4, 26, 8, 16, 7, 27, 20, 13, 2,
  41, 52, 31, 37, 47, 55, 30, 40, 51, 45, 33, 48,
  44, 49, 39, 56, 34, 53, 46, 42, 50, 36, 29, 32,
}

var expansionPermutation []int = []int{
  32, 1, 2, 3, 4, 5, 4, 5, 6, 7, 8, 9,
  8, 9, 10, 11, 12, 13, 12, 13, 14, 15, 16, 17,
  16, 17, 18, 19, 20, 21, 20, 21, 22, 23, 24, 25,
  24, 25, 26, 27, 28, 29, 28, 29, 30, 31, 32, 1,
}

var sboxes [][]int = [][]int{
  {
    14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7,
    0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8,
    4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0,
    15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13,
  },
  {
    15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10,
    3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5,
    0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15,
    13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9,
  },
  {
    10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8,
    13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1,
    13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7,
    1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12,
  },
  {
    7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15,
    13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9,
    10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4,
    3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14,
  },
  {
    2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9,
    14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6,
    4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14,
    11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3,
  },
  {
    12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11,
    10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8,
    9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6,
    4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13,
  },
  {
    4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1,
    13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6,
    1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2,
    6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12,
  },
  {
    13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7,
    1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2,
    7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8,
    2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11,
  },
}

var pBoxPermutation []int = []int{
  16, 7, 20, 21, 29, 12, 28, 17, 1, 15, 23, 26, 5, 18, 31, 10,
  2, 8, 24, 14, 32, 27, 3, 9, 19, 13, 30, 6, 22, 11, 4, 25,
}

func getBit(block []byte, pos int) (result byte) {
  posByte := pos >> 3
  // pos mod 8 = pos & (pow(2, 3) - 1)
  bitPos := pos & 7
  result = 1 & (block[posByte] >> (7 - uint(bitPos)))
  return result
}

func permutate(block []byte, permutationMatrix []int) (result []byte) {
  result = make([]byte, len(permutationMatrix)>>3)
  for i := range permutationMatrix {
    var bit byte = getBit(block, permutationMatrix[i]-1)
    resByteIndex := i >> 3
    // i mod 8 = i & (pow(2, 3) - 1)
    bitPos := i & 7
    result[resByteIndex] |= (bit << (7 - uint(bitPos)))
  }
  return result
}

func key64To56(key []byte) (result []byte) {
  result = permutate(key, keyPermutation)
  return result
}

func leftRoundShift(block []byte, shift int, beginOffset int, endOffset int) (result []byte) {
  result = make([]byte, len(block))
  copy(result, block)
  if endOffset > 0 {
    result[len(block)-1] = result[len(block)-1] & (0xff << uint(endOffset))
  }
  for i := 0; i < shift; i++ {
    var leading byte = 0x00
    for j := len(result) - 1; j >= 0; j-- {
      var nextLeading byte = (result[j] & 0x80) >> 7
      result[j] = (result[j] << 1) | leading
      leading = nextLeading
    }
    if endOffset > 0 {
      result[len(result)-1] = result[len(result)-1] | (leading << uint(endOffset))
    } else {
      result[len(result)-1] = result[len(result)-1] | leading
    }
  }
  return result
}

func joinKeys(left []byte, right []byte) (result []byte) {
  result = make([]byte, 7)
  for i := 0; i < 3; i++ {
    result[i] = left[i]
  }
  result[3] = (left[3] & 0xf0) | (0x0f & right[0])
  for i := 1; i < 4; i++ {
    result[3+i] = right[i]
  }
  return result
}

func getRoundSubkey(key []byte, round int) (result []byte) {
  var left []byte = make([]byte, 4)
  var right []byte = make([]byte, 4)
  copy(left, key[0:4])
  copy(right, key[3:])
  left = leftRoundShift(left, keyShiftsPerRound[round], 0, 4)
  right = leftRoundShift(right, keyShiftsPerRound[round], 4, 0)
  shiftedKey := joinKeys(left, right)
  result = permutate(shiftedKey, compressionPermutation)
  return result
}

func xorBlocks(block1 []byte, block2 []byte) (result []byte) {
  result = make([]byte, len(block1))
  for i := range block1 {
    result[i] = block1[i] ^ block2[i]
  }
  return result
}

func getSBoxValue(block []byte, step int) (result byte) {
  beginPos := step * 6
  endPos := beginPos + 5
  var col byte
  var row byte

  row = (getBit(block, beginPos) << 1) | (getBit(block, endPos))
  col =
    (getBit(block, beginPos+1) << 3) |
      (getBit(block, beginPos+2) << 2) |
      (getBit(block, beginPos+3) << 1) |
      getBit(block, beginPos+4)
  result = byte(sboxes[step][row*16+col])
  return result
}

func applySBoxes(block []byte) (result []byte) {
  result = make([]byte, 4)
  for i := 0; i < 8; i += 2 {
    leftNibble := getSBoxValue(block, i)
    rightNibble := getSBoxValue(block, i+1)
    result[i>>1] = leftNibble<<4 | rightNibble
  }
  return result
}

func feistelFunction(key []byte, block []byte) (result []byte) {
  expandedBlock := permutate(block, expansionPermutation)
  xoredBlock := xorBlocks(expandedBlock, key)
  substitutedBlock := applySBoxes(xoredBlock)
  result = permutate(substitutedBlock, pBoxPermutation)
  return result
}

func cipherStep(block []byte, roundKeys [][]byte) (result []byte) {
  block = permutate(block, initialPermutation)
  left := block[0:4]
  right := block[4:]
  for k := range roundKeys {
    lastRight := right
    transformedRight := feistelFunction(roundKeys[k], right)
    right = xorBlocks(left, transformedRight)
    left = lastRight
  }
  result = permutate(append(right, left...), finalPermutation)
  return result
}

// GenerateRoundsKeys generates the key for the cipher
// it receives the keys (one or three keys)
// operation can be des.ENCRYPT or des.DECRYPT
// this keys should be used in CipherBlock and the operation param is what decides if a encrypt or decrypt will be
// processed by the CipherBlock function
func GenerateRoundsKeys(keys [][]byte, operation int) (result [][][]byte) {
  keyCount := 1
  if len(keys) == 3 {
    keyCount = 3
  }
  result = make([][][]byte, keyCount)
  for k := 0; k < keyCount; k++ {
    result[k] = make([][]byte, 16)
    var key []byte
    // Invert key order in decryption
    if operation == DECRYPT {
      key = keys[keyCount-k-1]
    } else {
      key = keys[k]
    }
    transformedKey := key64To56(key)
    // Encription and k is even OR Decription and k is odd
    if (operation == ENCRYPT && k & 1 == 0) || (operation == DECRYPT && k & 1 != 0) {
      // Encryption round keys
      for i := 0; i < 16; i++ {
        result[k][i] = getRoundSubkey(transformedKey, i)
      }
    } else {
      // Decryption round keys
      for i := 0; i < 16; i++ {
        result[k][i] = getRoundSubkey(transformedKey, 15-i)
      }
    }
  }
  return result
}

// CipherBlock cipher a block using the roundKeys
// roundKeys must be the return of GenerateRoundsKeys
// If the length of roundKeys is 3 then Triple DES is applied, if not, only the first key is used
func CipherBlock(block []byte, roundKeys [][][]byte) (result []byte) {
  if len(roundKeys) == 3 {
    result = block
    for i := 0; i < 3; i++ {
      result = cipherStep(result, roundKeys[i])
    }
  } else {
    result = cipherStep(block, roundKeys[0])
  }
  return result
}

const BLOCK_SIZE = 8

var inputReader *bufio.Reader
var inputFile *os.File
var outputWriter *bufio.Writer
var outputFile *os.File
var inputParam string
var outputParam string
var keyString string
var operation int
var tripleDes bool = false

func askForKey() string {
  reader := bufio.NewReader(os.Stdin)
  fmt.Fprintf(os.Stderr, "Enter key: ")
  key, _ := reader.ReadString('\n')
  return strings.Trim(key, "\n")
}

func readArgs() {
  var encrypt, decrypt, newKey, help bool
  flag.BoolVar(&encrypt, "encrypt", false, "Encrypt input")
  flag.BoolVar(&decrypt, "decrypt", false, "Decrypt input")
  flag.BoolVar(&newKey, "newkey", false, "Outputs a new key")
  flag.BoolVar(&help, "h", false, "Display help and exit")
  flag.BoolVar(&tripleDes, "3des", false, "Encrypt/decrypt using Triple DES")
  flag.StringVar(&inputParam, "i", "", "Input file to encrypt/decrypt. Use \"-\" to standard input. Default \"-\"")
  flag.StringVar(&outputParam, "o", "", "Output file. Use \"-\" to standard output. Default \"-\"")
  flag.StringVar(&keyString, "k", "", "Cipher key")
  flag.Parse()
  if help {
    flag.Usage()
    os.Exit(0)
  }
  if newKey {
    keys := prepareKey()
    fmt.Println(keysToString(keys))
    os.Exit(0)
  }
  if !encrypt && !decrypt || encrypt {
    operation = ENCRYPT
  } else {
    operation = DECRYPT
  }
  if encrypt && decrypt {
    fmt.Fprintln(os.Stderr, "Choose between only one of -encrypt or -decrypt operations")
    os.Exit(1)
  }
  if inputParam == "" {
    inputParam = "-"
  }
  if outputParam == "" {
    outputParam = "-"
  }
  if keyString == "" && operation == DECRYPT {
    // Request key from user if input is not standard input
    if inputParam != "-" {
      keyString = askForKey()
    }
    if keyString == "" {
      fmt.Fprintln(os.Stderr, "** Cipher key is required for operation decrypt")
      flag.Usage()
      os.Exit(1)
    }
  }
}

func validateArgs() {
  // Input is readable
  if inputParam != "-" {
    if _, err := os.Stat(inputParam); os.IsNotExist(err) {
      fmt.Fprintln(os.Stderr, "Error: Input file not found")
    }
  }
}

func prepareKey() (result [][]byte) {
  keyCount := 1
  if tripleDes {
    keyCount = 3
  }
  if keyString == "" {
    result = make([][]byte, keyCount)
    for i := 0; i < keyCount; i++ {
      result[i] = make([]byte, BLOCK_SIZE)
      _, err := rand.Read(result[i])
      if err != nil {
        fmt.Fprintln(os.Stderr, "Error generating key")
        os.Exit(1)
      }
    }
  } else {
    splitedKey := strings.Split(keyString, ":")
    if len(splitedKey) < keyCount {
      fmt.Fprintln(os.Stderr, "Triple DES requires three keys")
      os.Exit(1)
    }
    if len(splitedKey) == 3 {
      tripleDes = true
      keyCount = 3
    }
    if len(splitedKey) != keyCount {
      fmt.Fprintln(os.Stderr, "Invalid key count. Allowed one or three keys")
      os.Exit(1)
    }
    result = make([][]byte, keyCount)
    for k := range splitedKey {
      key := fmt.Sprintf("%016s", splitedKey[k])
      key = key[0:16]
      var err error
      result[k], err = hex.DecodeString(key)
      if err != nil {
        fmt.Fprintf(os.Stderr, "Key %s is not a valid hexadecimal value\n", key)
        os.Exit(1)
      }
    }
  }
  return result
}

func blockToHexString(block []byte) string {
  return hex.EncodeToString(block)
}

func keysToString(keys [][]byte) string {
  res := make([]string, len(keys))
  for k := range keys {
    res[k] = blockToHexString(keys[k])
  }
  return strings.Join(res, ":")
}

func checkInputReader() {
  if inputReader != nil {
    return
  }
  if inputParam != "-" {
    file, err := os.Open(inputParam)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error opening input file: %s\n", err.Error())
      os.Exit(1)
    }
    inputFile = file
    inputReader = bufio.NewReader(file)
  } else {
    inputReader = bufio.NewReader(os.Stdin)
  }
}

func getNextInputBlock(block []byte) (end bool, size int, pad int) {
  end = false
  pad = 0
  checkInputReader()
  size, err := inputReader.Read(block)
  if err == io.EOF {
    end = true
  } else if err != nil {
    fmt.Fprintf(os.Stderr, "Error reading input file: %s\n", err.Error())
    os.Exit(1)
  }
  if size > 0 && size != BLOCK_SIZE {
    if operation == DECRYPT {
      pad = int(block[0])
      size = 0
    } else {
      pad = BLOCK_SIZE - size
    }
  }
  return end, size, pad
}

func checkOutputWriter() {
  if outputWriter != nil {
    return
  }
  if outputParam != "-" {
    file, err := os.OpenFile(outputParam, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error opening output file: %s\n", err.Error())
      os.Exit(1)
    }
    outputFile = file
    outputWriter = bufio.NewWriter(file)
  } else {
    outputWriter = bufio.NewWriter(os.Stdout)
  }
}

func writeToOutput(block []byte, pad int) {
  checkOutputWriter()
  var err error
  if pad == 0 {
    _, err = outputWriter.Write(block)
  } else {
    if operation == ENCRYPT {
      buff := make([]byte, BLOCK_SIZE+1)
      copy(buff, block)
      buff[len(buff)-1] = byte(pad)
      _, err = outputWriter.Write(buff)
    } else {
      buff := make([]byte, BLOCK_SIZE-pad)
      copy(buff, block)
      _, err = outputWriter.Write(buff)
    }
  }
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error writing output: %s\n", err.Error())
    os.Exit(1)
  }
  if outputParam == "-" {
    outputWriter.Flush()
  }
}

func closeFiles() {
  if outputWriter != nil {
    outputWriter.Flush()
  }
  if inputFile != nil {
    inputFile.Close()
  }
  if outputFile != nil {
    outputFile.Close()
  }
}

func encrypt(keys [][]byte) {
  roundKeys := GenerateRoundsKeys(keys, ENCRYPT)
  end := false
  size := 0
  pad := 0
  block := make([]byte, BLOCK_SIZE)
  for !end {
    end, size, pad = getNextInputBlock(block)
    if size > 0 {
      processedBlock := CipherBlock(block, roundKeys)
      writeToOutput(processedBlock, pad)
    }
  }
}

func decrypt(keys [][]byte) {
  end := false
  size := 0
  pad := 0
  roundKeys := GenerateRoundsKeys(keys, DECRYPT)
  var pendingBlock bool = false
  block := make([]byte, BLOCK_SIZE)
  var processedBlock []byte
  for !end {
    end, size, pad = getNextInputBlock(block)
    if pendingBlock {
      writeToOutput(processedBlock, pad)
      pendingBlock = false
    }
    if size > 0 {
      processedBlock = CipherBlock(block, roundKeys)
      pendingBlock = true
    }
  }
  if pendingBlock {
    writeToOutput(processedBlock, pad)
  }
}

func filenameObfuscator(path string) string {
  filenameArr := strings.Split(path, string(os.PathSeparator))
  filename := filenameArr[len(filenameArr)-1]
  path2 := strings.Join(filenameArr[:len(filenameArr)-1], string(os.PathSeparator))

  return path2 + string(os.PathSeparator) + filename + Ext

}
func filenameDeobfuscator(path string) string {
  //get the path for the output
  r := regexp.MustCompile(`^(.*)\`+Ext +`$`)
  opPath :=r.FindStringSubmatch(path)[1]
  // Divide filepath
  filenameArr := strings.Split(opPath, string(os.PathSeparator))
  //Get  filename
  filename := filenameArr[len(filenameArr)-1]
  // get parent dir
  path2 := strings.Join(filenameArr[:len(filenameArr)-1], string(os.PathSeparator))
  return path2 + string(os.PathSeparator) + filename
}
var Ext = ".enc"
func DesDecryption(path, key string){
  defer closeFiles()
  inputParam = path
  outputParam = filenameDeobfuscator(path)
  keyString = key
  keys := prepareKey()
  if keyString == "" {
    fmt.Fprintf(os.Stderr, "Using key: %s\n", keysToString(keys))
  }
    decrypt(keys)
}

func DesEncryption(path, key string){
  defer closeFiles()
  inputParam = path
  outputParam = filenameObfuscator(path)
  keyString = key
  keys := prepareKey()
  if keyString == "" {
    fmt.Fprintf(os.Stderr, "Using key: %s\n", keysToString(keys))
  }
    encrypt(keys)
}

func main() {
  DesDecryption("myfiles/myfile.enc", "db1e839774ccecdb")
  

}
