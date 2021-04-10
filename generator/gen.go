package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
)

// ReadGenRequest
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

type field struct {
	Name    string
	Type    string
	KeyType string
	Note    string
	Doc     string
	Label   descriptor.FieldDescriptorProto_Label
}

func (f field) isRepeated() bool {
	return f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

type message struct {
	Name   string
	Fields []field
	Label  descriptor.FieldDescriptorProto_Label
	Doc    string
}

type api struct {
	Method  string
	Path    string
	Doc     string
	Request *message
	Reply   *message
	Input   string
	Output  string
}

type Generator struct {
	output *bytes.Buffer

	// Map of all proto messages
	messages map[string]*message

	enums map[string]*protokit.EnumDescriptor

	// List of all APIs
	apis []*api

	// List of all service comments
	comments *protokit.Comment

	// Service name
	name string
}

func NewGenerator() Generator {
	return Generator{
		messages: map[string]*message{},
		enums:    map[string]*protokit.EnumDescriptor{},
		apis:     []*api{},
		output:   bytes.NewBuffer(nil),
	}
}

func (g *Generator) Generate(in *plugin.CodeGeneratorRequest) (resp plugin.CodeGeneratorResponse) {
	g.scanAllMessage(in, &resp)
	g.GenerateMarkdown(in, &resp)
	return
}

func (g *Generator) P(args ...string) {
	for _, v := range args {
		g.output.WriteString(v)
	}

	g.output.WriteByte('\n')
}

func (g *Generator) scanAllMessage(req *plugin.CodeGeneratorRequest, resp *plugin.CodeGeneratorResponse) {
	descriptors := protokit.ParseCodeGenRequest(req)
	for _, v := range descriptors {
		g.scanMessages(v)
	}
	return
}

func (g *Generator) scanMessages(d *protokit.FileDescriptor) {
	// TODO 支持enum
	for _, md := range d.GetMessages() {
		g.scanMessage(md)
	}
}

func (g *Generator) scanEnum(md *protokit.EnumDescriptor) {
	g.enums["."+md.GetFullName()] = md
}

func (g *Generator) scanMessage(md *protokit.Descriptor) {
	for _, smd := range md.GetMessages() {
		g.scanMessage(smd)
	}

	for _, ed := range md.GetEnums() {
		g.scanEnum(ed)
	}

	{
		fields := make([]field, len(md.GetMessageFields()))
		maps := make(map[string]*descriptor.DescriptorProto)

		for _, v := range md.NestedType {
			if v.Options.GetMapEntry() {
				pkg := md.GetPackage()
				name := fmt.Sprintf(".%s.%s.%s", pkg, md.GetName(), v.GetName())
				maps[name] = v
			}
		}

		for i, fd := range md.GetMessageFields() {
			typeName := fd.GetTypeName()
			if typeName == "" {
				typeName = fd.GetType().String()
			}

			f := field{
				Name:  fd.GetName(),
				Type:  typeName,
				Doc:   fd.GetComments().GetLeading(),
				Note:  fd.GetComments().GetTrailing(),
				Label: fd.GetLabel(),
			}

			if e, ok := g.enums[fd.GetTypeName()]; ok {
				f.Type = "TYPE_ENUM"
				parts := []string{}

				for _, v := range e.GetValues() {
					line := fmt.Sprintf("%s(=%d) %s", v.GetName(), v.GetNumber(), v.GetComments().GetTrailing())
					parts = append(parts, line)
				}

				f.Doc = strings.Join(parts, "\n")
			}

			if m, ok := maps[f.Type]; ok {
				for _, ff := range m.GetField() {
					switch ff.GetName() {
					case "key":
						f.KeyType = ff.GetType().String()
					case "value":
						typeName := ff.GetTypeName()
						if typeName == "" {
							typeName = ff.GetType().String()
						}
						f.Type = typeName
					}
				}
			}
			fields[i] = f
		}

		g.messages[md.GetFullName()] = &message{
			Name:   md.GetName(),
			Doc:    md.GetComments().GetTrailing(),
			Fields: fields,
		}
	}
}

func (g *Generator) scanService(d *protokit.ServiceDescriptor) {
	g.comments = d.Comments

	for _, md := range d.GetMethods() {
		api := api{}

		// TODO callback generate
		api.Method = "POST"
		api.Path = "" + "/" + d.GetFullName() + "/" + md.GetName()
		doc := md.GetComments().GetLeading()

		// 支持文档换行
		api.Doc = strings.Replace(doc, "\n", "\n\n", -1)

		inputType := md.GetInputType()[1:] // trim leading dot
		api.Request = g.messages[inputType]

		outputType := md.GetOutputType()[1:] // trim leading dot
		api.Reply = g.messages[outputType]

		g.apis = append(g.apis, &api)
	}
}

func (g *Generator) GenerateMarkdown(req *plugin.CodeGeneratorRequest, resp *plugin.CodeGeneratorResponse) {
	descriptors := protokit.ParseCodeGenRequest(req)

	// 一个 service 对应一个 markdown 文件
	for _, d := range descriptors {
		for _, sd := range d.GetServices() {
			g.scanService(sd)

			g.name = *sd.Name
			for _, api := range g.apis {
				api.Input = g.generateJsDocForMessage(api.Request)
				api.Output = g.generateJsDocForMessage(api.Reply)
			}

			g.generateDoc()

			// 输出的文件信息
			name := strings.Replace(d.GetName(), ".proto", ".md", 1)
			resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(name),
				Content: proto.String(g.output.String()),
			})
		}
	}
}

func (g *Generator) generateJsDocForField(f field) string {
	var js string
	var v, vt string
	disableDoc := false

	if f.Doc != "" {
		for _, line := range strings.Split(f.Doc, "\n") {
			js += "//" + line + "\n"
		}
	}

	if f.Type == "TYPE_STRING" {
		vt = "string"
		if f.isRepeated() {
			v = `["", ""]`
		} else if f.KeyType != "" {
			v = fmt.Sprintf(`{"%s": ""}`, getTypeValue(f.KeyType))
			vt = fmt.Sprintf(`map<%s, string>`, getType(f.KeyType))
		} else {
			v = `""`
		}
	} else if f.Type == "TYPE_DOUBLE" || f.Type == "TYPE_FLOAT" {
		vt = "float"
		if f.isRepeated() {
			v = "[0.0, 0.0]"
		} else if f.KeyType != "" {
			v = fmt.Sprintf(`{"%s":0.0}`, getTypeValue(f.KeyType))
			vt = fmt.Sprintf("map<%s,float>", getType(f.KeyType))
		} else {
			v = "0.0"
		}
	} else if f.Type == "TYPE_BOOL" {
		vt = "bool"
		if f.isRepeated() {
			v = "[false, false]"
		} else if f.KeyType != "" {
			v = fmt.Sprintf(`{"%s":false}`, getTypeValue(f.KeyType))
			vt = fmt.Sprintf("map<%s,bool>", getType(f.KeyType))
		} else {
			v = "false"
		}
	} else if f.Type == "TYPE_INT64" || f.Type == "TYPE_UINT64" {
		vt = "string(int64)"
		if f.isRepeated() {
			v = `["0", "0"]`
		} else if f.KeyType != "" {
			v = fmt.Sprintf(`{"%s":"0"}`, getTypeValue(f.KeyType))
			vt = fmt.Sprintf("map<%s,string(int64)>", getType(f.KeyType))
		} else {
			v = `"0"`
		}
	} else if f.Type == "TYPE_INT32" || f.Type == "TYPE_UINT32" {
		vt = "int"
		if f.isRepeated() {
			v = "[0, 0]"
		} else if f.KeyType != "" {
			v = fmt.Sprintf(`{"%s":0}`, getTypeValue(f.KeyType))
			vt = fmt.Sprintf("map<%s,int>", getType(f.KeyType))
		} else {
			v = "0"
		}
	} else if f.Type == "TYPE_ENUM" {
		vt = "string(enum)"
		if f.isRepeated() {
			v = `["", ""]`
		} else {
			v = `""`
		}
	} else if f.Type[0] == '.' {
		m := g.messages[f.Type[1:]]
		v = g.generateJsDocForMessage(m)
		if f.isRepeated() {
			doc := fmt.Sprintf("// type:<list<%s>>", m.Name)
			if f.Note != "" {
				doc = " " + f.Note
			}
			v = "[" + doc + "\n" + v + "]"
		} else if f.KeyType != "" {
			doc := fmt.Sprintf("// type:<map<%s,%s>>", getType(f.KeyType), m.Name)
			if f.Note != "" {
				doc = " " + f.Note
			}
			v = fmt.Sprintf("{%s\n\"%s\":%s}", doc, getTypeValue(f.KeyType), v)
		}
		disableDoc = true
	} else {
		v = "UNKNOWN"
	}

	if disableDoc {
		js += fmt.Sprintf("%s: %s", f.Name, v)
	} else {
		js += fmt.Sprintf("%s: %s, // type: <%s>", f.Name, v, vt)
		if f.Note != "" {
			js = js + ", " + f.Note
		}
	}

	js = strings.Trim(js, " ")
	js += "\n"
	return js
}

func (g *Generator) generateJsDocForMessage(m *message) string {
	var js string
	js += "{\n"

	for _, f := range m.Fields {
		js += g.generateJsDocForField(f)
	}

	js += "\n"

	return js
}

func (g *Generator) generateDoc() {
	options := jsbeautifier.DefaultOptions()
	g.P("# ", g.name)
	g.P()

	comments := strings.Split(g.comments.Leading, "\n")
	for _, v := range comments {
		g.P(v, " ")
	}
	g.P()

	for _, api := range g.apis {
		anchor := strings.Replace(api.Path, "/", "", -1)
		anchor = strings.Replace(anchor, ".", "", -1)
		anchor = strings.ToLower(anchor)
		g.P(fmt.Sprintf("- [%s](#%s)", api.Path, anchor))
	}
	g.P()

	for _, api := range g.apis {
		g.P("## ", api.Path)
		g.P()
		g.P(api.Doc)
		g.P()
		g.P("### Method")
		g.P()
		g.P(api.Method)
		g.P()
		g.P("### Request")
		g.P("```javascript")
		code, _ := jsbeautifier.Beautify(&api.Input, options)
		g.P(code)
		g.P("```")
		g.P()
		g.P("### Reply")
		g.P("```javascript")
		code, _ = jsbeautifier.Beautify(&api.Output, options)
		g.P(code)
		g.P("```")
	}
}

func getType(t string) string {
	switch t {
	case "TYPE_STRING":
		return "string"
	case "TYPE_DOUBLE", "TYPE_FLOAT":
		return "float"
	case "TYPE_BOOL":
		return "bool"
	case "TYPE_INT64", "TYPE_UINT64", "TYPE_INT32", "TYPE_UINT32":
		return "int"
	default:
		return t
	}
}

func getTypeValue(t string) string {
	switch t {
	case "TYPE_STRING":
		return ""
	case "TYPE_DOUBLE", "TYPE_FLOAT":
		return "0.0"
	case "TYPE_BOOL":
		return "false"
	case "TYPE_INT64", "TYPE_UINT64", "TYPE_INT32", "TYPE_UINT32":
		return "0"
	default:
		return ""
	}
}
