package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lix-it/pout/internal/proto/message"
	"github.com/lix-it/pout/internal/proto/registry"
	"github.com/lix-it/pout/pkg/pout"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var protoPath = flag.String("I", "./protos", "path to proto base folder")
var verbose = flag.Bool("v,verbose", false, "verbose mode. Display messages relating to proto registry and messages")
var debug = flag.Bool("d,debug", false, "debug mode. Display debug messages")
var help = flag.Bool("h,help", false, "help")

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	// set up logger
	log.SetPrefix("pout: ")
	loggingFlags := 0
	if *debug {
		loggingFlags = log.Llongfile
	}
	var start time.Time
	if *verbose {
		start = time.Now()
	}
	log.SetFlags(loggingFlags)
	if *protoPath == "" {
		log.Fatal("you must enter a valid proto base path")
	}

	protoPackage, msgName, err := pout.SplitIdentifier(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

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
			log.Fatal(err)
		}
	}
	defer r.Close()
	regy, err := registry.BuildProtoRegistry(config.Verbose, *protoPath)
	if err != nil {
		log.Fatal(fmt.Errorf("error BuildProtoRegistry(): %w", err))
	}

	msg, err := message.FromJSON(config.Verbose, regy, protoPackage, protoreflect.Name(msgName), r)
	if err != nil {
		wrapErr := fmt.Errorf("error converting JSON to proto: %w; check whether the message paths and types are correct", err)
		log.Fatal(wrapErr)
	}
	b, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)
	if *verbose {
		log.Printf("printing took %v", time.Since(start))
	}
}
