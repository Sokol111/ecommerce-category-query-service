package mongo

import (
	"time"
)

// categoryAttributeEntity represents a category attribute in MongoDB
type categoryAttributeEntity struct {
	AttributeID string `bson:"attributeId"`
	Role        string `bson:"role"`
	Required    bool   `bson:"required"`
	SortOrder   int    `bson:"sortOrder"`
	Filterable  bool   `bson:"filterable"`
	Searchable  bool   `bson:"searchable"`
	Enabled     bool   `bson:"enabled"`
}

// categoryViewEntity represents the MongoDB document structure for category views
type categoryViewEntity struct {
	ID         string                    `bson:"_id"`
	Version    int                       `bson:"version"`
	Name       string                    `bson:"name"`
	Enabled    bool                      `bson:"enabled"`
	Attributes []categoryAttributeEntity `bson:"attributes"`
	CreatedAt  time.Time                 `bson:"createdAt"`
	ModifiedAt time.Time                 `bson:"modifiedAt"`
}
