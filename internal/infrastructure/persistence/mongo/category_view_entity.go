package mongo

import (
	"time"
)

// attributeOptionEntity represents an attribute option in MongoDB
type attributeOptionEntity struct {
	Name      string  `bson:"name"`
	Slug      string  `bson:"slug"`
	ColorCode *string `bson:"colorCode,omitempty"`
	SortOrder int     `bson:"sortOrder"`
	Enabled   bool    `bson:"enabled"`
}

// categoryAttributeEntity represents a category attribute in MongoDB
type categoryAttributeEntity struct {
	AttributeID string                  `bson:"attributeId"`
	Name        string                  `bson:"name"`
	Slug        string                  `bson:"slug"`
	Type        string                  `bson:"type"`
	Unit        *string                 `bson:"unit,omitempty"`
	Options     []attributeOptionEntity `bson:"options"`
	Role        string                  `bson:"role"`
	Required    bool                    `bson:"required"`
	SortOrder   int                     `bson:"sortOrder"`
	Filterable  bool                    `bson:"filterable"`
	Searchable  bool                    `bson:"searchable"`
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
