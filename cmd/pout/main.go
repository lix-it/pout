package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lix-it/pout/internal/proto/message"
	"github.com/lix-it/pout/internal/proto/registry"
	"github.com/lix-it/pout/pkg/pout"
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

	protoPackage, msgName := pout.SplitIdentifier(flag.Arg(0))

	config := Config{
		Verbose: *verbose,
	}

	if config.Verbose {
		fmt.Println("proto root:", *protoPath)
		fmt.Println("proto package:", protoPackage)
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
	regy, err := registry.BuildProtoRegistry(config.Verbose, *protoPath)
	if err != nil {
		panic(fmt.Errorf("error BuildProtoRegistry(): %w", err))
	}

	msg, err := message.FromJSON(config.Verbose, regy, protoPackage, protoreflect.Name(msgName), r)
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
