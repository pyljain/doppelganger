package main

import (
	"context"
	"doppelganger"
	"doppelganger/pkg/datasource"
	"doppelganger/pkg/tool"
)

func main() {
	app := doppelganger.New()

	mongoConnection := datasource.NewMongoDataSource()

	err := mongoConnection.Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongoConnection.Close()

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
		Query: "{ code: \"$code\" }",
	}

	err = app.RegisterTool(tool)
	if err != nil {
		panic(err)
	}

	systemInstruction: = "You are a helpful assistant"
	prompt := "Can you validate that this transaction is valid?"
	model := "gpt-3.5-turbo"

	ctx := context.Background()
	result, err := app.MakeDecision(ctx, systemInstruction, prompt, model)
	if err != nil {
		panic(err)
	}

	println(result)
}
