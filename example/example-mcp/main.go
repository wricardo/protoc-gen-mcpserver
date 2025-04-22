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
		Result: "Hello, " + req.Name,
	}, nil
}

func (s *GreetServer) Tool3(ctx context.Context, req Tool3Request) (*Tool3Response, error) {
	return &Tool3Response{
		HisFavoriteFood: "Hello, " + req.WallaceFavoriteFood,
	}, nil
}

func main() {
	err := ServeStdio("protoc-example-mcp", "0.1.0", &GreetServer{})
	if err != nil {
		log.Fatal(err)
	}

	// s := server.NewMCPServer("git-mcp", "0.1.0",
	// 	server.WithToolCapabilities(true),
	// 	server.WithLogging(),
	// )

	// RegisterMyToolsMcpServer(s, greeter)

	// server.ServeStdio(s)
}
