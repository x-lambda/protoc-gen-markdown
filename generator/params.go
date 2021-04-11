package generator

import (
	"fmt"
	"strings"
)

type Options struct {
	Prefix   string
	Package  string
	Filename string
	// TODO callback
}

// ParseParams 解析参数
// 参数格式 protoc --markdown_out=prefix=api,package=demo.v0,filename=x.md
// prefix: 表示在
// package: 如果参数中未指定，则使用proto文件中定义
// filename: 生成的文件名
func ParseParams(params string) (opt Options, err error) {
	pamars := make(map[string]string)
	for _, v := range strings.Split(params, ",") {
		if v == "" {
			continue
		}

		i := strings.Index(v, "=")
		if i < 0 {
			err = fmt.Errorf("invalid parameter %s: expected format of parameter to be key=value", v)
			return
		}

		key := v[0:i]
		value := v[i+1:]
		if value == "" {
			err = fmt.Errorf("invalid parameter: value can't be empty")
			return
		}

		pamars[key] = value
	}

	for key, value := range pamars {
		switch key {
		case "prefix":
			opt.Prefix = value
		case "package":
			opt.Package = value
		case "filename":
			opt.Filename = value
		default:
		}
	}
	return
}
