package main

import (
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"

	pgm "github.com/x-lambda/protoc-gen-markdown/generator"
)

func main() {
	req := pgm.ReadGenRequest(os.Stdin)

	//params, err := pgm.ParseParams(req.GetParameter())
	//if err != nil {
	//	panic(err)
	//}

	g := pgm.NewGenerator()
	resp := g.Generate(req)

	out, err := proto.Marshal(&resp)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
	os.Stdout.Write(out)
}
