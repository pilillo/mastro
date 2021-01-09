package abstract

import "github.com/pilillo/mastro/utils/conf"

// FeatureSetDAOProvider ... The interface each dao must implement
type FeatureSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(fs *FeatureSet) error
	GetById(id string) (*FeatureSet, error)
	ListAllFeatureSets() (*[]FeatureSet, error)
	CloseConnection()
}
