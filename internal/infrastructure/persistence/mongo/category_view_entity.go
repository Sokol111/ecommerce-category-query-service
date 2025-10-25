package mongo

import (
	"time"
)

// categoryViewEntity represents the MongoDB document structure for category views
type categoryViewEntity struct {
	ID         string    `bson:"_id"`
	Version    int       `bson:"version"`
	Name       string    `bson:"name"`
	Enabled    bool      `bson:"enabled"`
	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}
