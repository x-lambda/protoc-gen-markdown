package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	req := readGenRequest(os.Stdin)

	// file_to_generate
	for _, v := range req.FileToGenerate {
		fmt.Println(v)
	}
}

func readGenRequest(r io.Reader) *plugin.CodeGeneratorRequest {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read err: %+v", err)
		panic(err)
	}

	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		fmt.Fprintf(os.Stderr, "proto unmarshal err: %+v", err)
		panic(err)
	}

	return req
}
