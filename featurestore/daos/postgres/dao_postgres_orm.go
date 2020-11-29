package postgres

import "time"

type FeatureSetPostgresDao struct {
	ID          int64
	InsertedAt  time.Time
	Version     string
	Features    []Feature
	Description string
	Labels      map[string]string
}

// Version ... definition of version for a feature set
type VersionPostgresDao struct{}

// Feature ... a named variable with a data type
type FeaturePostgresDao struct {
	ID       int64
	Name     string
	Value    interface{}
	DataType string
}
