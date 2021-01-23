# Mastro
## Data Catalogue
Data providers can describe and publish data using a shared definition format.
Consequently, data definitions can be crawled from networked and distributed file systems, as well as directly published to a common endpoint.

### Catalogue API
A Catalogue service endpoint implements the following interface:

```go
type Service interface {
	Init(cfg *conf.Config) *errors.RestErr
	UpsertAssets(assets *[]abstract.Asset) (*[]abstract.Asset, *errors.RestErr)
	GetAssetByID(assetID string) (*abstract.Asset, *errors.RestErr)
	GetAssetByName(name string) (*abstract.Asset, *errors.RestErr)
	SearchAssetsByTags(tags []string) (*[]abstract.Asset, *errors.RestErr)
	ListAllAssets() (*[]abstract.Asset, *errors.RestErr)
}
```

This can be easily mapped to a DAO backend:
```go
type AssetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Upsert(asset *Asset) error
	GetById(id string) (*Asset, error)
	GetByName(id string) (*Asset, error)
	SearchAssetsByTags(tags []string) (*[]Asset, error)
	ListAllAssets() (*[]Asset, error)
	CloseConnection()
}
```

Have a look at `catalogue/daos/*` for example implementations.