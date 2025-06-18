package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/model"
	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type categoryUpdatedHandler struct {
	service model.CategoryListService
	log     *zap.Logger
}

func newCategoryUpdatedHandler(service model.CategoryListService, log *zap.Logger) consumer.Handler[payload.CategoryUpdated] {
	return &categoryUpdatedHandler{
		service: service,
		log:     log,
	}
}

func (h *categoryUpdatedHandler) Process(ctx context.Context, e *event.Event[payload.CategoryUpdated]) error {
	return h.service.ProcessCategoryUpdatedEvent(ctx, e)
}

func (h *categoryUpdatedHandler) Validate(payload *payload.CategoryUpdated) error {
	return nil
}
