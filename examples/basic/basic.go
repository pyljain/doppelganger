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

	mongoConnection := datasource.NewMongoDataSource()

	err := mongoConnection.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongoConnection.Close(ctx)

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

	systemInstruction := "You are a helpful assistant"
	prompt := "Can you validate if this swift code exists? Swift Code: UBSWCHZH80A"
	model := "gpt-4.1"

	result, err := app.MakeDecision(ctx, systemInstruction, prompt, model)
	if err != nil {
		panic(err)
	}

	println(result)
}
