package model

import (
	"context"

	"github.com/Sokol111/ecommerce-commons/pkg/event"
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
)

type CategoryListViewDTO struct {
	Categories []CategoryDTO
}

type CategoryDTO struct {
	ID   string
	Name string
}

type CategoryListService interface {
	ProcessCategoryCreatedEvent(ctx context.Context, e *event.Event[payload.CategoryCreated]) error

	ProcessCategoryUpdatedEvent(ctx context.Context, e *event.Event[payload.CategoryUpdated]) error

	GetActiveCategories(ctx context.Context) (*CategoryListViewDTO, error)
}
