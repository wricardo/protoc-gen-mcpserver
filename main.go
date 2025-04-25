package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
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
		"toLower":      strings.ToLower,
		"fieldType":    getFieldType,
		"parseFunc":    getParseFunction,
		"mcpType":      getMcpType,
		"isRepeated":   isRepeated,
		"formatOutput": formatOutput,
		"defaultValue": getDefaultValue,
		"getBaseType":  getBaseType,
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

// isRepeated checks if a field is a repeated field (array/slice)
func isRepeated(field *protogen.Field) bool {
	return field.Desc.Cardinality() == protoreflect.Repeated && !field.Desc.IsMap()
}

// getFieldType returns the Go type of a protobuf field
func getFieldType(field *protogen.Field) string {
	if isRepeated(field) {
		baseType := getBaseType(field)
		return "[]" + baseType
	}

	return getBaseType(field)
}

// getBaseType returns the base Go type of a protobuf field
func getBaseType(field *protogen.Field) string {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "[]byte"
	case protoreflect.MessageKind:
		return "map[string]interface{}"
	default:
		return "string" // default to string for unsupported types
	}
}

// getMcpType returns the appropriate MCP type for a protobuf field
func getMcpType(field *protogen.Field) string {
	// For repeated fields, use WithArray
	if isRepeated(field) {
		return "WithArray"
	}

	// For non-repeated fields
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "WithBoolean"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return "WithNumber"
	case protoreflect.StringKind:
		return "WithString"
	case protoreflect.MessageKind:
		return "WithObject"
	case protoreflect.BytesKind:
		return "WithString" // Bytes represented as base64 string
	default:
		return "WithString" // Default to string for unsupported types
	}
}

// getParseFunction returns the appropriate parse function for a protobuf field
func getParseFunction(field *protogen.Field) string {
	// For repeated fields, we would need custom handling in the template
	if isRepeated(field) {
		// The template will need to handle array parsing specially
		return "ParseArray"
	}

	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "ParseBoolean"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "ParseInt32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "ParseUInt32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "ParseInt64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "ParseUInt64"
	case protoreflect.FloatKind:
		return "ParseFloat32"
	case protoreflect.DoubleKind:
		return "ParseFloat64"
	case protoreflect.StringKind:
		return "ParseString"
	case protoreflect.BytesKind:
		return "ParseString" // Parse as string, convert to []byte later
	case protoreflect.MessageKind:
		return "ParseStringMap"
	default:
		return "ParseString" // Default to string for unsupported types
	}
}

// getDefaultValue returns the default value for a field type
func getDefaultValue(field *protogen.Field) string {
	if isRepeated(field) {
		return "nil"
	}

	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.StringKind:
		return "\"\""
	case protoreflect.BytesKind:
		return "nil"
	case protoreflect.MessageKind:
		return "nil"
	default:
		// For all numeric types
		return "0"
	}
}

// formatOutput provides the correct formatting for output fields
func formatOutput(field *protogen.Field) string {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "fmt.Sprintf(\"%v\", res.{{ $field.GoName }})"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "fmt.Sprintf(\"%d\", res.{{ $field.GoName }})"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "fmt.Sprintf(\"%d\", res.{{ $field.GoName }})"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "fmt.Sprintf(\"%f\", res.{{ $field.GoName }})"
	case protoreflect.StringKind:
		return "res.{{ $field.GoName }}"
	case protoreflect.BytesKind:
		return "string(res.{{ $field.GoName }})"
	default:
		return "fmt.Sprintf(\"%v\", res.{{ $field.GoName }})"
	}
}

const mcpServerTemplate = `
// Code generated by protoc-gen-mcpserver. DO NOT EDIT.
package {{ .PackageName }}

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

{{- range $service := .Services }}
type {{ $service.GoName }}McpServer interface {
	{{- range $method := $service.Methods }}
	{{ $method.GoName }}(ctx context.Context, req *{{ $method.Input.GoIdent.GoName }}) (*{{ $method.Output.GoIdent.GoName }}, error)
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
			mcp.{{ mcpType $field }}("{{ $field.GoName }}", mcp.Description("Parameter {{ $field.GoName }}")),
			{{- end }}
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			req := &{{ $method.Input.GoIdent.GoName }}{}
			
			{{- range $field := $method.Input.Fields }}
			{{- if isRepeated $field }}
			// Handle repeated field {{ $field.GoName }}
			if arr, ok := request.Params.Arguments["{{ $field.GoName }}"]; ok && arr != nil {
				if arrValue, ok := arr.([]interface{}); ok {
					for _, v := range arrValue {
						{{- if eq (getBaseType $field) "string" }}
						if strVal, ok := v.(string); ok {
							req.{{ $field.GoName }} = append(req.{{ $field.GoName }}, strVal)
						}
						{{- else if eq (getBaseType $field) "bool" }}
						if boolVal, ok := v.(bool); ok {
							req.{{ $field.GoName }} = append(req.{{ $field.GoName }}, boolVal)
						}
						{{- else if or (eq (getBaseType $field) "int32") (eq (getBaseType $field) "int64") }}
						if numVal, ok := v.(float64); ok { // JSON numbers are float64
							req.{{ $field.GoName }} = append(req.{{ $field.GoName }}, {{ getBaseType $field }}(numVal))
						}
						{{- else if or (eq (getBaseType $field) "uint32") (eq (getBaseType $field) "uint64") }}
						if numVal, ok := v.(float64); ok && numVal >= 0 {
							req.{{ $field.GoName }} = append(req.{{ $field.GoName }}, {{ getBaseType $field }}(numVal))
						}
						{{- else if or (eq (getBaseType $field) "float32") (eq (getBaseType $field) "float64") }}
						if numVal, ok := v.(float64); ok {
							req.{{ $field.GoName }} = append(req.{{ $field.GoName }}, {{ getBaseType $field }}(numVal))
						}
						{{- else }}
						// Unsupported array element type
						{{- end }}
					}
				}
			}
			{{- else }}
			req.{{ $field.GoName }} = mcp.{{ parseFunc $field }}(request, "{{ $field.GoName }}", {{ defaultValue $field }})
			{{- end }}
			{{- end }}
			
			res, err := srv.{{ $method.GoName }}(ctx, req)
			if err != nil {
				return nil, err
			}
			
			result := &mcp.CallToolResult{
				Result:  mcp.Result{},
				Content: []mcp.Content{},
				IsError: false,
			}
			
			{{- range $field := $method.Output.Fields }}
			{{- if isRepeated $field }}
			// Format repeated field
			if len(res.{{ $field.GoName }}) > 0 {
				arrayStr := "["
				for i, v := range res.{{ $field.GoName }} {
					if i > 0 {
						arrayStr += ", "
					}
					{{- if eq (getBaseType $field) "string" }}
					arrayStr += fmt.Sprintf("%q", v)
					{{- else }}
					arrayStr += fmt.Sprintf("%v", v)
					{{- end }}
				}
				arrayStr += "]"
				result.Content = append(result.Content, mcp.NewTextContent("{{ $field.GoName }}: " + arrayStr))
			} else {
				result.Content = append(result.Content, mcp.NewTextContent("{{ $field.GoName }}: []"))
			}
			{{- else }}
			// Format non-repeated field
			{{- if eq (fieldType $field) "string" }}
			result.Content = append(result.Content, mcp.NewTextContent("{{ $field.GoName }}: " + res.{{ $field.GoName }}))
			{{- else if eq (fieldType $field) "[]byte" }}
			result.Content = append(result.Content, mcp.NewTextContent("{{ $field.GoName }}: " + string(res.{{ $field.GoName }})))
			{{- else }}
			result.Content = append(result.Content, mcp.NewTextContent("{{ $field.GoName }}: " + fmt.Sprintf("%v", res.{{ $field.GoName }})))
			{{- end }}
			{{- end }}
			{{- end }}
			
			return result, nil
		},
	)
	{{- end }}
}
{{- end }}

func ServeStdio(
name, 
version string, 
{{- range $service := .Services }}
srv{{ $service.GoName }} {{ $service.GoName }}McpServer,
{{- end }}
) error {
	s := server.NewMCPServer(name, version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)
{{- range $service := .Services }}
	Register{{ $service.GoName }}McpServer(s, srv{{ $service.GoName }})
{{- end }}
	return server.ServeStdio(s)
}
`
