package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type categoryUpdatedHandler struct {
	log *zap.Logger
}

func newCategoryUpdatedHandler(log *zap.Logger) consumer.Handler[payload.CategoryUpdated] {
	return &categoryUpdatedHandler{
		log: log,
	}
}

func (h *categoryUpdatedHandler) Process(ctx context.Context, e *event.Event[payload.CategoryUpdated]) error {
	h.log.Info("message received", zap.String("message", fmt.Sprintf("%v", e)))
	return nil
}

func (h *categoryUpdatedHandler) Validate(payload *payload.CategoryUpdated) error {
	return nil
}
