package categorylist

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/pkg/model"
	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
)

type service struct {
	store Store
}

func newService(store Store) model.CategoryListService {
	return &service{store: store}
}

func (s *service) ProcessCategoryCreatedEvent(ctx context.Context, e *event.Event[payload.CategoryCreated]) error {
	return s.store.Upsert(ctx, e.Payload.CategoryID, e.Payload.Name, e.Payload.Version, e.Payload.Enabled)
}

func (s *service) ProcessCategoryUpdatedEvent(ctx context.Context, e *event.Event[payload.CategoryUpdated]) error {
	return s.store.Upsert(ctx, e.Payload.CategoryID, e.Payload.Name, e.Payload.Version, e.Payload.Enabled)
}

func (s *service) GetActiveCategories(ctx context.Context) (*model.CategoryListViewDTO, error) {
	return s.store.GetAllEnabled(ctx)
}
