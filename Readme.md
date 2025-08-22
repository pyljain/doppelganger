# Doppelganger

A Go library for creating AI agents that can interact with multiple data sources through configurable tools. Doppelganger enables LLMs to query databases, cloud storage, and other data sources to provide intelligent responses.

## Features

- **Multi-Provider LLM Support**: Compatible with OpenAI GPT models and Anthropic Claude models
- **Pluggable Data Sources**: Support for MongoDB, Google Cloud Storage, and extensible architecture for custom data sources
- **Tool Registration System**: Define custom tools with JSON Schema validation for parameters
- **Template-Based Queries**: Use Go templates to construct dynamic queries from user parameters
- **Automatic Tool Calling**: Handle LLM tool calls and responses automatically in a conversation loop

## Architecture

- **Core Engine** (`doppelganger.go`): Main orchestration logic for LLM interactions and tool calling
- **Data Sources** (`pkg/datasource/`): Abstraction layer for different data backends
- **Tools** (`pkg/tool/`): Wrapper system to expose data sources as LLM-callable functions
- **LLM Providers** (`pkg/llm/`): Support for multiple LLM providers

## Quick Start

```go
package main

import (
    "context"
    "doppelganger"
    "doppelganger/pkg/datasource"
    "doppelganger/pkg/tool"
)

func main() {
    ctx := context.Background()
    app := doppelganger.New()
    
    // Setup MongoDB connection
    mongoConnection := datasource.NewMongoDataSource()
    err := mongoConnection.Connect(ctx, "mongodb://localhost:27017")
    if err != nil {
        panic(err)
    }
    defer mongoConnection.Close(ctx)
    
    // Register a tool
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
    
    // Make a decision
    result, err := app.MakeDecision(ctx, 
        "You are a helpful assistant",
        "Can you validate if this swift code exists? Swift Code: UBSWCHZH80A",
        "gpt-4.1")
    if err != nil {
        panic(err)
    }
    
    println(result)
}
```

## Supported Data Sources

### MongoDB
- Connect to MongoDB instances
- Execute findOne, find, and other MongoDB operations
- Template-based query construction

### Google Cloud Storage
- List files in GCS buckets
- Retrieve file contents
- Support for various file formats

### Custom Data Sources
Implement the `DataSource` interface to add support for additional backends:

```go
type DataSource interface {
    Connect(ctx context.Context, connectionString string) error
    Close(ctx context.Context) error
    Query(ctx context.Context, database, method, collection, query string) ([]string, error)
    Type() string
}
```

## Supported LLM Providers

- **OpenAI**: GPT models (gpt-3.5-turbo, gpt-4, etc.)
- **Anthropic**: Claude models (claude-3-sonnet, claude-3-opus, etc.)

## Requirements

- Go 1.24.0 or later
- Valid API keys for chosen LLM provider (set via environment variables)
- Access to configured data sources

## Examples

See the `examples/` directory for complete working examples:
- `examples/basic/`: MongoDB integration with SWIFT code validation
- `examples/storage/`: Google Cloud Storage integration for document retrieval

## Dependencies

- [langchaingo](https://github.com/tmc/langchaingo): LLM provider abstraction
- [mongo-driver](https://go.mongodb.org/mongo-driver): MongoDB connectivity
- [gojsonschema](https://github.com/xeipuuv/gojsonschema): JSON Schema validation
- [json-iterator](https://github.com/json-iterator/go): High-performance JSON processing