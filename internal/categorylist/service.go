package categorylist

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
)

type Service interface {
	ProcessCategoryCreatedEvent(ctx context.Context, e *event.Event[payload.CategoryCreated]) error

	ProcessCategoryUpdatedEvent(ctx context.Context, e *event.Event[payload.CategoryUpdated]) error

	GetActiveCategories(ctx context.Context) (*CategoryListViewDTO, error)
}

type service struct {
	store Store
}

func newService(store Store) Service {
	return &service{store: store}
}

func (s *service) ProcessCategoryCreatedEvent(ctx context.Context, e *event.Event[payload.CategoryCreated]) error {
	return s.store.Upsert(ctx, e.Payload.CategoryID, e.Payload.Name, e.Payload.Version, e.Payload.Enabled)
}

func (s *service) ProcessCategoryUpdatedEvent(ctx context.Context, e *event.Event[payload.CategoryUpdated]) error {
	return s.store.Upsert(ctx, e.Payload.CategoryID, e.Payload.Name, e.Payload.Version, e.Payload.Enabled)
}

func (s *service) GetActiveCategories(ctx context.Context) (*CategoryListViewDTO, error) {
	return s.store.GetAllEnabled(ctx)
}
