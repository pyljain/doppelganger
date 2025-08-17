package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDataSource struct {
	client *mongo.Client
}

func NewMongoDataSource() *MongoDataSource {
	return &MongoDataSource{}
}

func (m *MongoDataSource) Type() string {
	return "mongo"
}

func (m *MongoDataSource) Connect(connectionString string) error {
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	m.client = client
	return nil
}

func (m *MongoDataSource) Close(ctx context.Context) error {
	err := m.client.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDataSource) Query(ctx context.Context, database, method, collection, query string) ([]string, error) {
	mc := m.client.Database(database).Collection(collection)

	var bsonObject bson.D
	err := json.Unmarshal([]byte(query), &bsonObject)
	if err != nil {
		return nil, err
	}

	switch method {
	case "find":
		cursor, err := mc.Find(ctx, bsonObject)
		if err != nil {
			return nil, err
		}

		defer cursor.Close(ctx)

		var records []string
		for cursor.Next(ctx) {
			var result bson.D
			if err := cursor.Decode(&result); err != nil {
				return nil, err
			}
			resBytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}

			records = append(records, string(resBytes))

		}

		return records, nil
	case "findOne":
		res := mc.FindOne(ctx, bsonObject)
		if res.Err() != nil {
			log.Printf("Error in findOne %s", res.Err())
			return nil, res.Err()
		}

		var result bson.D
		if err := res.Decode(&result); err != nil {
			return nil, err
		}

		resBytes, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		return []string{string(resBytes)}, nil

	}

	return nil, fmt.Errorf("method not supported")
}
