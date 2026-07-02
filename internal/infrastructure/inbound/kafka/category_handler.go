package kafka

import (
	"context"

	catalog_eventsv1 "github.com/Sokol111/ecommerce-catalog-service-api/gen/events/catalog/v1"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/samber/lo"
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

func (h *categoryHandler) HandleCategoryUpdated(ctx context.Context, e *catalog_eventsv1.CategoryUpdatedEvent) error {
	attributes := mapCategoryAttributes(e.Attributes)

	view := categoryview.Reconstruct(
		e.CategoryId,
		e.Version,
		e.Name,
		e.Enabled,
		attributes,
		e.CreatedAt.AsTime().UTC(),
		e.ModifiedAt.AsTime().UTC(),
	)

	cmd := categoryview.UpsertCategoryCommand{Category: view}
	if err := h.upsertHandler.Handle(ctx, cmd); err != nil {
		return err
	}

	h.log(ctx).Debug("category view updated",
		zap.String("categoryID", e.CategoryId),
		zap.Int64("version", e.Version))

	return nil
}

func mapCategoryAttributes(eventAttrs []*catalog_eventsv1.CategoryAttribute) []categoryview.CategoryAttribute {
	return lo.Map(eventAttrs, func(attr *catalog_eventsv1.CategoryAttribute, _ int) categoryview.CategoryAttribute {
		return categoryview.CategoryAttribute{
			AttributeID: attr.AttributeId,
			Slug:        attr.AttributeSlug,
			Role:        attr.Role.String(),
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		}
	})
}

func (h *categoryHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "category-handler"))
}
