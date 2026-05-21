package application

import (
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/attributeview"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/categoryview"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
func Module() fx.Option {
	return fx.Options(
		// Query handlers
		fx.Provide(
			categoryview.NewGetCategoryByIDHandler,
			categoryview.NewGetAllActiveCategoriesHandler,
		),
		// Command handlers
		fx.Provide(
			categoryview.NewUpsertCategoryHandler,
			attributeview.NewUpsertAttributeHandler,
		),
	)
}
