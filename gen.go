package protoc_gen_markdown

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func ReadGenRequest(r io.Reader) *plugin.CodeGeneratorRequest {
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
