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
