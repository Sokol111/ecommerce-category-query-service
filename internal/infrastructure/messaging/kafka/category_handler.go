package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-category-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	commonsevents "github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/events"
	"go.uber.org/zap"
)

type categoryHandler struct {
	repo categoryview.Repository
}

func newCategoryHandler(repo categoryview.Repository) *categoryHandler {
	return &categoryHandler{
		repo: repo,
	}
}

func (h *categoryHandler) Process(ctx context.Context, event any) error {
	e, ok := event.(commonsevents.Event)
	if !ok {
		return fmt.Errorf("event does not implement Event interface: %T: %w", event, consumer.ErrSkipMessage)
	}

	// Now switch on concrete types - exhaustive linter will warn if any Event type is missing
	switch evt := e.(type) {
	case *events.CategoryCreatedEvent:
		return h.handleCategoryCreated(ctx, evt)
	case *events.CategoryUpdatedEvent:
		return h.handleCategoryUpdated(ctx, evt)
	default:
		// If exhaustive linter is enabled and all Event types are handled above,
		// this case should theoretically never be reached
		return fmt.Errorf("unhandled event type: %T: %w", event, consumer.ErrSkipMessage)
	}
}

func (h *categoryHandler) handleCategoryCreated(ctx context.Context, e *events.CategoryCreatedEvent) error {
	view := categoryview.NewCategoryView(
		e.Payload.CategoryID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Enabled,
		e.Payload.CreatedAt,
		e.Payload.ModifiedAt,
	)

	if err := h.repo.Upsert(ctx, view); err != nil {
		return fmt.Errorf("failed to upsert category view: %w", err)
	}

	h.log(ctx).Debug("category view created",
		zap.String("categoryID", e.Payload.CategoryID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

func (h *categoryHandler) handleCategoryUpdated(ctx context.Context, e *events.CategoryUpdatedEvent) error {
	view := categoryview.NewCategoryView(
		e.Payload.CategoryID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Enabled,
		e.Payload.CreatedAt,
		e.Payload.ModifiedAt,
	)

	if err := h.repo.Upsert(ctx, view); err != nil {
		return fmt.Errorf("failed to upsert category view: %w", err)
	}

	h.log(ctx).Debug("category view updated",
		zap.String("categoryID", e.Payload.CategoryID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

func (h *categoryHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "category-handler"))
}
