package categoryview

import (
	"context"
	"fmt"
)

type UpsertCategoryCommand struct {
	Category *CategoryView
}

type UpsertCategoryCommandHandler interface {
	Handle(ctx context.Context, cmd UpsertCategoryCommand) error
}

type upsertCategoryHandler struct {
	repo Repository
}

func NewUpsertCategoryHandler(repo Repository) UpsertCategoryCommandHandler {
	return &upsertCategoryHandler{repo: repo}
}

func (h *upsertCategoryHandler) Handle(ctx context.Context, cmd UpsertCategoryCommand) error {
	if err := h.repo.Upsert(ctx, cmd.Category); err != nil {
		return fmt.Errorf("failed to upsert category view: %w", err)
	}
	return nil
}
