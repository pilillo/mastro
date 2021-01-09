package mongo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pilillo/mastro/utils/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var requiredFields = map[string]string{
	"username": "username",
	"password": "password",
	"host":     "host",
	// surely needed the DB and the target collection
	"database":   "database",
	"collection": "collection",
}

// NewMongoConnector ... Factory
func NewMongoConnector() *Connector {
	return &Connector{}
}

// todo: find a way not to export this struct outside

// Connector ... struct containing info on how to connect to a mongo db
type Connector struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

// ValidateDataSourceDefinition ... validates the provided data source definition
func (c *Connector) ValidateDataSourceDefinition(def *conf.DataSourceDefinition) error {
	// check all required fields are available
	var missingFields []string
	for _, reqvalue := range requiredFields {
		if _, exist := def.Settings.Values[reqvalue]; !exist {
			missingFields = append(missingFields, reqvalue)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("The following %d fields are missing from the data source configuration: %s", len(missingFields), strings.Join(missingFields[:], ","))
	}

	log.Println("Successfully validated data source definition")
	return nil
}

// InitConnection ... Instantiate the connection with the remote DB
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {

	// todo: mongo connection string varies a lot, maybe just pass that from a secret rather than composing it here??
	connectionString := fmt.Sprintf(
		"mongodb://%s:%s@%s",
		def.Settings.Values[requiredFields["username"]],
		def.Settings.Values[requiredFields["password"]],
		def.Settings.Values[requiredFields["host"]],
	)
	log.Println("Connecting to", connectionString)

	var err error
	ctx := context.Background()
	//c.DBCLient, err = mongo.NewClient(options.Client().ApplyURI(connectionString))
	//err = c.DBCLient.Connect(context.Background())
	c.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)
	} else {
		if err = c.Client.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatal(err)
		} else {
			log.Println("Successfully connected to db")
		}
	}

	// set target db and connections
	c.Database = c.Client.Database(def.Settings.Values[requiredFields["database"]])
	c.Collection = c.Database.Collection(def.Settings.Values[requiredFields["collection"]])
}

// CloseConnection ... Disconnects and deallocates resources
func (c *Connector) CloseConnection() {
	ctx := context.Background()
	c.Client.Disconnect(ctx)
}
