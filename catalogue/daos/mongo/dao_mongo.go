package mongo

import (
	"sync"

	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/sources/mongo"
	"github.com/pilillo/mastro/utils/conf"
)

var once sync.Once
var instance *dao

type dao struct {
	Connector *mongo.Connector
}

// GetSingleton ... get an instance of the dao backend
func GetSingleton() abstract.AssetDAOProvider {
	// once.do is lazy, we use it to return an instance of the DAO
	once.Do(func() {
		instance = &dao{}
	})
	return instance
}

// Init ... Initialize connection to elastic search and target index
func (dao *dao) Init(def *conf.DataSourceDefinition) {
	return
}

// Create ... Create asset on ES
func (dao *dao) Create(fs *abstract.Asset) error {
	return nil
}

// ListAllFeatureSets ... Return all assets in index
func (dao *dao) ListAllAssets() (*[]abstract.Asset, error) {
	return nil, nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.Asset, error) {
	return nil, nil
}

// CloseConnection ... Terminates the connection to ES for the DAO
func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}
