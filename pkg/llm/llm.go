package llm

import (
	"errors"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/openai"
)

var ErrModelNotFound = errors.New("model not found")

func GetProvider(model string) (llms.Model, error) {
	if strings.HasPrefix(model, "gpt") {
		return openai.New(openai.WithModel(model))
	} else if strings.HasPrefix(model, "claude") {
		return anthropic.New(anthropic.WithModel(model))
	}

	return nil, ErrModelNotFound
}
