```
(\ 
\'\ 
 \'\     __________  		___  ___          _             
 / '|   ()_________)		|  \/  |         | |            
 \ '/    \ ~~~~~~~~ \		| .  . | __ _ ___| |_ _ __ ___  
   \       \ ~~~~~~   \		| |\/| |/ _  / __| __| __/  _  \ 
   ==).      \__________\	| |  | | (_| \__ \ |_| | | (_) |
  (__)       ()__________)	\_|  |_/\__,_|___/\__|_|  \___/ 
```
---
Data and Feature Catalogue in Go


## Connectors

A Data source is defined in the `abstract` package as follows:

```go
type ConnectorProvider interface {
	ValidateDataSourceDefinition(*conf.DataSourceDefinition) error
	InitConnection(*conf.DataSourceDefinition)
	CloseConnection()
}
```

Have a look at the `sources/*` packages for specific implementations of the interface.

A factory is generally used to instantiate the connector with default settings:

```go
func NewElasticConnector() *Connector {
	return &Connector{}
}
```
The connector can be then used for any of the implemented DAOs to be started.

## Feature Store

A feature store is a service to store and version features.

### FeatureSets and FeatureStates

A Feature can either be computed on a dataset or a data stream, respectively using a batch or a stream processing pipeline.
This is due to the different life cycle and performance requirements for collecting and serving those data to end applications.

A data access object (DAO) for a featureSet is defined as follows:

```go
type FeatureSetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(fs *FeatureSet) error
	GetById(id string) (*FeatureSet, error)
	ListAllFeatureSets() (*[]FeatureSet, error)
	CloseConnection()
}
```

The interface is then implemented for specific targets in the `featurestore/daos/*` packages.

Each DAO also implements a lazy singleton using `sync.once` (see [blog post](https://medium.com/@ishagirdhar/singleton-pattern-in-golang-9f60d7fdab23)).
This way, all DAO implementations can be efficiently linked from a `dao_mappings.go` file, for instance:

```go
var availableDAOs = map[string]func() abstract.FeatureSetDAOProvider{
	"mongo":   mongo.GetSingleton,
	"elastic": elastic.GetSingleton,
}
```

As for the exposed service, the `featurestore/service.go` defines a basic interface to retrieve featureSets:

```go
type Service interface {
	Init(cfg *conf.Config) *errors.RestErr
	CreateFeatureSet(fs abstract.FeatureSet) (*abstract.FeatureSet, *errors.RestErr)
	GetFeatureSetByID(fsID string) (*abstract.FeatureSet, *errors.RestErr)
	ListAllFeatureSets() (*[]abstract.FeatureSet, *errors.RestErr)
}
```

This is for instance how to add a new featureSet calculated in the test environment of a fictional project.

*PUT* on `localhost:8085/featureset` with body:
```json
{
	"version" : "test-v1.0",
	"description" : "example feature set for testing purposes",
	"labels" : {
	    "refers-to" : "project-gilberto",
	    "environment" : "test"
	},
	"features" : [
		{
			"name":"feature1",
			"value":"10",
			"data-type":"int"
		},
		{
			"name":"feature2",
			"value":"true",
			"data-type":"bool"
		}
	]
}
```

with the service adding a date time for additional versioning and finally replying with:
```json
{
    "inserted_at": "2020-11-29T17:24:01.747543Z",
    "version": "test-v1.0",
    "features": [
        {
            "name": "feature1",
            "value": "10",
            "data-type": "int"
        },
        {
            "name": "feature2",
            "value": "true",
            "data-type": "bool"
        }
    ],
    "description": "example feature set for testing purposes",
    "labels": {
        "environment": "test",
        "refers-to": "project-gilberto"
    }
}
```

### Example
Have a look at the conf folder for an example configuration, using either ElasticSearch or Mongo as backend.

```
./mastro --configpath conf/featurestore/elastic/example_elastic.cfg
```


## Data Catalogue
Data providers can describe and publish data using a shared definition format.
Consequently, data definitions can be crawled from networked and distributed file systems, as well as directly published to a common endpoint.

### Catalogue API
A Catalogue servuce endpoint implements the following interface:

```go
type Service interface {
	Init(cfg *conf.Config) *errors.RestErr
	CreateAsset(fs abstract.Asset) (*abstract.Asset, *errors.RestErr)
	GetAssetByID(assetID string) (*abstract.Asset, *errors.RestErr)
	ListAllAssets() (*[]abstract.Asset, *errors.RestErr)
}
```

This can be easily mapped to a DAO backend:
```go
type AssetDAOProvider interface {
	Init(*conf.DataSourceDefinition)
	Create(asset *Asset) error
	GetById(id string) (*Asset, error)
	ListAllAssets() (*[]Asset, error)
	CloseConnection()
}
```

Have a look at `catalogue/daos/*` for example implementations.

### Crawlers
A crawler is an agent traversing file systems to seek for asset definition files.
Crawlers implement the Crawler interface:

```go
type Crawler interface {
	InitConnection(cfg *conf.Config) (Crawler, error)
	WalkWithFilter(root string, filenameFilter string) ([]Asset, error)
}
```

Specifically, the crawler inits the connection to a volume (e.g., hdfs, s3) whereas in the WalkWithFilter it traverses the file system starting from the provided root path.
A filter is provided to only select specific metadata files, whose naming follows a reserved global setting such as `MANIFEST.yml`. Selected files are then marshalled and returned using the `abstract.Asset` definition:

```go
// Asset ... managed resource
type Asset struct {
	// asset publication datetime
	PublishedOn timeutils.Time `yaml:"published-on" json:"published-on"`
	// name of the asset
	Name string `yaml:"name" json:"name"`
	// description of the asset
	Description string `yaml:"description" json:"description"`
	// the list of assets this depends on
	DependsOn []string `yaml:"depends-on" json:"depends-on"`
	// the actual asset metadata
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	// tags are flags used to simplify asset search
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// Metadata ... Asset metadata
type Metadata struct {
	// asset type - refers to the type.name
	TypeName string `yaml:"type" json:"type"`
	// asset attributes, key-val list
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

// Type ... Asset type
type Type struct {
	Name string `yaml:"name" json:"name"`
	ID   string `yaml:"id" json:"id"`
}

// Item ... Managed item
type Item struct {
	LastDiscoveredAt timeutils.Time
	Asset            Asset
}
```

The package also provide means to parse and validate assets:
```go
func ParseAsset(data []byte) (*Asset, error) {}
func ValidateAsset(asset Asset) (*Asset, error) {}
```

## Disclaimer

Mastro is still on development and largely untested.