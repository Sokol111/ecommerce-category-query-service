package kafka

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/model"
	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type categoryCreatedHandler struct {
	service model.CategoryListService
	log     *zap.Logger
}

func newCategoryCreatedHandler(service model.CategoryListService, log *zap.Logger) consumer.Handler[payload.CategoryCreated] {
	return &categoryCreatedHandler{
		service: service,
		log:     log,
	}
}

func (h *categoryCreatedHandler) Process(ctx context.Context, e *event.Event[payload.CategoryCreated]) error {
	return h.service.ProcessCategoryCreatedEvent(ctx, e)
}

func (h *categoryCreatedHandler) Validate(payload *payload.CategoryCreated) error {
	return nil
}
