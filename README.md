# Installation
- `cd cmd`
- `go run cmd.go`

# Features 
- [x] CLI interface with usage
- [x] Sending data over TCP connection
- [x] Sending file over TCP connection
- [x] Encrypt file with AES or DES and share key
- [x] Show encryption and decryption time 
- [x] Generate some random files with different size
- [ ] Make CSV file for timing 
- [ ] Openssl
- [ ] Plot the CSV file (optional)
- [ ] Windows compatibility (optional) 
- [ ] Unit testing for app (optional)
- [ ] Generate wiki with Godoc (optional)

# Usage

```
NAME:
   encrypted-file-sender - send encrypted file over network

USAGE:
   cmd [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command
   Encryption:
     encrypt  Encryption -f <filespath> -alg <algorithm> -k <keypath>
     decrypt  Decryption -f <filespath> -alg <algorithm> -k <keypath>
   Markov Chain Generator:
     filegen  filegen -minwords <int> -maxwords <int> -count <int>
   Network:
     connect  Connect [-p <PORT>] <IP>
     listen   listen [-p <PORT>]

GLOBAL OPTIONS:
   --verbose, -V  verbose output (default: false)
   --help, -h     show help (default: false)
   --version, -v  output version information (default: false)
```
