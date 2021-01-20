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

## Disclaimer

Mastro is still on development and largely untested. Please fork the repo and extend it at wish.

## TL-DR

Terminology:
* FeatureStore - service to manage features (i.e., featureSets and featureStates);
* Catalogue - service to manage data assets (i.e., static data definitions and their relationships);
* Crawler - any agent able to list and walk a file system, filter and parse asset definitions (i.e. manifest files) and push them to the catalogue;

Help:
* [PlantUML Diagram of the repo](https://www.dumels.com/diagram/2e5f820a-1822-4852-8259-4811deefa789)

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

```go
// FeatureState ... a versioned set of features refered to a window over a reference time series or stream
type FeatureState struct {
	Description string            `json:"description,omitempty"`
	Features    []Feature         `json:"features,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// FeatureSet ... a versioned set of features
type FeatureSet struct {
	InsertedAt  time.Time         `json:"inserted_at,omitempty"`
	Version     string            `json:"version,omitempty"`
	Features    []Feature         `json:"features,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Version ... definition of version for a feature set
type Version struct{}

// Feature ... a named variable with a data type
type Feature struct {
	Name     string      `json:"name,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	DataType string      `json:"data-type,omitempty"`
}
```

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

## Configuration

The package `conf` defines the structure of the Yaml configuration, to be provided as input.
The config can be used to start one of the three different types: i) crawler, ii) catalogue or iii) featurestore.
This is defined using the `ConfigType`, an alias for those cases.
Additional `Details` are also provided as a map to start the component.
Each component is defined by a `DataSourceDefinition` defining the connection details to a backend persistence service.

```go
// Config ... Defines a model for the input config files
type Config struct {
	ConfigType           ConfigType           `yaml:"type"`
	Details              map[string]string    `yaml:"details,omitempty"`
	DataSourceDefinition DataSourceDefinition `yaml:"backend"`
}

// ConfigType ... config type
type ConfigType string

const (
	// Crawler ... crawler agent config type
	Crawler ConfigType = "crawler"
	// Catalogue ... catalogue config type
	Catalogue = "catalogue"
	// FeatureStore ... featurestore config type
	FeatureStore = "featurestore"
)
```

The `DataSourceDefinition` is defined as a user-selected `name` and a `type`.

```go
// DataSourceDefinition ... connection details for a data source connector
type DataSourceDefinition struct {
	Name              string            `yaml:"name"`
	Type              string            `yaml:"type"`
	CrawlerDefinition CrawlerDefinition `yaml:"crawler,omitempty"`
	Settings          map[string]string `yaml:"settings,omitempty"`
	// optional kerberos section
	KerberosDetails *KerberosDetails `yaml:"kerberos"`
	// optional tls section
	TLSDetails *TLSDetails `yaml:"tls"`
}
```

A `CrawlerDefinition` is optionally provided to the `crawler` component to determine scraping information.

```go
// CrawlerDefinition ... Config for a Crawler service
type CrawlerDefinition struct {
	RootFolder     string `yaml:"root-folder"`
	FilterFilename string `yaml:"filter-filename"`
	ScheduleEvery  Period `yaml:"schedule-period"`
	ScheduleValue  uint64 `yaml:"schedule-value"`
}
```

### Feature store

An example configuration for a feature store is defined below:

```yaml
type: featurestore
details:
  port: 8085
backend:
  name: test-mongo
  type: mongo
  settings:
    username: mongo
    password: test
    host: "localhost:27017"
    schema: features
```

### Catalogue

An example configuration for a mongo-based catalogue service is defined below:

```yaml
type: catalogue
details:
  port: 8085
backend:
  name: test-mongo
  type: mongo
  settings:
    username: mongo
    password: test
    host: "localhost:27017"
    schema: catalogue
```

### Crawler

An example configuration for an S3 crawler is defined below:

```yaml
type: crawler
details:
  endpoint: localhost
  port: 8085
backend:
  name: public-minio-s3
  type: s3
  crawler:
    root-folder: ""
    filter-filename: "MANIFEST.yaml"
  settings:
    endpoint: "play.min.io"
    access-key-id: "Q3AM3UQ867SPQQA43P2F"
    secret-access-key: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
    use-ssl: "true"
```