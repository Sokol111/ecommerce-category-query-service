package application

import (
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/query"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
func Module() fx.Option {
	return fx.Options(
		// Query handlers
		fx.Provide(
			query.NewGetCategoryByIDHandler,
			query.NewGetAllActiveCategoriesHandler,
		),
	)
}
