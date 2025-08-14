package tool

import "doppelganger/pkg/datasource"

type DataSourceTool struct {
	Source      datasource.DataSource
	Name        string
	Description string
	Parameters  map[string]interface{}
	Collection  string
	Method      string
	Query       string
}

func (dst *DataSourceTool) Execute(params map[string]interface{}) (string, error) {
	dst.Source.Query()
}

/*



 */
