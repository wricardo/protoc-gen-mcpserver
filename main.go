package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.1.0"

func main() {
	flagVersion := flag.Bool("version", false, "Print the version and exit")
	flag.Parse()

	if *flagVersion {
		fmt.Println("protoc-gen-mcpserver version:", version)
		os.Exit(0)
	}

	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}
			generateFile(gen, file)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + ".mcpserver.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
	}

	tmpl, err := template.New("mcpserver").Funcs(funcMap).Parse(mcpServerTemplate)
	if err != nil {
		g.P("// Error parsing template: ", err)
		return
	}

	var data = struct {
		PackageName string
		Services    []*protogen.Service
		Methods     map[string][]*protogen.Method
	}{
		PackageName: string(file.GoPackageName),
		Services:    file.Services,
		Methods:     make(map[string][]*protogen.Method),
	}

	for _, service := range file.Services {
		data.Methods[service.GoName] = service.Methods
	}

	var builder strings.Builder
	if err := tmpl.Execute(&builder, data); err != nil {
		g.P("// Error executing template: ", err)
		return
	}

	g.P(builder.String())
}

const mcpServerTemplate = `package {{ .PackageName }}

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

{{- range $service := .Services }}
type {{ $service.GoName }}McpServer interface {
	{{- range $method := $service.Methods }}
	{{ $method.GoName }}(ctx context.Context, req {{ $method.Input.GoIdent.GoName }}) (*{{ $method.Output.GoIdent.GoName }}, error)
	{{- end }}
}

func Register{{ $service.GoName }}McpServer(s *server.MCPServer, srv {{ $service.GoName }}McpServer) {
	{{- range $method := $service.Methods }}
	s.AddTool(
		mcp.NewTool(
			"{{ $method.GoName }}",
			mcp.WithDescription("{{ $method.GoName }} description"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title: "{{ $method.GoName }}",
			}),
			{{- range $field := $method.Input.Fields }}
			mcp.WithString("{{ $field.GoName }}", mcp.Description("The name of the person to greet")),
			{{- end }}
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := srv.{{ $method.GoName }}(ctx, {{ $method.Input.GoIdent.GoName }}{
				{{- range $field := $method.Input.Fields }}
				{{ $field.GoName }}: mcp.ParseString(request, "{{ $field.GoName }}", ""),
				{{- end }}
			})
			if err != nil {
				return nil, err
			}
			result := &mcp.CallToolResult{
				Result:  mcp.Result{},
				Content: []mcp.Content{},
				IsError: false,
			}
			{{- range $field := $method.Output.Fields }}
			result.Content = append(result.Content, mcp.NewTextContent("{{ $field.GoName }}: "+res.{{ $field.GoName }}))
			{{- end }}
			return result, nil
		},
	)
	{{- end }}
}

func ServeStdio(name, version string, srv {{ $service.GoName }}McpServer) error {
	s := server.NewMCPServer(name, version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)
	Register{{ $service.GoName }}McpServer(s, srv)
	return server.ServeStdio(s)
}
{{- end }}
`
