package tool

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	tt := []struct {
		description       string
		query             string
		queryReturnsError bool
		params            map[string]interface{}
		expectedError     bool
		expectedResult    []string
	}{
		{
			description:       "When query does not have correct go template syntax",
			query:             "{{ .test }",
			params:            map[string]interface{}{},
			queryReturnsError: false,
			expectedError:     true,
			expectedResult:    nil,
		},
		{
			description:       "When LLM does not pass the right params should return error",
			query:             "This is a {{ .test }}",
			params:            map[string]interface{}{},
			queryReturnsError: false,
			expectedError:     true,
			expectedResult:    nil,
		},
		{
			description: "When LLM passes the right params should not return error",
			query:       "This is the sample code: {{ .code }}",
			params: map[string]interface{}{
				"code": "abc",
			},
			queryReturnsError: false,
			expectedError:     false,
			expectedResult:    []string{"This is the sample code: abc"},
		},
		{
			description: "When query execution fails an error should be returned",
			query:       "This is the sample code: {{ .code }}",
			params: map[string]interface{}{
				"code": "abc",
			},
			queryReturnsError: true,
			expectedError:     true,
			expectedResult:    nil,
		},
	}

	for _, test := range tt {
		ctx := context.Background()
		t.Run(test.description, func(t *testing.T) {
			dst := DataSourceTool{
				Source: &mockDatasource{
					returnError: test.queryReturnsError,
				},
				Name:        "gcs",
				Description: "gcs",
				Parameters:  map[string]interface{}{},
				Collection:  "",
				Database:    "",
				Method:      "get",
				Query:       test.query,
			}

			res, err := dst.Execute(ctx, test.params)
			if test.expectedError {
				require.NotNil(t, err)
				return
			}

			require.Equal(t, test.expectedResult, res)
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
