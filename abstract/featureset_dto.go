package abstract

import (
	"time"

	"github.com/pilillo/mastro/utils/errors"
)

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

// Feature ... a named variable with a data type
type Feature struct {
	Name     string      `json:"name,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	DataType string      `json:"data-type,omitempty"`
}

// Validate ... validate a featureSet
func (fs *FeatureSet) Validate() *errors.RestErr {
	return nil
}

// Validate ... validate a feature
func (fs *Feature) Validate() *errors.RestErr {
	// todo: validate data type
	return nil
}
