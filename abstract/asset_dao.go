package abstract

import "github.com/pilillo/mastro/utils/conf"

// AssetDAOProvider ... The interface each dao must implement
type AssetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(asset *Asset) error
	GetById(id string) (*Asset, error)
	ListAllAssets() (*[]Asset, error)
	CloseConnection()
}
