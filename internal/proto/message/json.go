package message

import (
	"fmt"
	"io"
	"io/ioutil"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// FromJSON creates a dynamic Proto messages from JSON input.
func FromJSON(verbose bool, registry *protoregistry.Files, protoPackage string, msgName protoreflect.Name, jsonReader io.Reader) (proto.Message, error) {
	msg, err := MakeDynamicMessage(verbose, registry, protoPackage, msgName)
	if err != nil {
		return nil, fmt.Errorf("error MakeDynamicMessage(): %w", err)
	}
	inB, err := ioutil.ReadAll(jsonReader)
	if err != nil {
		panic(err)
	}

	if err := protojson.Unmarshal(inB, msg); err != nil {
		return nil, fmt.Errorf("error protojson.Unmarshal(): %w", err)
	}
	return msg, nil
}
