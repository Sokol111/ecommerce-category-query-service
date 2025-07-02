package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/model"
	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
)

type categoryCreatedHandler struct {
	categoryListService model.CategoryListService
}

func newCategoryCreatedHandler(categoryListService model.CategoryListService) consumer.Handler[payload.CategoryCreated] {
	return &categoryCreatedHandler{
		categoryListService: categoryListService,
	}
}

func (h *categoryCreatedHandler) Process(ctx context.Context, e *event.Event[payload.CategoryCreated]) error {
	return h.categoryListService.ProcessCategoryCreatedEvent(ctx, e)
}

func (h *categoryCreatedHandler) Validate(payload *payload.CategoryCreated) error {
	return nil
}
