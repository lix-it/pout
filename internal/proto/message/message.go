package message

import (
	"github.com/lix-it/pout/internal/proto/registry"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func UnmarshalProto(verbose bool, rgy *protoregistry.Files, protoPackage string, msgName protoreflect.Name, in []byte) []byte {
	msg, err := registry.MakeDynamicMessage(verbose, rgy, protoPackage, msgName)
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(in, msg)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := protojson.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}
