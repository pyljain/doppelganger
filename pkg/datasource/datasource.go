package datasource

import "context"

type DataSource interface {
	Connect(connectionString string) error
	Close(ctx context.Context) error
	Query(ctx context.Context, database, method, collection, query string) ([]string, error)
	Type() string
}

// Mongo method = FindOne, Find, Aggregate, collection=collection
// Postgres method = Select, collection="", query="query"
