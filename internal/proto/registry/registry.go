package registry

import (
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
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
