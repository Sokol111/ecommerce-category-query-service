package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/model"
	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
)

type categoryUpdatedHandler struct {
	categoryListService model.CategoryListService
}

func newCategoryUpdatedHandler(categoryListService model.CategoryListService) consumer.Handler[payload.CategoryUpdated] {
	return &categoryUpdatedHandler{
		categoryListService: categoryListService,
	}
}

func (h *categoryUpdatedHandler) Process(ctx context.Context, e *event.Event[payload.CategoryUpdated]) error {
	return h.categoryListService.ProcessCategoryUpdatedEvent(ctx, e)
}

func (h *categoryUpdatedHandler) Validate(payload *payload.CategoryUpdated) error {
	return nil
}
