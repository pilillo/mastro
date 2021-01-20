package abstract

import (
	"time"

	"gopkg.in/yaml.v2"
)

// Asset ... managed resource
type Asset struct {
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

// Item ... Managed item
type Item struct {
	LastDiscoveredAt time.Time
	Asset            Asset
}

// ParseAsset ... Parse an asset specification file
func ParseAsset(data []byte) (*Asset, error) {
	asset := Asset{}

	err := yaml.Unmarshal(data, &asset)

	return &asset, err
}

// ValidateAsset ... Validate asset specification file
func ValidateAsset(asset Asset) (*Asset, error) {
	return &Asset{}, nil
}
