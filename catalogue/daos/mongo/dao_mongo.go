package mongo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/sources/mongo"
	"github.com/pilillo/mastro/utils/conf"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type assetMongoDao struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// asset discovery datetime
	LastDiscoveredAt time.Time `bson:"last-discovered-at"`
	// asset publication datetime
	PublishedOn time.Time `bson:"published-on"`
	// name of the asset
	Name string `bson:"name"`
	// description of the asset
	Description string `bson:"description"`
	// the list of assets this depends on
	DependsOn []string `bson:"depends-on"`
	// asset type
	Type abstract.AssetType `bson:"type"`
	// asset labels
	Labels map[string]interface{} `bson:"labels,omitempty"`
	// tags are flags used to simplify asset search
	Tags []string `bson:"tags,omitempty"`
}

func convertAssetDTOtoDAO(as *abstract.Asset) *assetMongoDao {
	asmd := &assetMongoDao{}

	// id not set at the time of insert (DTO->DAO)
	// however id is set if we are updating an existing asset
	asmd.LastDiscoveredAt = as.LastDiscoveredAt
	asmd.PublishedOn = as.PublishedOn
	asmd.Name = as.Name
	asmd.Description = as.Description
	asmd.DependsOn = as.DependsOn

	asmd.Type = as.Type
	asmd.Labels = as.Labels

	asmd.Tags = as.Tags
	return asmd
}

func convertAssetDAOtoDTO(asmd *assetMongoDao) *abstract.Asset {
	as := &abstract.Asset{}

	as.LastDiscoveredAt = asmd.LastDiscoveredAt
	as.PublishedOn = asmd.PublishedOn
	as.Name = asmd.Name
	as.Description = asmd.Description
	as.DependsOn = asmd.DependsOn

	as.Type = asmd.Type
	as.Labels = asmd.Labels

	as.Tags = asmd.Tags

	return as
}

func convertAllAssets(inAssets *[]assetMongoDao) []abstract.Asset {
	var assets []abstract.Asset
	for _, element := range *inAssets {
		assets = append(assets, *convertAssetDAOtoDTO(&element))
	}
	return assets
}

var once sync.Once
var instance *dao

type dao struct {
	Connector *mongo.Connector
}

var timeout = 5 * time.Second

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
	dao.Connector = mongo.NewMongoConnector()
	if err := dao.Connector.ValidateDataSourceDefinition(def); err != nil {
		panic(err)
	}
	dao.Connector.InitConnection(def)
}

// Upsert ... Upsert asset
func (dao *dao) Upsert(as *abstract.Asset) error {
	asmd := convertAssetDTOtoDAO(as)

	bsonVal, err := bson.Marshal(asmd)
	if err != nil {
		return err
	}

	// https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents
	// When working with the ReplaceOne function, update operators such as $set cannot be used since
	// it is a complete replace rather than an update of particular fields.
	// replace all fields in a document while maintaining the id of that document
	// https://stackoverflow.com/questions/59311020/upsert-not-working-when-using-updateone-with-the-mongodb-golang-driver

	opts := options.Replace().SetUpsert(true) //.Update().SetUpsert(true)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"name": as.Name}
	result, err := dao.Connector.Collection.ReplaceOne(ctx, filter, bsonVal, opts)
	if err != nil {
		return fmt.Errorf("Error while upserting asset :: %v", err)
	}
	id := result.UpsertedID
	log.Printf("Upserted Asset %d", id)

	return nil
}

func (dao *dao) getOneDocumentUsingFilter(filter interface{}) (*abstract.Asset, error) {
	var result assetMongoDao
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := dao.Connector.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving asset :: %v", err)
	}
	return convertAssetDAOtoDTO(&result), nil
}

func (dao *dao) getAnyDocumentUsingFilter(filter interface{}) (*[]abstract.Asset, error) {
	var assets []assetMongoDao

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cursor, err := dao.Connector.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &assets); err != nil {
		return nil, err
	}

	var resultAssets []abstract.Asset = convertAllAssets(&assets)
	return &resultAssets, nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.Asset, error) {
	filter := bson.M{"_id": id}
	return dao.getOneDocumentUsingFilter(filter)
}

// GetByName ... Retrieve document by given name
func (dao *dao) GetByName(name string) (*abstract.Asset, error) {
	filter := bson.D{{"name", name}}
	return dao.getOneDocumentUsingFilter(filter)
}

func (dao *dao) SearchAssetsByTags(tags []string) (*[]abstract.Asset, error) {
	// https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline
	// https://docs.mongodb.com/manual/tutorial/query-arrays/#match-an-array
	// find all docs whose tags field contains all the elements provided as tags []string in input
	// without regard of the order
	filter := bson.M{"tags": bson.M{"$all": tags}}
	return dao.getAnyDocumentUsingFilter(filter)
}

// ListAllFeatureSets ... Return all assets in index
func (dao *dao) ListAllAssets() (*[]abstract.Asset, error) {
	filter := bson.D{{}}
	return dao.getAnyDocumentUsingFilter(filter)
}

// CloseConnection ... Terminates the connection to ES for the DAO
func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}
