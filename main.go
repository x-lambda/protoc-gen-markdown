package main

import (
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	pgm "github.com/x-lambda/protoc-gen-markdown/generator"
)

func main() {
	req, err := pgm.ReadGenRequest(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	opt, err := pgm.ParseParams(req.GetParameter())
	if err != nil {
		log.Fatal(err)
	}

	g := pgm.NewGenerator(opt)
	resp := g.Generate(&req)
	out, err := proto.Marshal(&resp)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(out)
	return
}
