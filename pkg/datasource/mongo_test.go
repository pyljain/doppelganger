package datasource

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestConnectAndDisconnect(t *testing.T) {
	connectionString := "mongodb://localhost:27017"
	ctx := context.Background()
	m := NewMongoDataSource()

	err := m.Connect(ctx, connectionString)
	require.Nil(t, err)

	err = m.Close(ctx)
	require.Nil(t, err)
}

func TestQuery(t *testing.T) {
	tt := []struct {
		description     string
		collection      string
		recordsToInsert []bson.M
		query           string
		method          string
		expectError     bool
		recordsReturned int
	}{
		{
			description: "When find is called and matching records are present, returns records without any error",
			collection:  "test",
			recordsToInsert: []bson.M{
				{
					"name": "one",
					"type": "record",
				},
				{
					"name": "two",
					"type": "record",
				},
			},
			query:           "{ \"type\": \"record\" }",
			method:          "find",
			expectError:     false,
			recordsReturned: 2,
		},
		{
			description: "When findOne is called and matching records are present, returns the matched record without any error",
			collection:  "test",
			recordsToInsert: []bson.M{
				{
					"name": "one",
					"type": "record",
				},
				{
					"name": "two",
					"type": "record",
				},
			},
			query:           "{ \"name\": \"one\" }",
			method:          "findOne",
			expectError:     false,
			recordsReturned: 1,
		},
		{
			description: "When an unsupported method is called, an error is returned saying method not supported",
			collection:  "test",
			recordsToInsert: []bson.M{
				{
					"name": "one",
					"type": "record",
				},
				{
					"name": "two",
					"type": "record",
				},
			},
			query:           "{ \"name\": \"one\" }",
			method:          "findTwo",
			expectError:     true,
			recordsReturned: 0,
		},
		{
			description: "When the query is not formed correctly, an error is returned",
			collection:  "test",
			recordsToInsert: []bson.M{
				{
					"name": "one",
					"type": "record",
				},
				{
					"name": "two",
					"type": "record",
				},
			},
			query:           "{ \"name: \"one\" }",
			method:          "findOne",
			expectError:     true,
			recordsReturned: 0,
		},
	}

	m := NewMongoDataSource()
	ctx := context.Background()
	err := m.Connect(ctx, "mongodb://localhost:27017")
	require.Nil(t, err)
	defer m.Close(ctx)

	for _, test := range tt {
		collection := uuid.New().String()

		_, err := m.client.Database("test").Collection(collection).InsertMany(ctx, test.recordsToInsert)
		require.Nil(t, err)

		records, err := m.Query(ctx, "test", test.method, collection, test.query)
		if !test.expectError {
			require.Nil(t, err)
			require.Equal(t, test.recordsReturned, len(records))
		} else {
			require.NotNil(t, err)
		}

		err = m.client.Database("test").Collection(collection).Drop(ctx)
		require.Nil(t, err)
	}
}
