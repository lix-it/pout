package registry

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// TODO: don't use protoc
func createProtoRegistry(srcDir string, filename string) (*protoregistry.Files, error) {
	// Create descriptors using the protoc binary.
	// Imported dependencies are included so that the descriptors are self-contained.
	tmpFile := path.Base(filename) + "-tmp.pb"
	cmd := exec.Command("protoc",
		"--include_imports",
		"--descriptor_set_out="+tmpFile,
		"-I"+srcDir,
		path.Join(srcDir, filename))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	marshalledDescriptorSet, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		return nil, err
	}
	descriptorSet := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(marshalledDescriptorSet, &descriptorSet)
	if err != nil {
		return nil, err
	}

	files, err := protodesc.NewFiles(&descriptorSet)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func UnmarshalProto(protoPath, protoFile string, msgName protoreflect.Name, in []byte) []byte {
	msg, err := MakeDynamicMessage(protoPath, protoFile, msgName)
	if err != nil {
		panic(err)
	}

	err = proto.Unmarshal(in, msg)
	jsonBytes, err := protojson.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}

// MakeDynamicMessage creates an unhydrated proto message using the registry information
func MakeDynamicMessage(protoPath, protoFile string, msgName protoreflect.Name) (proto.Message, error) {
	registry, err := createProtoRegistry(protoPath, protoFile)
	if err != nil {
		return nil, fmt.Errorf("error createProtoRegistry(): %w", err)
	}

	desc, err := registry.FindFileByPath(protoFile)
	if err != nil {
		return nil, fmt.Errorf("error registry.FindFileByPath(): %w", err)
	}
	fd := desc.Messages()
	req := fd.ByName(msgName)
	msg := dynamicpb.NewMessage(req)
	return msg, nil
}

// FromJSON creates a dynamic Proto messages from JSON input.
func FromJSON(protoPath, protoFile string, msgName protoreflect.Name, jsonReader io.Reader) (proto.Message, error) {
	msg, err := MakeDynamicMessage(protoPath, protoFile, msgName)
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
