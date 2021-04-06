package pout

import (
	"strings"
)

// SplitIdentifier splits an input string into the package identifier
// and the message name.
// e.g. google.protobuf.Timestamp - package google.protobuf; message Timestamp
//
// returns ErrInvalidIdentifier if the input does not conform to this standard
func SplitIdentifier(input string) (string, string, error) {
	if err := ValidateIdentifier(input); err != nil {
		return "", "", err
	}
	pathIdentifier := strings.Split(input, ".")

	protoPackage := strings.Join(pathIdentifier[0:len(pathIdentifier)-1], ".")
	return protoPackage, pathIdentifier[len(pathIdentifier)-1], nil
}

func ValidateIdentifier(input string) error {
	if input[len(input)-1] == '.' {
		return ErrInvalidIdentifier{Detail: ErrTrailingDot}
	}
	return nil
}
