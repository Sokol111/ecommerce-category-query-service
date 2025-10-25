package categoryview

import "time"

// CategoryView - read model for category queries (CQRS query side)
// Unlike the domain Category in the command service, this is a denormalized view optimized for reads
type CategoryView struct {
	ID         string
	Version    int
	Name       string
	Enabled    bool
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// Reconstruct creates a CategoryView from persistence data
func Reconstruct(id string, version int, name string, enabled bool, createdAt, modifiedAt time.Time) *CategoryView {
	return &CategoryView{
		ID:         id,
		Version:    version,
		Name:       name,
		Enabled:    enabled,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
	}
}

// NewCategoryView creates a new category view from event data
func NewCategoryView(id string, version int, name string, enabled bool, createdAt, modifiedAt time.Time) *CategoryView {
	return &CategoryView{
		ID:         id,
		Version:    version,
		Name:       name,
		Enabled:    enabled,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
	}
}
