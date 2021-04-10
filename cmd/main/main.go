package main

import (
	"fmt"
	"os"

	proto "github.com/x-lambda/protoc-gen-markdown"
)

func main() {
	req := proto.ReadGenRequest(os.Stdin)

	// file_to_generate
	for _, v := range req.FileToGenerate {
		fmt.Println(v)
	}
}
