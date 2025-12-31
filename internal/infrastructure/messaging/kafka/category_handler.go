package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-category-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
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
	switch evt := event.(type) {
	case *events.CategoryCreatedEvent:
		return h.handleCategoryCreated(ctx, evt)
	case *events.CategoryUpdatedEvent:
		return h.handleCategoryUpdated(ctx, evt)
	default:
		return fmt.Errorf("unhandled event type: %T: %w", event, consumer.ErrSkipMessage)
	}
}

func (h *categoryHandler) handleCategoryCreated(ctx context.Context, e *events.CategoryCreatedEvent) error {
	attributes := mapEventAttributes(e.Payload.Attributes)

	view := categoryview.NewCategoryView(
		e.Payload.CategoryID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Enabled,
		attributes,
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
	attributes := mapEventAttributes(e.Payload.Attributes)

	view := categoryview.NewCategoryView(
		e.Payload.CategoryID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Enabled,
		attributes,
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

func mapEventAttributes(eventAttrs []events.CategoryAttribute) []categoryview.CategoryAttribute {
	attributes := make([]categoryview.CategoryAttribute, 0, len(eventAttrs))
	for _, attr := range eventAttrs {
		options := make([]categoryview.AttributeOption, 0, len(attr.Options))
		for _, opt := range attr.Options {
			options = append(options, categoryview.AttributeOption{
				Name:      opt.Name,
				Slug:      opt.Slug,
				ColorCode: opt.ColorCode,
				SortOrder: opt.SortOrder,
				Enabled:   opt.Enabled,
			})
		}
		attributes = append(attributes, categoryview.CategoryAttribute{
			AttributeID: attr.AttributeID,
			Name:        attr.Name,
			Slug:        attr.Slug,
			Type:        attr.Type,
			Unit:        attr.Unit,
			Options:     options,
			Role:        attr.Role,
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
			Enabled:     attr.Enabled,
		})
	}
	return attributes
}

func (h *categoryHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "category-handler"))
}
