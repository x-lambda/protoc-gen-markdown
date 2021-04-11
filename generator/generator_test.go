package generator

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/stretchr/testify/assert"
)

func TestReadGenRequest(t *testing.T) {
	// 命令行参数
	params := "prefix=/api,package=a"

	// pb文件
	name := "test.proto"

	// pb文件中定义的 package
	packageName := "demo.v0"

	// 分别对应pb文件中定义的message
	messageReqFieldName := "name"
	messageReqFieldNumber := int32(1)
	messageReqFieldLabel := descriptor.FieldDescriptorProto_LABEL_OPTIONAL
	messageReqFieldType := descriptor.FieldDescriptorProto_TYPE_STRING
	messageReqJsonName := "name"
	// message 中的 field 定义
	messageReqField := &descriptor.FieldDescriptorProto{
		Name:     &messageReqFieldName,
		Number:   &messageReqFieldNumber,
		Label:    &messageReqFieldLabel,
		Type:     &messageReqFieldType,
		JsonName: &messageReqJsonName,
	}
	// message 定义
	messageReqName := "SearchReq"
	messageReq := &descriptor.DescriptorProto{
		Name:  &messageReqName,
		Field: []*descriptor.FieldDescriptorProto{messageReqField},
	}

	// TODO 定义一堆 field
	// message 定义
	messageDataName := "ReplyData"
	messageData := &descriptor.DescriptorProto{
		Name: &messageDataName,
	}

	// TODO 定义一堆 field
	// message 定义
	messageRespName := "SearchResp"
	messageResp := &descriptor.DescriptorProto{
		Name: &messageRespName,
	}

	// 对应pb文件中对应的service
	serviceName := ""
	serviceMethodName := "Search"       // 对应pb中的rpc name
	inputType := ".demo.v0.SearchReq"   // 对应message定义，需要处理一下
	outputType := ".demo.v0.SearchResp" // 对应message定义，需要处理一下
	serviceMethod := &descriptor.MethodDescriptorProto{
		Name:       &serviceMethodName,
		InputType:  &inputType,
		OutputType: &outputType,
	}
	service := &descriptor.ServiceDescriptorProto{
		Name:   &serviceName,
		Method: []*descriptor.MethodDescriptorProto{serviceMethod},
	}

	syntax := "proto3"

	desc := &descriptor.FileDescriptorProto{
		Name:    &name,
		Package: &packageName,
		MessageType: []*descriptor.DescriptorProto{
			messageReq, messageResp, messageData,
		},
		Service: []*descriptor.ServiceDescriptorProto{service},
		Syntax:  &syntax,

		// TODO SourceCodeInfo ?
		// SourceCodeInfo: nil,
	}

	// protoc 解析 pb 文件之后，跟生成一个这样的对象
	req := plugin.CodeGeneratorRequest{
		FileToGenerate: []string{"test.proto"},
		Parameter:      &params,
		ProtoFile:      []*descriptor.FileDescriptorProto{desc},
	}

	// 经过 proto 序列化，生成字节对象
	data, err := proto.Marshal(&req)
	assert.Nil(t, err)

	// 通过io发送给插件
	// 插件收到数据之后，反序列化之后就可以得到上面的CodeGeneratorRequest对象
	// 然后可以生成指定的模板代码/文档等等
	buf := bytes.NewReader(data)
	expect, err := ReadGenRequest(buf)
	assert.Nil(t, err)

	assert.Equal(t, *expect.Parameter, *req.Parameter)
	assert.Equal(t, len(expect.FileToGenerate), len(req.FileToGenerate))
	assert.Equal(t, expect.FileToGenerate[0], req.FileToGenerate[0])

	// TODO test generate markdown
}
