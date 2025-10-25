package query

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
)

type GetCategoryByIDQuery struct {
	ID string
}

type GetAllActiveCategoriesQuery struct{}

type GetCategoryByIDQueryHandler interface {
	Handle(ctx context.Context, query GetCategoryByIDQuery) (*categoryview.CategoryView, error)
}

type GetAllActiveCategoriesQueryHandler interface {
	Handle(ctx context.Context, query GetAllActiveCategoriesQuery) ([]*categoryview.CategoryView, error)
}
