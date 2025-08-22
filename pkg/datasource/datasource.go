package datasource

import "context"

type DataSource interface {
	Connect(ctx context.Context, connectionString string) error
	Close(ctx context.Context) error
	Query(ctx context.Context, database, method, collection, query string) ([]string, error)
	Type() string
}
