package pout

import (
	"strings"
)

// SplitIdentifier splits an input string into the package identifier
// and the message name.
// e.g. google.protobuf.Timestamp - package google.protobuf; message Timestamp
func SplitIdentifier(input string) (string, string) {
	pathIdentifier := strings.Split(input, ".")

	protoPackage := strings.Join(pathIdentifier[0:len(pathIdentifier[0])-2], ".")
	return protoPackage, pathIdentifier[len(pathIdentifier)-1]
}
