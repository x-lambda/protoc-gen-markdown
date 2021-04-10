package protoc_gen_markdown

import (
	"fmt"
	"strings"
)

type Params struct {
	PathPrefix string
}

// parseParams 解析参数
// 参数格式 protoc --markdown_out=path_prefix=/tmp,package_name=api:.
func parseParams(params string) (p Params, err error) {
	temp := make(map[string]string)
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

		temp[key] = value
	}

	for key, value := range temp {
		if key == "path_prefix" {
			p.PathPrefix = value
		}
	}
	return
}
