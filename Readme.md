pOut - for hackers who like types // Protocol Buffer

pOut is a printer for Protocol Buffers. It can read from stdin or files and produce bytes printed to Stdout. These can then be piped to other programs. Using Protocol Buffers in this way enables rapid development and debugging, a faster feedback cycle, and 

Features
- Protobuf *or* JSON input

- can be piped into too, so you can do neat things like pipe a cURL response directly into pOut and into your program
- Detection of Protobuf types with standard protobuf layouts

Installation & Requirements
- Protoc - I am working on using *buf* to remove this dependency and hardcode this into the binary

`brew cask install pout` 

Background & Motivation
- It is really hard to hack with protobuf, there are so many things to compile
- type and folder structure layouts can often be confusing to beginners
- gRPC is a great tool, but because of the typed nature it is often hard to debug. It's not as easy as you think to load up a server and simply send a few commands to it. There are 'curl' options, but it becomes complicated to send messages there.
- compiling Protoc Buffers can break flow, especially if you want to test data types
- you can u se pOut with gRPC curl
- pOut is designed to take slightly longer to process a file, but make it super simple to get going. 
- lots of proto files means lots of cognitive overhead
- The new protocol buffers Go layout enabling easy reflection https://blog.golang.org/protobuf-apiv2
- UNIX philosophy of extremely small programs. pOut in combination with pIn allow you to chain many commands together.

Usage
`pout [flags] `
Flags:
Resolving the message:
There are 3 variables that any system needs to resolve a proto message:
`-I A proto base path`
`-f A proto file & package`
only -n is required 

Args:
`[message name] A message name`
Input
`"--" will take the input from stdin`
`[filename] - will load a json or .pb file in`
`-`

pOut will find the first message name in any of the packages you supply. 

Examples
NodeJS
Python
Go
cURL

check the /examples folder for more!

Notes

Thank you to the excellent GoReleaser package, which allows releases, including protoc builds, to be automated for lots of operating systems.