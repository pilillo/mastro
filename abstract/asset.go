package abstract

import (
	"time"

	"github.com/pilillo/mastro/utils/errors"
	"gopkg.in/yaml.v2"
)

// Asset ... managed resource
type Asset struct {
	// asset last found by crawler at - only added by service (not crawler/manifest itself, i.e. no yaml)
	LastDiscoveredAt time.Time `json:"last-discovered-at"`
	// asset publication datetime
	PublishedOn time.Time `yaml:"published-on" json:"published-on"`
	// name of the asset
	Name string `yaml:"name" json:"name"`
	// description of the asset
	Description string `yaml:"description" json:"description"`
	// the list of assets this depends on
	DependsOn []string `yaml:"depends-on" json:"depends-on"`
	// asset type
	Type AssetType `yaml:"type,omitempty" json:"type,omitempty"`
	// labels for the specific asset
	Labels map[string]interface{} `yaml:"labels,omitempty" json:"labels,omitempty"`
	// tags are flags used to simplify asset search
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// AssetType ... Asset type information
type AssetType string
const (
	_Database AssetType = "database"
	_Dataset = "dataset"
	_FeatureSet = "featureset"
	_Model = "model"
	_Notebook = "notebook"
	_Pipeline = "pipeline"
	_Report = "report"
	_Service = "service"
	_Stream = "stream"
	_Table = "table"
	_User = "user"
	_Workflow = "workflow"
)


// ParseAsset ... Parse an asset specification file
func ParseAsset(data []byte) (*Asset, error) {
	asset := Asset{}

	err := yaml.Unmarshal(data, &asset)

	return &asset, err
}

// Validate ... Validate asset specification file
func (asset *Asset) Validate() *errors.RestErr {

	return nil
}
