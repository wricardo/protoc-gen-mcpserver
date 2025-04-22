package main

import (
	"context"
	"log"

	. "github.com/wricardo/protoc-gen-mcpserver/example"
)

type GreetServer struct {
}

func (s *GreetServer) Tool1(ctx context.Context, req Tool1Request) (*Tool1Response, error) {
	return &Tool1Response{
		Fullname: "Hello, " + req.Firstname + " " + req.Lastname,
	}, nil
}

func (s *GreetServer) Tool2(ctx context.Context, req Tool2Request) (*Tool2Response, error) {
	return &Tool2Response{
		Result: "Hello 22",
	}, nil
}

func (s *GreetServer) Tool3(ctx context.Context, req Tool3Request) (*Tool3Response, error) {
	return &Tool3Response{
		HisFavoriteFood: "Hello 33",
	}, nil
}

func main() {
	err := ServeStdio("example2-mcp", "0.1.0", &GreetServer{})
	if err != nil {
		log.Fatal(err)
	}
}
