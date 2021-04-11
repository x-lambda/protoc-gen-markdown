package generator

import (
	"fmt"
	"testing"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func TestReadGenRequest(t *testing.T) {
	// TODO 手动构造 req 对象
	req := plugin.CodeGeneratorRequest{}
	fmt.Println(req)
}
