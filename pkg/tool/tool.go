package tool

import (
	"bytes"
	"context"
	"doppelganger/pkg/datasource"
	"log"
	"text/template"
)

type DataSourceTool struct {
	Source      datasource.DataSource
	Name        string
	Description string
	Parameters  map[string]interface{}
	Database    string
	Collection  string
	Method      string
	Query       string
}

func (dst *DataSourceTool) Execute(ctx context.Context, params map[string]interface{}) ([]string, error) {
	tmpl, err := template.New("test").Parse(dst.Query)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, params)
	if err != nil {
		return nil, err
	}
	log.Printf("Query is %s", buf.String())

	records, err := dst.Source.Query(ctx, dst.Database, dst.Method, dst.Collection, buf.String())
	if err != nil {
		return nil, err
	}

	return records, nil
}
