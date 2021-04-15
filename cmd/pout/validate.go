package main

import (
	"errors"
	"flag"
	"log"
)

func validateArgs() error {
	if len(flag.Args()) < 2 {
		return errors.New("you need to specify a file and a message type")
	}
	return nil
}

func validateFlags() error {
	if *protoPath == "" {
		log.Fatal("you must enter a valid proto base path")
	}
	return nil
}
