package registry

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// BuildProtoRegistry walks through every file in the proto folder
// to build up a protoregistry.Files array with every single file in here.
// TODO: use a user-defined filter on the paths walked to improve performance
// TODO: don't use protoc
func BuildProtoRegistry(verbose bool, srcDir string) (*protoregistry.Files, error) {
	f := []string{}
	err := filepath.WalkDir(srcDir,
		func(p string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(p) == ".proto" {
				// remove root from dir
				pp := strings.SplitN(p, "/", 2)
				f = append(f, path.Join(srcDir, pp[1]))
			}

			return nil
		})
	if err != nil {
		panic(err)
	} // Create descriptors using the protoc binary.

	// Imported dependencies are included so that the descriptors are self-contained.
	tmpFile := path.Base("descriptor_set") + "-tmp.pb"
	args := []string{"--include_imports", "--descriptor_set_out=" + tmpFile, "-I" + srcDir}
	args = append(args, f...)
	cmd := exec.Command("protoc",
		args...)

	if verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	// how do I get a proto file for the file descriptor set?
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

func UnmarshalProto(verbose bool, protoPath, protoPackage string, msgName protoreflect.Name, in []byte) []byte {
	msg, err := MakeDynamicMessage(verbose, protoPath, protoPackage, msgName)
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

// MakeDynamicMessage creates an unhydrated proto message using the registry information
func MakeDynamicMessage(verbose bool, protoPath, protoPackage string, msgName protoreflect.Name) (proto.Message, error) {
	registry, err := BuildProtoRegistry(verbose, protoPath)
	if err != nil {
		return nil, fmt.Errorf("error BuildProtoRegistry(): %w", err)
	}
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

// FromJSON creates a dynamic Proto messages from JSON input.
func FromJSON(verbose bool, protoPath, protoFile string, msgName protoreflect.Name, jsonReader io.Reader) (proto.Message, error) {
	msg, err := MakeDynamicMessage(verbose, protoPath, protoFile, msgName)
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
