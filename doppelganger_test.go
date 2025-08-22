package doppelganger

import (
	"context"
	"doppelganger/pkg/tool"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/llms"
)

func TestRegisterTool(t *testing.T) {
	tt := []struct {
		description   string
		toolDef       tool.DataSourceTool
		expectedError bool
	}{
		{
			description: "When a valid tool definition is passed, it should get successful added and no error should be returned",
			toolDef: tool.DataSourceTool{
				Name:        "mockFunction",
				Description: "A function to interact with the Mock tool",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"code": map[string]any{
							"type": "string",
						},
					},
				},
				Database:   "",
				Collection: "",
				Query:      "",
				Source:     &mockDatasource{},
			},
			expectedError: false,
		},
		{
			description: "When an invalid tool definition is passed, it should throw an error",
			toolDef: tool.DataSourceTool{
				Name:        "mockFunction",
				Description: "A function to interact with the Mock tool",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"code": "",
					},
				},
				Database:   "",
				Collection: "",
				Query:      "",
				Source:     &mockDatasource{},
			},
			expectedError: true,
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			dg := New()
			err := dg.RegisterTool(test.toolDef)
			if test.expectedError {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
		})
	}
}

type mockDatasource struct {
	returnError bool
}

func (m *mockDatasource) Connect(ctx context.Context, connectionString string) error {
	return nil
}

func (m *mockDatasource) Close(ctx context.Context) error {
	return nil
}

func (m *mockDatasource) Query(ctx context.Context, database, method, collection, query string) ([]string, error) {
	mockError := errors.New("Mock Error")
	if m.returnError == true {
		return nil, mockError
	}
	return []string{query}, nil
}

func (m *mockDatasource) Type() string {
	return "mock"
}

func TestMakeDecision(t *testing.T) {
	tt := []struct {
		description           string
		tools                 []tool.DataSourceTool
		providerGeneratorFunc ProviderGeneratorFunc
		expectedError         bool
	}{
		{
			description: "when no tools are added llm should respond with a string message",
			tools:       nil,
			providerGeneratorFunc: func(model string) (llms.Model, error) {
				return &mockProvider{
					responses: []*llms.ContentResponse{
						{
							Choices: []*llms.ContentChoice{
								{
									Content: "LLMs response",
								},
							},
						},
					},
					err: nil,
				}, nil
			},
			expectedError: false,
		},
		{
			description: "when tools are added and the LLM requests a tool call, it should be executed without error.",
			tools: []tool.DataSourceTool{
				{
					Name:        "mockFunction",
					Description: "A function to interact with the Mock tool",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"code": map[string]any{
								"type": "string",
							},
						},
					},
					Database:   "",
					Collection: "",
					Query:      "",
					Source:     &mockDatasource{},
				},
			},
			providerGeneratorFunc: func(model string) (llms.Model, error) {
				return &mockProvider{
					responses: []*llms.ContentResponse{
						{
							Choices: []*llms.ContentChoice{
								{
									ToolCalls: []llms.ToolCall{
										{
											ID: "123",
											FunctionCall: &llms.FunctionCall{
												Name:      "mockFunction",
												Arguments: "{ \"code\": \"abc\" }",
											},
										},
									},
								},
							},
						},
						{
							Choices: []*llms.ContentChoice{
								{
									Content: "abc",
								},
							},
						},
					},
					err: nil,
				}, nil
			},
			expectedError: false,
		},
		{
			description: "when llm returns error should return error",
			tools:       nil,
			providerGeneratorFunc: func(model string) (llms.Model, error) {
				return &mockProvider{
					responses: []*llms.ContentResponse{
						{
							Choices: []*llms.ContentChoice{
								{
									Content: "abc",
								},
							},
						},
					},
					err: fmt.Errorf("random error"),
				}, nil
			},
			expectedError: true,
		},
		{
			description: "when invalid provider passed returns error",
			tools:       nil,
			providerGeneratorFunc: func(model string) (llms.Model, error) {
				return nil, fmt.Errorf("invalid provider")
			},
			expectedError: true,
		},
		{
			description: "when arguments are passed as invalid json should return error",
			tools: []tool.DataSourceTool{
				{
					Name:        "mockFunction",
					Description: "A function to interact with the Mock tool",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"code": map[string]any{
								"type": "string",
							},
						},
					},
					Database:   "",
					Collection: "",
					Query:      "",
					Source:     &mockDatasource{},
				},
			},
			providerGeneratorFunc: func(model string) (llms.Model, error) {
				return &mockProvider{
					responses: []*llms.ContentResponse{
						{
							Choices: []*llms.ContentChoice{
								{
									ToolCalls: []llms.ToolCall{
										{
											ID: "123",
											FunctionCall: &llms.FunctionCall{
												Name:      "mockFunction",
												Arguments: "{ \"code\": }",
											},
										},
									},
								},
							},
						},
					},
					err: nil,
				}, nil
			},
			expectedError: true,
		},
		{
			description: "when cant execute tool returns error",
			tools: []tool.DataSourceTool{
				{
					Name:        "mockFunction",
					Description: "A function to interact with the Mock tool",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"code": map[string]any{
								"type": "string",
							},
						},
					},
					Database:   "",
					Collection: "",
					Query:      "{{ .co }",
					Source:     &mockDatasource{},
				},
			},
			providerGeneratorFunc: func(model string) (llms.Model, error) {
				return &mockProvider{
					responses: []*llms.ContentResponse{
						{
							Choices: []*llms.ContentChoice{
								{
									ToolCalls: []llms.ToolCall{
										{
											ID: "123",
											FunctionCall: &llms.FunctionCall{
												Name:      "mockFunction",
												Arguments: "{ \"code\": \"abc\" }",
											},
										},
									},
								},
							},
						},
					},
					err: nil,
				}, nil
			},
			expectedError: true,
		},
		{
			description: "when invalid tool name is passed returns error",
			tools: []tool.DataSourceTool{
				{
					Name:        "mockFunction",
					Description: "A function to interact with the Mock tool",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"code": map[string]any{
								"type": "string",
							},
						},
					},
					Database:   "",
					Collection: "",
					Query:      "",
					Source:     &mockDatasource{},
				},
			},
			providerGeneratorFunc: func(model string) (llms.Model, error) {
				return &mockProvider{
					responses: []*llms.ContentResponse{
						{
							Choices: []*llms.ContentChoice{
								{
									ToolCalls: []llms.ToolCall{
										{
											ID: "123",
											FunctionCall: &llms.FunctionCall{
												Name:      "functionMock",
												Arguments: "{ \"code\": \"abc\" }",
											},
										},
									},
								},
							},
						},
					},
					err: nil,
				}, nil
			},
			expectedError: true,
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			d := New()
			d.providerGeneratorFunc = test.providerGeneratorFunc
			ctx := context.Background()

			// Register testing tools
			for _, tl := range test.tools {
				err := d.RegisterTool(tl)
				require.Nil(t, err)
			}

			res, err := d.MakeDecision(ctx, "abc", "efg", "mock")
			if test.expectedError {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.NotEqual(t, "", res)
		})
	}
}

type mockProvider struct {
	responses []*llms.ContentResponse
	err       error
	counter   int
}

func (m *mockProvider) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	response := m.responses[m.counter]
	m.counter += 1
	return response, m.err
}

func (m *mockProvider) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return "", fmt.Errorf("not supported")
}
