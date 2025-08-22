package doppelganger

import (
	"context"
	"doppelganger/pkg/llm"
	"doppelganger/pkg/tool"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/tmc/langchaingo/llms"
	"github.com/xeipuuv/gojsonschema"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ProviderGeneratorFunc func(model string) (llms.Model, error)

type Doppelganger struct {
	Tools                 []tool.DataSourceTool
	providerGeneratorFunc ProviderGeneratorFunc
	toolsMap              map[string]tool.DataSourceTool
}

func New() *Doppelganger {
	return &Doppelganger{
		providerGeneratorFunc: llm.GetProvider,
		toolsMap:              make(map[string]tool.DataSourceTool),
	}
}

func (d *Doppelganger) RegisterTool(tool tool.DataSourceTool) error {
	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true
	sl.Draft = gojsonschema.Draft7
	sl.AutoDetect = false

	err := sl.AddSchemas(gojsonschema.NewGoLoader(tool.Parameters))
	if err != nil {
		return err
	}

	// Save tool definition to the base struct
	d.Tools = append(d.Tools, tool)
	d.toolsMap[tool.Name] = tool

	return nil
}

func (d *Doppelganger) MakeDecision(ctx context.Context, systemInstruction, userInstruction, model string) (string, error) {
	// Get Provider
	provider, err := d.providerGeneratorFunc(model)
	if err != nil {
		return "", err
	}

	// Construct history
	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemInstruction),
		llms.TextParts(llms.ChatMessageTypeHuman, userInstruction),
	}

	// Inject tool definitions available to the callout
	var toolDef []llms.Tool

	for _, tool := range d.Tools {
		toolDef = append(toolDef, llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		})
	}

	// Start inifinite loop
	for {
		res, err := provider.GenerateContent(ctx, messageHistory, llms.WithTools(toolDef))
		if err != nil {
			return "", err
		}

		// Parse response to check if tool calls requested
		var isToolCalled bool
		for _, choice := range res.Choices {
			for _, toolCall := range choice.ToolCalls {

				isToolCalled = true

				// Append tool_use to messageHistory
				aiResponse := llms.MessageContent{
					Role: llms.ChatMessageTypeAI,
					Parts: []llms.ContentPart{
						llms.ToolCall{
							ID:   toolCall.ID,
							Type: toolCall.Type,
							FunctionCall: &llms.FunctionCall{
								Name:      toolCall.FunctionCall.Name,
								Arguments: toolCall.FunctionCall.Arguments,
							},
						},
					},
				}
				messageHistory = append(messageHistory, aiResponse)

				// Call tools if requested
				toolResult, err := d.callTool(ctx, toolCall)
				if err != nil {
					return "", err
				}

				// Write back
				response := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolCall.ID,
							Name:       toolCall.FunctionCall.Name,
							Content:    toolResult,
						},
					},
				}

				messageHistory = append(messageHistory, response)
			}
		}

		if !isToolCalled {
			return res.Choices[0].Content, nil
		}
	}
}

func (d *Doppelganger) callTool(ctx context.Context, toolRequested llms.ToolCall) (string, error) {

	rt, exists := d.toolsMap[toolRequested.FunctionCall.Name]
	if !exists {
		return "", fmt.Errorf("invalid tool")
	}

	var params map[string]interface{}
	err := json.Unmarshal([]byte(toolRequested.FunctionCall.Arguments), &params)
	if err != nil {
		return "", err
	}

	result, err := rt.Execute(ctx, params)
	if err != nil {
		return "", err
	}

	resBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resBytes), nil

}
