package sources

import (
	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/sources/elastic"
	"github.com/pilillo/mastro/sources/mongo"
	"github.com/pilillo/mastro/sources/postgres"
)

// AvailableConnectors ... available connectors that can be instantiated
var AvailableConnectors = map[string]abstract.ConnectorProvider{
	"postgres": &postgres.Connector,
	"elastic":  &elastic.Connector,
	"mongo":    &mongo.Connector,
}
