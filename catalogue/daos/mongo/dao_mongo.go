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
	"gopkg.in/mgo.v2/bson"
)

// AssetMongoDao ... managed resource
type AssetMongoDao struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// asset publication datetime
	PublishedOn time.Time `bson:"published-on"`
	// name of the asset
	Name string `bson:"name"`
	// description of the asset
	Description string `bson:"description"`
	// the list of assets this depends on
	DependsOn []string `bson:"depends-on"`
	// the actual asset metadata
	Metadata MetadataMongoDao `bson:"metadata"`
	// tags are flags used to simplify asset search
	Tags []string `bson:"tags,omitempty"`
}

// MetadataMongoDao ... Asset metadata
type MetadataMongoDao struct {
	// asset type - refers to the type.name
	TypeName string `bson:"type"`
	// asset attributes, key-val list
	Labels map[string]string `bson:"labels,omitempty"`
}

// TypeMongoDao ... Asset type
type TypeMongoDao struct {
	Name string `bson:"name"`
	ID   string `bson:"id"`
}

// ItemMongoDao ... Managed item
type ItemMongoDao struct {
	LastDiscoveredAt time.Time
	Asset            AssetMongoDao
}

func convertAssetDTOtoDAO(as *abstract.Asset) *AssetMongoDao {
	asmd := &AssetMongoDao{}

	// id not set at the time of insert (DTO->DAO)
	asmd.PublishedOn = as.PublishedOn
	asmd.Name = as.Name
	asmd.Description = as.Description
	asmd.DependsOn = as.DependsOn

	asmdMetadata := MetadataMongoDao{}
	asmdMetadata.TypeName = as.Metadata.TypeName
	asmdMetadata.Labels = as.Metadata.Labels
	asmd.Metadata = asmdMetadata

	asmd.Tags = as.Tags
	return asmd
}

func convertAssetDAOtoDTO(asmd *AssetMongoDao) *abstract.Asset {
	as := &abstract.Asset{}

	as.PublishedOn = asmd.PublishedOn
	as.Name = asmd.Name
	as.Description = asmd.Description
	as.DependsOn = asmd.DependsOn

	asMetadata := abstract.Metadata{}
	asMetadata.TypeName = asmd.Metadata.TypeName
	asMetadata.Labels = asmd.Metadata.Labels
	as.Metadata = asMetadata

	as.Tags = asmd.Tags

	return as
}

func convertAllAssets(inAssets *[]AssetMongoDao) []abstract.Asset {
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
	dao.Connector.InitConnection(def)
}

// Create ... Create asset
func (dao *dao) Create(as *abstract.Asset) error {
	asmd := convertAssetDTOtoDAO(as)

	bsonVal, err := bson.Marshal(asmd)
	if err != nil {
		return err
	}

	// insert
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := dao.Connector.Collection.InsertOne(ctx, bsonVal)
	if err != nil {
		return fmt.Errorf("Error while creating asset :: %v", err)
	}
	id := res.InsertedID
	log.Printf("Inserted Asset %d", id)
	return nil
}

// GetById ... Retrieve document by given id
func (dao *dao) GetById(id string) (*abstract.Asset, error) {
	var result AssetMongoDao

	filter := bson.M{"_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := dao.Connector.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving asset :: %v", err)
	}

	return convertAssetDAOtoDTO(&result), nil
}

// ListAllFeatureSets ... Return all assets in index
func (dao *dao) ListAllAssets() (*[]abstract.Asset, error) {
	var assets []AssetMongoDao

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cursor, err := dao.Connector.Collection.Find(
		ctx,
		bson.D{{}},
	)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &assets); err != nil {
		return nil, err
	}

	var resultAssets []abstract.Asset = convertAllAssets(&assets)
	return &resultAssets, nil
}

// CloseConnection ... Terminates the connection to ES for the DAO
func (dao *dao) CloseConnection() {
	dao.Connector.CloseConnection()
}
