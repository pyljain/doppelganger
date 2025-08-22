package datasource

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

type GCS struct {
	bucket *storage.BucketHandle
	client *storage.Client
}

func NewGCS(ctx context.Context) *GCS {
	return &GCS{}
}

func (g *GCS) Connect(ctx context.Context, bucketName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	g.bucket = bucket
	g.client = client

	return nil
}

func (g *GCS) Close(ctx context.Context) error {
	err := g.client.Close()
	return err
}

func (g *GCS) Type() string {
	return "gcs"
}

func (g *GCS) Query(ctx context.Context, database, method, collection, query string) ([]string, error) {
	switch method {
	case "list":
		objectNames := []string{}
		iterator := g.bucket.Objects(ctx, nil)
		for {
			attr, err := iterator.Next()
			if err != nil {
				break
			}

			objectNames = append(objectNames, attr.Name)
		}

		return objectNames, nil
	case "get":
		objectHandle := g.bucket.Object(query)
		reader, err := objectHandle.NewReader(ctx)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		objectBytes, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		return []string{string(objectBytes)}, nil

	}

	return nil, fmt.Errorf("method not supported")
}
