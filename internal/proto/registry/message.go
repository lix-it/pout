package registry

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// MakeDynamicMessage creates an unhydrated proto message using the registry information
func MakeDynamicMessage(verbose bool, registry *protoregistry.Files, protoPackage string, msgName protoreflect.Name) (proto.Message, error) {
	var req protoreflect.MessageDescriptor
	registry.RangeFilesByPackage(protoreflect.FullName(protoPackage), func(fd protoreflect.FileDescriptor) bool {
		req = fd.Messages().ByName(msgName)
		if verbose {
			fmt.Println("package file:", fd.Name())
		}
		return req == nil
	})
	if req == nil {
		panic("no message found!")
	}
	msg := dynamicpb.NewMessage(req)
	return msg, nil
}
