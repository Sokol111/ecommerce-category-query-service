package kafka

import (
	"context"

	catalog_eventsv1 "github.com/Sokol111/ecommerce-catalog-service-api/gen/events/catalog/v1"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/attributeview"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type attributeHandler struct {
	upsertHandler attributeview.UpsertAttributeCommandHandler
}

func newAttributeHandler(upsertHandler attributeview.UpsertAttributeCommandHandler) *attributeHandler {
	return &attributeHandler{
		upsertHandler: upsertHandler,
	}
}

func (h *attributeHandler) HandleAttributeUpdated(ctx context.Context, e *catalog_eventsv1.AttributeUpdatedEvent) error {
	view := attributeview.Reconstruct(
		e.AttributeId,
		e.Version,
		e.Slug,
		e.Name,
		attributeview.AttributeType(e.Type),
		e.Unit,
		e.Enabled,
		e.ModifiedAt.AsTime().UTC(),
		lo.Map(e.Options, mapAttributeOption),
	)

	cmd := attributeview.UpsertAttributeCommand{Attribute: view}
	if err := h.upsertHandler.Handle(ctx, cmd); err != nil {
		return err
	}

	h.log(ctx).Debug("attribute view updated",
		zap.String("attributeID", e.AttributeId),
		zap.Int64("version", e.Version))

	return nil
}

func (h *attributeHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx)
}

func mapAttributeOption(opt *catalog_eventsv1.AttributeOption, _ int) attributeview.AttributeOption {
	return attributeview.AttributeOption{
		Slug:      opt.Slug,
		Name:      opt.Name,
		ColorCode: opt.ColorCode,
		SortOrder: opt.SortOrder,
	}
}
