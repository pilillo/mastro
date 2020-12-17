package catalogue

import (
	"log"

	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/utils/conf"
	"github.com/pilillo/mastro/utils/errors"
)

// Service ... Service Interface listing implemented methods
type Service interface {
	Init(cfg *conf.Config) *errors.RestErr
	CreateAsset(fs abstract.Asset) (*abstract.Asset, *errors.RestErr)
	GetAssetByID(assetID string) (*abstract.Asset, *errors.RestErr)
	ListAllAssets() (*[]abstract.Asset, *errors.RestErr)
}

// assetServiceType ... Service Type
type assetServiceType struct{}

// assetService ... Group all service methods in a kind FeatureSetServiceType implementing the FeatureSetService
var assetService Service = &assetServiceType{}

// selected dao for the featureSetService
var dao abstract.AssetDAOProvider

// Init ... initializes the service
func (s *assetServiceType) Init(cfg *conf.Config) *errors.RestErr {
	// select target DAO based on used connector
	// set a connector to the selected backend here
	var err error
	// select dao using mapping function in same package
	dao, err = selectDao(cfg)
	if err != nil {
		log.Panicln(err)
	}
	dao.Init(&cfg.DataSourceDefinition)
	return nil
}

// CreateAsset ... Adds and asset description
func (s *assetServiceType) CreateAsset(asset abstract.Asset) (*abstract.Asset, *errors.RestErr) {
	return nil, nil
}

// GetAssetById ... Retrieves an asset by its name ID
func (s *assetServiceType) GetAssetByID(assetID string) (*abstract.Asset, *errors.RestErr) {
	asset, err := dao.GetById(assetID)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return asset, nil
}

// ListAllAssets ... Retrieves all stored assets
func (s *assetServiceType) ListAllAssets() (*[]abstract.Asset, *errors.RestErr) {
	asset, err := dao.ListAllAssets()
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	return asset, nil
}
