package example

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MyToolsMcpServer interface {
	Tool1(ctx context.Context, req Tool1Request) (*Tool1Response, error)
	Tool2(ctx context.Context, req Tool2Request) (*Tool2Response, error)
	Tool3(ctx context.Context, req Tool3Request) (*Tool3Response, error)
}

func RegisterMyToolsMcpServer(s *server.MCPServer, srv MyToolsMcpServer) {
	s.AddTool(
		mcp.NewTool(
			"Tool1",
			mcp.WithDescription("Tool1 description"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title: "Tool1",
			}),
			mcp.WithString("Firstname", mcp.Description("The name of the person to greet")),
			mcp.WithString("Lastname", mcp.Description("The name of the person to greet")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := srv.Tool1(ctx, Tool1Request{
				Firstname: mcp.ParseString(request, "Firstname", ""),
				Lastname:  mcp.ParseString(request, "Lastname", ""),
			})
			if err != nil {
				return nil, err
			}
			// Get the field value from response - we need to handle different field names
			var responseText string
			responseText = res.Fullname
			return &mcp.CallToolResult{
				Result: mcp.Result{},
				Content: []mcp.Content{
					mcp.NewTextContent(responseText),
				},
				IsError: false,
			}, nil
		},
	)
	s.AddTool(
		mcp.NewTool(
			"Tool2",
			mcp.WithDescription("Tool2 description"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title: "Tool2",
			}),
			mcp.WithString("Name", mcp.Description("The name of the person to greet")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := srv.Tool2(ctx, Tool2Request{
				Name: mcp.ParseString(request, "Name", ""),
			})
			if err != nil {
				return nil, err
			}
			// Get the field value from response - we need to handle different field names
			var responseText string
			responseText = res.Result
			return &mcp.CallToolResult{
				Result: mcp.Result{},
				Content: []mcp.Content{
					mcp.NewTextContent(responseText),
				},
				IsError: false,
			}, nil
		},
	)
	s.AddTool(
		mcp.NewTool(
			"Tool3",
			mcp.WithDescription("Tool3 description"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title: "Tool3",
			}),
			mcp.WithString("WallaceFavoriteFood", mcp.Description("The name of the person to greet")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := srv.Tool3(ctx, Tool3Request{
				WallaceFavoriteFood: mcp.ParseString(request, "WallaceFavoriteFood", ""),
			})
			if err != nil {
				return nil, err
			}
			// Get the field value from response - we need to handle different field names
			var responseText string
			responseText = res.HisFavoriteFood
			return &mcp.CallToolResult{
				Result: mcp.Result{},
				Content: []mcp.Content{
					mcp.NewTextContent(responseText),
				},
				IsError: false,
			}, nil
		},
	)
}

func ServeStdio(name, version string, srv MyToolsMcpServer) error {
	s := server.NewMCPServer(name, version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)
	RegisterMyToolsMcpServer(s, srv)
	return server.ServeStdio(s)
}
