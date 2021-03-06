package abstract

import "github.com/pilillo/mastro/utils/conf"

// AssetDAOProvider ... The interface each dao must implement
type AssetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Upsert(asset *Asset) error
	GetById(id string) (*Asset, error)
	GetByName(id string) (*Asset, error)
	SearchAssetsByTags(tags []string) (*[]Asset, error)
	ListAllAssets() (*[]Asset, error)
	CloseConnection()
}
