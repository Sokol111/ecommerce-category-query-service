package categoryview

import "time"

// AttributeOption represents an option for single/multiple type attributes
type AttributeOption struct {
	Name      string
	Slug      string
	ColorCode *string
	SortOrder int
}

// CategoryAttribute represents an attribute assigned to a category
type CategoryAttribute struct {
	AttributeID string
	Name        string
	Slug        string
	Type        string
	Unit        *string
	Options     []AttributeOption
	Role        string
	Required    bool
	SortOrder   int
	Filterable  bool
	Searchable  bool
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

// NewCategoryView creates a new category view from event data
func NewCategoryView(id string, version int, name string, enabled bool, attributes []CategoryAttribute, createdAt, modifiedAt time.Time) *CategoryView {
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
