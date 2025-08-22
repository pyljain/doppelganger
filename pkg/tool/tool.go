package tool

import (
	"bytes"
	"context"
	"doppelganger/pkg/datasource"
	"sync"
	"text/template"
)

type DataSourceTool struct {
	Source         datasource.DataSource
	Name           string
	Description    string
	Parameters     map[string]interface{}
	Database       string
	Collection     string
	Method         string
	Query          string
	parsedTemplate *template.Template
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

func (dst *DataSourceTool) Execute(ctx context.Context, params map[string]interface{}) ([]string, error) {
	if dst.parsedTemplate == nil {
		tmpl, err := template.New("test").Option("missingkey=error").Parse(dst.Query)
		if err != nil {
			return nil, err
		}

		dst.parsedTemplate = tmpl
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	err := dst.parsedTemplate.Execute(buf, params)
	if err != nil {
		return nil, err
	}

	records, err := dst.Source.Query(ctx, dst.Database, dst.Method, dst.Collection, buf.String())
	if err != nil {
		return nil, err
	}

	return records, nil
}
