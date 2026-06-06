package kafka

import (
	"context"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"go.uber.org/zap"
)

type categoryHandler struct {
	upsertHandler categoryview.UpsertCategoryCommandHandler
}

func newCategoryHandler(upsertHandler categoryview.UpsertCategoryCommandHandler) *categoryHandler {
	return &categoryHandler{
		upsertHandler: upsertHandler,
	}
}

func (h *categoryHandler) HandleCategoryUpdated(ctx context.Context, e *events.CategoryUpdatedEvent) error {
	attributes := mapAttributes(e.Payload.Attributes)

	view := categoryview.Reconstruct(
		e.Payload.CategoryID,
		e.Payload.Version,
		e.Payload.Name,
		e.Payload.Enabled,
		attributes,
		e.Payload.CreatedAt,
		e.Payload.ModifiedAt,
	)

	cmd := categoryview.UpsertCategoryCommand{Category: view}
	if err := h.upsertHandler.Handle(ctx, cmd); err != nil {
		return err
	}

	h.log(ctx).Debug("category view updated",
		zap.String("categoryID", e.Payload.CategoryID),
		zap.String("eventID", e.Metadata.EventID),
		zap.Int("version", e.Payload.Version))

	return nil
}

func mapAttributes(eventAttrs []events.CategoryAttribute) []categoryview.CategoryAttribute {
	return lo.Map(eventAttrs, func(attr events.CategoryAttribute, _ int) categoryview.CategoryAttribute {
		return categoryview.CategoryAttribute{
			AttributeID: attr.AttributeID,
			Slug:        attr.AttributeSlug,
			Role:        attr.Role,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		}
	})
}

func (h *categoryHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "category-handler"))
}
