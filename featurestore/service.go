package featurestore

import (
	"fmt"
	"log"

	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/featurestore/daos/elastic"
	"github.com/pilillo/mastro/featurestore/daos/mongo"
	"github.com/pilillo/mastro/utils/conf"
	"github.com/pilillo/mastro/utils/date"
	"github.com/pilillo/mastro/utils/errors"
)

// FeatureSetService ... Group all service methods in a kind FeatureSetServiceType implementing the FeatureSetServiceInterface
var featureSetService FeatureSetServiceInterface = &FeatureSetServiceType{}

// FeatureSetServiceType ... Service Type
type FeatureSetServiceType struct{}

// FeatureSetServiceInterface ... Service Interface listing implemented methods
type FeatureSetServiceInterface interface {
	Init(cfg *conf.Config) *errors.RestErr
	CreateFeatureSet(fs abstract.FeatureSet) (*abstract.FeatureSet, *errors.RestErr)
	GetFeatureSetByID(fsID string) (*abstract.FeatureSet, *errors.RestErr)
	ListAllFeatureSets() (*[]abstract.FeatureSet, *errors.RestErr)
}

// selected dao for the featureSetService
var dao abstract.FeatureSetDAOProvider

// Mapping function -----
// available backends
var availableDAOs = map[string]abstract.FeatureSetDAOProvider{
	"mongo":   &mongo.MongoDAO{},
	"elastic": &elastic.ElasticDAO{},
}

// todo: remove pre-allocated dao and use a switch in the selectDao function to allocate just 1
func selectDao(cfg *conf.Config) (abstract.FeatureSetDAOProvider, error) {
	if dao, ok := availableDAOs[cfg.DataSourceDefinition.Type]; ok {
		return dao, nil
	}
	return nil, fmt.Errorf("Impossible to find specified DAO connector %s", cfg.DataSourceDefinition.Type)
}

// -----------------------

// Init ... Initializes the connector by validating the config and initializing the connection
func (s *FeatureSetServiceType) Init(cfg *conf.Config) *errors.RestErr {
	// select target DAO based on used connector
	// set a connector to the selected backend here
	var err error
	dao, err = selectDao(cfg)
	if err != nil {
		log.Panicln(err)
	}
	dao.Init(&cfg.DataSourceDefinition)
	return nil
}

// CreateFeatureSet ... Create a FeatureSet entry
func (s *FeatureSetServiceType) CreateFeatureSet(fs abstract.FeatureSet) (*abstract.FeatureSet, *errors.RestErr) {
	if restErr := fs.Validate(); restErr != nil {
		return nil, restErr
	}
	// set insert time to current date, then insert using selected dao
	fs.InsertedAt = date.GetNow()
	err := dao.Create(&fs)
	if err != nil {
		return nil, errors.GetBadRequestError(err.Error())
	}
	// what should we actually return of the newly inserted object?
	return &fs, nil
}

// GetFeatureSetByID ... Retrieves a FeatureSet
func (s *FeatureSetServiceType) GetFeatureSetByID(fsID string) (*abstract.FeatureSet, *errors.RestErr) {
	fset, err := dao.GetById(fsID)
	if err != nil {
		return nil, errors.GetNotFoundError(err.Error())
	}
	return fset, nil
}

/*
func (s *FeatureSetServiceType) GetFeatureSetByTagsAndTime() {

}
*/

// ListAllFeatureSets ... Retrieves all FeatureSets
func (s *FeatureSetServiceType) ListAllFeatureSets() (*[]abstract.FeatureSet, *errors.RestErr) {
	fset, err := dao.ListAllFeatureSets()
	if err != nil {
		return nil, errors.GetInternalServerError(err.Error())
	}
	return fset, nil
}
