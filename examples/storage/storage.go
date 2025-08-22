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

	gcsDataSource := datasource.NewGCS(ctx)

	err := gcsDataSource.Connect(ctx, "sparkbox")
	if err != nil {
		panic(err)
	}
	defer gcsDataSource.Close(ctx)

	listTool := tool.DataSourceTool{
		Source:      gcsDataSource,
		Name:        "list_policies",
		Description: "List banking policies",
		Parameters:  map[string]interface{}{},
		Method:      "list",
	}

	err = app.RegisterTool(listTool)
	if err != nil {
		panic(err)
	}

	getTool := tool.DataSourceTool{
		Source:      gcsDataSource,
		Name:        "get_policy_document",
		Description: "Get banking policy document by name of file",
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

	err = app.RegisterTool(getTool)
	if err != nil {
		panic(err)
	}

	systemInstruction := "You are a helpful assistant"
	prompt := "Can you tell me about the special rules that apply to large deposits at Caja Test España?"
	model := "gpt-4.1"

	result, err := app.MakeDecision(ctx, systemInstruction, prompt, model)
	if err != nil {
		panic(err)
	}

	println(result)
}
