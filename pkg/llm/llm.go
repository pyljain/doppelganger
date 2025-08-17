package llm

import (
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/openai"
)

func GetProvider(model string) (llms.Model, error) {
	if strings.HasPrefix(model, "gpt-4.1") {
		return openai.New(openai.WithModel(model))
	} else if strings.HasPrefix(model, "claude") {
		return anthropic.New(anthropic.WithModel(model))
	}

	return nil, fmt.Errorf("Model not found")
}
