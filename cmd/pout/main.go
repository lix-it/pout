package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lix-it/pout/internal/proto/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var protoPath = flag.String("I", "./protos", "path to proto base folder")
var verbose = flag.Bool("v", false, "verbose mode. display all messages")

func main() {
	flag.Parse()
	if *protoPath == "" {
		panic("you must enter a valid proto base path")
	}

	pkgMsg := strings.Split(flag.Arg(0), ":")
	protoFile := filepath.Clean(pkgMsg[0])
	msgName := pkgMsg[1]

	config := Config{
		Verbose: *verbose,
	}

	if config.Verbose {
		fmt.Println("proto path:", path.Join(*protoPath, protoFile))
		fmt.Println("message name:", msgName)
	}
	// use file or stdin
	var useStdin bool
	// -- isn't parsed by flag.Args()
	if os.Args[len(os.Args)-1] == "--" {
		useStdin = true
	}
	var r io.ReadCloser
	r = os.Stdin
	if !useStdin {
		var err error
		r, err = os.Open(flag.Arg(1))
		if err != nil {
			panic(err)
		}
	}
	defer r.Close()

	msg, err := registry.FromJSON(*protoPath, protoFile, protoreflect.Name(msgName), r)
	if err != nil {
		wrapErr := fmt.Errorf("error converting JSON to proto: %w; check whether the message paths and types are correct", err)
		panic(wrapErr)
	}
	b, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", b)
}
