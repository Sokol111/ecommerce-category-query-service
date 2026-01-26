package categoryview

import "time"

// CategoryAttribute represents an attribute assignment to a category
// Contains only immutable references and category-specific settings
// Mutable data (name, type, unit, options) should be joined from attribute master data
type CategoryAttribute struct {
	AttributeID string // Reference to attribute definition (UUID)
	Slug        string // Attribute URL-friendly identifier (immutable)
	Role        string // variant or specification
	Required    bool   // Whether required for products in this category
	SortOrder   int    // Sort order for display in this category
	Filterable  bool   // Whether filterable for this category
	Searchable  bool   // Whether searchable for this category
}

// CategoryView - read model for category queries (CQRS query side)
// Unlike the domain Category in the command service, this is a denormalized view optimized for reads
type CategoryView struct {
	ID         string
	Version    int
	Name       string
	Enabled    bool
	Attributes []CategoryAttribute
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// Reconstruct creates a CategoryView from persistence data
func Reconstruct(id string, version int, name string, enabled bool, attributes []CategoryAttribute, createdAt, modifiedAt time.Time) *CategoryView {
	return &CategoryView{
		ID:         id,
		Version:    version,
		Name:       name,
		Enabled:    enabled,
		Attributes: attributes,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
	}
}
