package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type categoryCreatedHandler struct {
	log *zap.Logger
}

func newCategoryCreatedHandler(log *zap.Logger) consumer.Handler[payload.CategoryCreated] {
	return &categoryCreatedHandler{
		log: log,
	}
}

func (h *categoryCreatedHandler) Process(ctx context.Context, e *event.Event[payload.CategoryCreated]) error {
	h.log.Info("message received", zap.String("message", fmt.Sprintf("%v", e)))
	return nil
}

func (h *categoryCreatedHandler) Validate(payload *payload.CategoryCreated) error {
	return nil
}
