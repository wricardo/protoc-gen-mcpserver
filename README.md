# protoc-gen-mcpserver

A Protobuf compiler plugin that generates Go code for serving MCP (Model-Code-Protocol) tools from Protobuf service definitions.

## Overview

This tool auto-generates the necessary Go code to create MCP server with tools tools, which can be consumed by AI assistants like Claude, all starting from a proto definition. It enables you to:

- Define your tools using Protocol Buffers
- Automatically generate the necessary MCP server code
- Implement only the business logic for your service methods

## Installation

Clone this repository and install the plugin:

```bash
go install github.com/wricardo/protoc-gen-mcpserver@latest
```

Or build from source:

```bash
git clone https://github.com/wricardo/protoc-gen-mcpserver.git
cd protoc-gen-mcpserver
go install
```

Make sure `protoc-gen-mcpserver` is in your PATH.

## Usage

### 1. Define your service in a .proto file

Create a Protocol Buffer file that defines your service and message types:

```protobuf
syntax = "proto3";

package yourpackage;

option go_package = "yourmodule/yourpackage;yourpackage";

service YourService {
  rpc YourMethod(YourRequest) returns (YourResponse);
}

message YourRequest {
  string parameter1 = 1;
  string parameter2 = 2;
}

message YourResponse {
  string result = 1;
}
```

### 2. Configure buf.gen.yaml

Use [buf](https://buf.build/) to generate code from your protobuf definitions. Create a `buf.gen.yaml` file:

```yaml
version: v2
managed:
  enabled: true
plugins:
  - name: go
    out: ./
    opt: paths=source_relative
  - name: mcpserver
    out: ./
    opt: paths=source_relative
```

### 3. Generate the code

Run the code generation:

```bash
buf generate
```

This will create:
- Standard Go Protobuf code
- An additional `.mcpserver.go` file with MCP server integration

### 4. Implement your service

Create a struct that implements the interface defined in the generated code:

```go
package main

import (
	"context"
	"log"
	
	. "yourmodule/yourpackage"
)

type YourServiceImpl struct {}

func (s *YourServiceImpl) YourMethod(ctx context.Context, req YourRequest) (*YourResponse, error) {
	// Implement your business logic here
	return &YourResponse{
		Result: "Processed: " + req.Parameter1 + " " + req.Parameter2,
	}, nil
}

func main() {
	// Use the generated ServeStdio function to start the MCP server
	err := ServeStdio("your-mcp-tool", "1.0.0", &YourServiceImpl{})
	if err != nil {
		log.Fatal(err)
	}
}
```

### 5. Build and run your MCP server

Build your implementation:

```bash
go build -o your-mcp-tool
```

Your tool can now be executed by MCP-compatible clients.

## How It Works

The plugin generates:

1. Interface definitions for your service
2. Registration functions to add your methods as MCP tools
3. Helper functions for serving the MCP protocol over stdio

Each method in your gRPC service becomes an MCP tool, with request fields automatically mapped to tool parameters.

## Example

Check out the included example to see a complete working implementation:

```bash
cd example
buf generate
cd example-mcp
go build -o example-mcp
./example-mcp
```

## License

[MIT License](LICENSE)
