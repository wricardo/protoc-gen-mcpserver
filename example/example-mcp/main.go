package main

import (
	"context"
	"log"

	. "github.com/wricardo/protoc-gen-mcpserver/example"
)

type GreetServer struct {
}

// CalculateSum implements example.ExampleServiceMcpServer.
func (s *GreetServer) CalculateSum(ctx context.Context, req *CalculateSumRequest) (*CalculateSumResponse, error) {
	return &CalculateSumResponse{
		Sum:     req.Number1 + req.Number2,
		Product: float64(req.Number1) * float64(req.Number2) * float64(req.Factor),
	}, nil
}

func (s *GreetServer) CheckStatus(ctx context.Context, req *CheckStatusRequest) (*CheckStatusResponse, error) {
	return &CheckStatusResponse{
		Success: true,
		Message: "Status is active",
	}, nil
}

func (s *GreetServer) ComplexOperation(ctx context.Context, req *ComplexOperationRequest) (*ComplexOperationResponse, error) {
	return &ComplexOperationResponse{
		Success:     true,
		OperationId: "12345",
		StatusCode:  200,
		Results:     []string{"Result1", "Result2"},
		Average:     10.5,
	}, nil
}

func (s *GreetServer) GreetPerson(ctx context.Context, req *GreetPersonRequest) (*GreetPersonResponse, error) {
	return &GreetPersonResponse{
		Greeting: "Hello, " + req.FirstName + " " + req.LastName,
	}, nil
}

// ProcessNames implements example.ExampleServiceMcpServer.
func (s *GreetServer) ProcessNames(ctx context.Context, req *ProcessNamesRequest) (*ProcessNamesResponse, error) {
	return &ProcessNamesResponse{
		ProcessedNames:  []string{"Processed1", "Processed2"},
		ProcessedCounts: []int32{1, 2},
		Summary:         "Processed 2 names",
	}, nil
}

func (s *GreetServer) Tool1(ctx context.Context, req *Tool1Request) (*Tool1Response, error) {
	return &Tool1Response{
		Fullname: "Hello, " + req.Firstname + " " + req.Lastname,
	}, nil
}

func (s *GreetServer) Tool2(ctx context.Context, req *Tool2Request) (*Tool2Response, error) {
	return &Tool2Response{
		Result: "Hello, " + req.Name,
	}, nil
}

func (s *GreetServer) Tool3(ctx context.Context, req *Tool3Request) (*Tool3Response, error) {
	return &Tool3Response{
		HisFavoriteFood: "Hello, " + req.WallaceFavoriteFood,
	}, nil
}

func main() {
	greeter := &GreetServer{}

	err := ServeStdio("protoc-example-mcp", "0.1.0", greeter, greeter)
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
