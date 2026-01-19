package kafka

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	catalogapi "github.com/Sokol111/ecommerce-catalog-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-category-query-service/internal/infrastructure/client"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/zap"
)

type categoryHandler struct {
	repo       categoryview.Repository
	attrClient client.AttributeClient
}

func newCategoryHandler(repo categoryview.Repository, attrClient client.AttributeClient) *categoryHandler {
	return &categoryHandler{
		repo:       repo,
		attrClient: attrClient,
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
	attributes, err := h.enrichAndMapAttributes(ctx, e.Payload.Attributes)
	if err != nil {
		return fmt.Errorf("failed to enrich attributes: %w", err)
	}

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
	attributes, err := h.enrichAndMapAttributes(ctx, e.Payload.Attributes)
	if err != nil {
		return fmt.Errorf("failed to enrich attributes: %w", err)
	}

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

// enrichAndMapAttributes fetches attribute data from catalog-service and builds domain attributes
func (h *categoryHandler) enrichAndMapAttributes(ctx context.Context, eventAttrs []events.CategoryAttribute) ([]categoryview.CategoryAttribute, error) {
	if len(eventAttrs) == 0 {
		return []categoryview.CategoryAttribute{}, nil
	}

	// Collect attribute IDs
	ids := make([]string, len(eventAttrs))
	for i, attr := range eventAttrs {
		ids[i] = attr.AttributeID
	}

	// Fetch attribute data from catalog-service
	attrDataMap, err := h.attrClient.GetAttributesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	attributes := make([]categoryview.CategoryAttribute, 0, len(eventAttrs))
	for _, attr := range eventAttrs {
		attrData := attrDataMap[attr.AttributeID]

		var unit *string
		if attrData.Unit.IsSet() {
			unit = &attrData.Unit.Value
		}

		options := mapAttributeOptions(attrData.Options)

		attributes = append(attributes, categoryview.CategoryAttribute{
			AttributeID: attr.AttributeID,
			Name:        attrData.Name,
			Slug:        attrData.Slug,
			Type:        string(attrData.Type),
			Unit:        unit,
			Options:     options,
			Role:        string(attr.Role),
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		})
	}

	return attributes, nil
}

func mapAttributeOptions(options []catalogapi.AttributeOption) []categoryview.AttributeOption {
	if options == nil {
		return []categoryview.AttributeOption{}
	}
	result := make([]categoryview.AttributeOption, 0, len(options))
	for _, opt := range options {
		var colorCode *string
		if opt.ColorCode.IsSet() {
			colorCode = &opt.ColorCode.Value
		}

		result = append(result, categoryview.AttributeOption{
			Name:      opt.Name,
			Slug:      opt.Slug,
			ColorCode: colorCode,
			SortOrder: opt.SortOrder,
			Enabled:   opt.Enabled,
		})
	}
	return result
}

func (h *categoryHandler) log(ctx context.Context) *zap.Logger {
	return logger.Get(ctx).With(zap.String("component", "category-handler"))
}
