package conf

// DataSourceDefinition ... connection details for a data source connector
type DataSourceDefinition struct {
	Name     string  `yaml:"name"`
	Type     string  `yaml:"type"`
	Settings Details `yaml:"settings,omitempty"`
}
