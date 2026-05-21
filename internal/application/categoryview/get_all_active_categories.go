package categoryview

import (
	"context"
	"fmt"
)

type GetAllActiveCategoriesQuery struct{}

type GetAllActiveCategoriesQueryHandler interface {
	Handle(ctx context.Context, query GetAllActiveCategoriesQuery) ([]*CategoryView, error)
}

type getAllActiveCategoriesHandler struct {
	repo Repository
}

func NewGetAllActiveCategoriesHandler(repo Repository) GetAllActiveCategoriesQueryHandler {
	return &getAllActiveCategoriesHandler{repo: repo}
}

func (h *getAllActiveCategoriesHandler) Handle(ctx context.Context, query GetAllActiveCategoriesQuery) ([]*CategoryView, error) {
	categories, err := h.repo.FindAllEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active categories: %w", err)
	}
	return categories, nil
}
