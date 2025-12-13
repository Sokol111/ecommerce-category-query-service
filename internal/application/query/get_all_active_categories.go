package query

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
)

type GetAllActiveCategoriesQuery struct{}

type GetAllActiveCategoriesQueryHandler interface {
	Handle(ctx context.Context, query GetAllActiveCategoriesQuery) ([]*categoryview.CategoryView, error)
}

type getAllActiveCategoriesHandler struct {
	repo categoryview.Repository
}

func NewGetAllActiveCategoriesHandler(repo categoryview.Repository) GetAllActiveCategoriesQueryHandler {
	return &getAllActiveCategoriesHandler{repo: repo}
}

func (h *getAllActiveCategoriesHandler) Handle(ctx context.Context, query GetAllActiveCategoriesQuery) ([]*categoryview.CategoryView, error) {
	categories, err := h.repo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active categories: %w", err)
	}
	return categories, nil
}
