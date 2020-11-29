package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pilillo/mastro/abstract"
	mastromongo "github.com/pilillo/mastro/sources/mongo"
	"github.com/pilillo/mastro/utils/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FeatureSetMongoDao ... DAO for the FeatureSet in Mongo
type FeatureSetMongoDao struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	InsertedAt  time.Time          `bson:"inserted-at,omitempty"`
	Version     string             `bson:"version,omitempty"`
	Features    []FeatureMongoDao  `bson:"features,omitempty"`
	Description string             `bson:"description,omitempty"`
	Labels      map[string]string  `bson:"labels,omitempty"`
}

// VersionMongoDao ... definition of version for a feature set
type VersionMongoDao struct{}

// FeatureMongoDao ... a named variable with a data type
type FeatureMongoDao struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Value    interface{}        `bson:"value,omitempty"`
	DataType string             `bson:"data-type,omitempty"`
}

type MongoDAO struct {
	Connector *mastromongo.MongoConnector
}

func convertFeatureDTOtoDAO(f *abstract.Feature) *FeatureMongoDao {
	fmd := &FeatureMongoDao{}

	//fmd.ID = f.ID // not set at time of insert
	fmd.Name = f.Name
	fmd.Value = f.Value
	fmd.DataType = f.DataType

	return fmd
}

func convertFeatureSetDTOtoDAO(fs *abstract.FeatureSet) *FeatureSetMongoDao {
	fsmd := &FeatureSetMongoDao{}

	//fsmd.ID = fs.ID // not set at time of insert
	fsmd.InsertedAt = fs.InsertedAt
	fsmd.Version = fs.Version

	var feats []FeatureMongoDao
	for _, element := range fs.Features {
		feats = append(feats, *convertFeatureDTOtoDAO(&element))
	}
	fsmd.Features = feats

	fsmd.Description = fs.Description
	fsmd.Labels = fs.Labels

	return fsmd
}

func convertFeatureDAOToDTO(fmd *FeatureMongoDao) *abstract.Feature {
	f := &abstract.Feature{}

	//f.ID = fmd.ID.String() // set it in DAO, propagate to DTO?
	f.Name = fmd.Name
	f.Value = fmd.Value
	f.DataType = fmd.DataType

	return f
}

func convertAllFeatureSets(inFeats *[]FeatureSetMongoDao) []abstract.FeatureSet {
	var feats []abstract.FeatureSet
	for _, element := range *inFeats {
		feats = append(feats, *convertFeatureSetDAOToDTO(&element))
	}
	return feats
}

func convertAllFeatures(inFeats *[]FeatureMongoDao) []abstract.Feature {
	var feats []abstract.Feature
	for _, element := range *inFeats {
		feats = append(feats, *convertFeatureDAOToDTO(&element))
	}
	return feats
}

func convertFeatureSetDAOToDTO(fsmd *FeatureSetMongoDao) *abstract.FeatureSet {
	fs := &abstract.FeatureSet{}

	//fs.ID = fsmd.ID.String() // set it in DAO, propagate to DTO?
	fs.InsertedAt = fsmd.InsertedAt
	fs.Version = fsmd.Version

	fs.Features = convertAllFeatures(&fsmd.Features)
	fs.Description = fsmd.Description
	fs.Labels = fsmd.Labels

	return fs
}

//collection := client.Database("testing").Collection("numbers")

func (dao *MongoDAO) Init(def *conf.DataSourceDefinition) {
	// create mongo connector
	dao.Connector = &mastromongo.MongoConnector{}
	// init mongo connector
	dao.Connector.InitConnection(def)
}

func (dao *MongoDAO) CloseConnection() {
	dao.Connector.CloseConnection()
}

func (dao *MongoDAO) Create(fs *abstract.FeatureSet) error {
	// convert DTO to DAO
	//bsonVal := bson.M{"name": "pi", "value": 3.14159}
	fsmd := convertFeatureSetDTOtoDAO(fs)

	bsonVal, err := bson.Marshal(fsmd)
	if err != nil {
		return err
	}

	// insert
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := dao.Connector.Collection.InsertOne(ctx, bsonVal)
	if err != nil {
		return fmt.Errorf("Error while creating feature set :: %v", err)
	}
	id := res.InsertedID
	log.Println("Inserted FeatureSet %d", id)
	return nil
}

func (dao *MongoDAO) GetById(id string) (*abstract.FeatureSet, error) {
	var result FeatureSetMongoDao
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := dao.Connector.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving feature set :: %v", err)
	}

	// convert DAO to DTO

	return nil, nil
}

func (dao *MongoDAO) ListAllFeatureSets() (*[]abstract.FeatureSet, error) {
	var features []FeatureSetMongoDao
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := dao.Connector.Collection.Find(
		ctx,
		bson.D{{}},
	)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &features); err != nil {
		return nil, err
	}

	var resultFeats []abstract.FeatureSet = convertAllFeatureSets(&features)
	return &resultFeats, nil
}
