package datasource

type MongoDataSource struct {
}

func NewMongoDataSource() *MongoDataSource {
	return &MongoDataSource{}
}

func (m *MongoDataSource) Type() string {
	return "mongo"
}

func (m *MongoDataSource) Connect(connectionString string) error {
	return nil
}

func (m *MongoDataSource) Close() error {
	return nil
}

func (m *MongoDataSource) Query(query string) ([]Record, error) {
	return nil, nil
}
