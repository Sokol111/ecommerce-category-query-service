package query

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type GetCategoryByIDQuery struct {
	ID string
}

type GetCategoryByIDQueryHandler interface {
	Handle(ctx context.Context, query GetCategoryByIDQuery) (*categoryview.CategoryView, error)
}

type getCategoryByIDHandler struct {
	repo categoryview.Repository
}

func NewGetCategoryByIDHandler(repo categoryview.Repository) GetCategoryByIDQueryHandler {
	return &getCategoryByIDHandler{repo: repo}
}

func (h *getCategoryByIDHandler) Handle(ctx context.Context, query GetCategoryByIDQuery) (*categoryview.CategoryView, error) {
	c, err := h.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, persistence.ErrEntityNotFound) {
			return nil, persistence.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to get category view: %w", err)
	}
	return c, nil
}
