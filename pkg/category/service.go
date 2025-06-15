package category

import (
	"context"
)

type Service interface {

	// can return NotFoundError
	GetById(ctx context.Context, id string) (*Category, error)

	GetAll(ctx context.Context) ([]*Category, error)
}
