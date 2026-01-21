package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
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
	attributes := mapAttributes(e.Payload.Attributes)

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
	attributes := mapAttributes(e.Payload.Attributes)

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

// mapAttributes converts event attributes to domain attributes using data from message
func mapAttributes(eventAttrs []events.CategoryAttribute) []categoryview.CategoryAttribute {
	if len(eventAttrs) == 0 {
		return []categoryview.CategoryAttribute{}
	}

	attributes := make([]categoryview.CategoryAttribute, len(eventAttrs))
	for i, attr := range eventAttrs {
		attributes[i] = categoryview.CategoryAttribute{
			AttributeID: attr.AttributeID,
			Name:        attr.AttributeName,
			Slug:        attr.AttributeSlug,
			Type:        attr.AttributeType,
			Unit:        attr.AttributeUnit,
			Options:     mapAttributeOptions(attr.AttributeOptions),
			Role:        attr.Role,
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		}
	}

	return attributes
}

func mapAttributeOptions(options []events.AttributeOption) []categoryview.AttributeOption {
	if options == nil {
		return []categoryview.AttributeOption{}
	}
	result := make([]categoryview.AttributeOption, len(options))
	for i, opt := range options {
		result[i] = categoryview.AttributeOption{
			Name:      opt.Name,
			Slug:      opt.Slug,
			ColorCode: opt.ColorCode,
			SortOrder: opt.SortOrder,
		}
	}
	return result
}

func (h *categoryHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "category-handler"))
}
