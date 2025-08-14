package datasource

type Record map[string]interface{}

type DataSource interface {
	Connect(connectionString string) error
	Close() error
	Query(query string) ([]Record, error)
	Type() string
}
