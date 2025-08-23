# Doppelganger

Doppelganger is a Go library that enables seamless integration between Large Language Models (LLMs) and various data sources. It allows you to create AI-powered applications that can access and query external data sources like MongoDB and Google Cloud Storage through a tool-based approach.

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/doppelganger.svg)](https://pkg.go.dev/github.com/yourusername/doppelganger)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/doppelganger)](https://goreportcard.com/report/github.com/yourusername/doppelganger)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- 🔌 Connect LLMs to external data sources (MongoDB, Google Cloud Storage)
- 🛠️ Register custom tools with JSON schema validation
- 🤖 Supports multiple LLM providers (OpenAI, Anthropic)
- 🔄 Handles tool calling and response processing automatically
- 📝 Template-based query generation

## Installation

```bash
go get github.com/pyljain/doppelganger
```

## Quick Start

Here's a simple example of using Doppelganger to validate a SWIFT code using MongoDB:

```go
package main

import (
	"context"
	"doppelganger"
	"doppelganger/pkg/datasource"
	"doppelganger/pkg/tool"
	"fmt"
)

func main() {
	ctx := context.Background()
	app := doppelganger.New()

	// Connect to MongoDB
	mongoConnection := datasource.NewMongoDataSource()
	err := mongoConnection.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongoConnection.Close(ctx)

	// Register a tool to validate SWIFT codes
	tool := tool.DataSourceTool{
		Source:      mongoConnection,
		Name:        "validate_swift_code",
		Description: "Validates whether a swift code is valid",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"code": map[string]interface{}{
					"type": "string",
				},
			},
		},
		Query:      "{ \"swift_code\": \"{{ .code }}\" }",
		Database:   "my_database",
		Collection: "swift_codes",
		Method:     "findOne",
	}

	err = app.RegisterTool(tool)
	if err != nil {
		panic(err)
	}

	// Make a decision using the LLM
	systemInstruction := "You are a helpful assistant"
	prompt := "Can you validate if this swift code exists? Swift Code: UBSWCHZH80A"
	model := "gpt-4.1"

	result, err := app.MakeDecision(ctx, systemInstruction, prompt, model)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
```

## Tutorials

### 1. Connecting to MongoDB

```go
// Create a new MongoDB data source
mongoDS := datasource.NewMongoDataSource()

// Connect to MongoDB
err := mongoDS.Connect(ctx, "mongodb://localhost:27017")
if err != nil {
    panic(err)
}
defer mongoDS.Close(ctx)
```

### 2. Connecting to Google Cloud Storage

```go
// Create a new GCS data source
gcsDS := datasource.NewGCS(ctx)

// Connect to a specific bucket
err := gcsDS.Connect(ctx, "my-bucket-name")
if err != nil {
    panic(err)
}
defer gcsDS.Close(ctx)
```

### 3. Creating and Registering Tools

```go
// Create a new Doppelganger instance
app := doppelganger.New()

// Define a tool for listing files in GCS
listTool := tool.DataSourceTool{
    Source:      gcsDataSource,
    Name:        "list_files",
    Description: "List files in a storage bucket",
    Parameters:  map[string]interface{}{},
    Method:      "list",
}

// Register the tool
err = app.RegisterTool(listTool)
if err != nil {
    panic(err)
}

// Define a tool for getting file content
getTool := tool.DataSourceTool{
    Source:      gcsDataSource,
    Name:        "get_file_content",
    Description: "Get content of a file by name",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "name": map[string]interface{}{
                "type": "string",
            },
        },
    },
    Query:  "{{ .name }}",
    Method: "get",
}

// Register the tool
err = app.RegisterTool(getTool)
if err != nil {
    panic(err)
}
```

### 4. Making Decisions with LLMs

```go
// Define system instructions and user prompt
systemInstruction := "You are a helpful assistant that can access files in a storage bucket."
prompt := "Can you list all the policy documents available and then show me the content of the newest one?"
model := "gpt-4.1" // or "claude-3-opus-20240229"

// Make a decision
result, err := app.MakeDecision(ctx, systemInstruction, prompt, model)
if err != nil {
    panic(err)
}

// Print the result
fmt.Println(result)
```

## API Reference

### Doppelganger

#### `New() *Doppelganger`

Creates a new Doppelganger instance.

```go
app := doppelganger.New()
```

#### `RegisterTool(tool tool.DataSourceTool) error`

Registers a new tool with the Doppelganger instance.

```go
err := app.RegisterTool(myTool)
```

#### `MakeDecision(ctx context.Context, systemInstruction, userInstruction, model string) (string, error)`

Makes a decision using the specified LLM model, system instructions, and user prompt.

```go
result, err := app.MakeDecision(ctx, systemInstruction, prompt, "gpt-4.1")
```

### DataSourceTool

The `DataSourceTool` struct connects a data source to the LLM:

```go
type DataSourceTool struct {
    Source      datasource.DataSource
    Name        string
    Description string
    Parameters  map[string]interface{}
    Database    string
    Collection  string
    Method      string
    Query       string
}
```

- `Source`: The data source implementation
- `Name`: Name of the tool (used by the LLM)
- `Description`: Description of what the tool does
- `Parameters`: JSON schema for the tool parameters
- `Database`: Database name (for MongoDB)
- `Collection`: Collection name (for MongoDB)
- `Method`: Method to use (e.g., "findOne", "list", "get")
- `Query`: Template string for the query

### DataSource Interface

All data sources implement the `DataSource` interface:

```go
type DataSource interface {
    Connect(ctx context.Context, connectionString string) error
    Close(ctx context.Context) error
    Query(ctx context.Context, database, method, collection, query string) ([]string, error)
    Type() string
}
```

## Supported Data Sources

### MongoDB

```go
mongoDS := datasource.NewMongoDataSource()
err := mongoDS.Connect(ctx, "mongodb://localhost:27017")
```

### Google Cloud Storage

```go
gcsDS := datasource.NewGCS(ctx)
err := gcsDS.Connect(ctx, "bucket-name")
```

## Supported LLM Providers

Doppelganger supports the following LLM providers:

- OpenAI (models starting with "gpt-")
- Anthropic (models starting with "claude-")

## Advanced Usage

### Custom Query Templates

You can use Go's template syntax in your queries:

```go
tool := tool.DataSourceTool{
    // ...
    Query: "{ \"user_id\": \"{{ .userId }}\", \"status\": \"{{ .status }}\" }",
    // ...
}
```

### Error Handling

Always check for errors when registering tools and making decisions:

```go
err := app.RegisterTool(tool)
if err != nil {
    // Handle error
}

result, err := app.MakeDecision(ctx, systemInstruction, prompt, model)
if err != nil {
    // Handle error
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.