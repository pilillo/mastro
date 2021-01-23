package featurestore

import (
	"fmt"

	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/featurestore/daos/elastic"
	"github.com/pilillo/mastro/featurestore/daos/mongo"
	"github.com/pilillo/mastro/utils/conf"
)

// available backends - lazy loaded singleton DAOs
var availableDAOs = map[string]func() abstract.FeatureSetDAOProvider{
	"mongo":   mongo.GetSingleton,
	"elastic": elastic.GetSingleton,
}

func selectDao(cfg *conf.Config) (abstract.FeatureSetDAOProvider, error) {
	if singletonDao, ok := availableDAOs[cfg.DataSourceDefinition.Type]; ok {
		// call singleton constructor on dao
		return singletonDao(), nil
	}
	return nil, fmt.Errorf("Impossible to find specified DAO connector %s", cfg.DataSourceDefinition.Type)
}
