package llm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProvider(t *testing.T) {
	tt := []struct {
		description   string
		model         string
		expectedError error
	}{
		{
			description:   "When a valid OpenAI model prefix is passed, returns a pointer to provider",
			model:         "gpt-4.1",
			expectedError: nil,
		},
		{
			description:   "When a valid Anthrpoic model prefix is passed, returns a pointer to provider",
			model:         "claude-sonnet-4-20250514",
			expectedError: nil,
		},
		{
			description:   "When an invalid model prefix is passed, returns an error",
			model:         "gemini2.5-pro",
			expectedError: ErrModelNotFound,
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			_, err := GetProvider(test.model)
			require.Equal(t, test.expectedError, err)
		})
	}
}
