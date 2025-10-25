package categoryview

import "context"

type Repository interface {
	// Upsert inserts or updates a category view (for event processing)
	Upsert(ctx context.Context, category *CategoryView) error

	// FindByID retrieves a category view by ID
	FindByID(ctx context.Context, id string) (*CategoryView, error)

	// FindAllEnabled retrieves all enabled categories
	FindAllEnabled(ctx context.Context) ([]*CategoryView, error)
}
