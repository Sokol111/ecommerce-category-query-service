package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
)

type categoryCreatedHandler struct {
	repo categoryview.Repository
}

func newCategoryCreatedHandler(repo categoryview.Repository) consumer.Handler[messaging.CategoryCreated] {
	return &categoryCreatedHandler{
		repo: repo,
	}
}

func (h *categoryCreatedHandler) Process(ctx context.Context, e *messaging.Event[messaging.CategoryCreated]) error {
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

func (h *categoryCreatedHandler) Validate(payload *messaging.CategoryCreated) error {
	return nil
}
