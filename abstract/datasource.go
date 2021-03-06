package abstract

import "github.com/pilillo/mastro/utils/conf"

// ConnectorProvider ... The interface each connector must implement
type ConnectorProvider interface {
	ValidateDataSourceDefinition(*conf.DataSourceDefinition) error
	InitConnection(*conf.DataSourceDefinition)
	CloseConnection()
}
