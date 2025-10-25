package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
)

type categoryUpdatedHandler struct {
	repo categoryview.Repository
}

func newCategoryUpdatedHandler(repo categoryview.Repository) consumer.Handler[messaging.CategoryUpdated] {
	return &categoryUpdatedHandler{
		repo: repo,
	}
}

func (h *categoryUpdatedHandler) Process(ctx context.Context, e *messaging.Event[messaging.CategoryUpdated]) error {
	view := categoryview.NewCategoryView(
		e.Payload.CategoryID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Enabled,
		e.Payload.CreatedAt,
		e.Payload.ModifiedAt,
	)

	return h.repo.Upsert(ctx, view)
}

func (h *categoryUpdatedHandler) Validate(payload *messaging.CategoryUpdated) error {
	return nil
}
