package connect

import (
	"net/http"

	"connectrpc.com/connect"
	categoryqueryv1connect "github.com/Sokol111/ecommerce-category-query-service-api/gen/connect/category_query/v1/categoryqueryv1connect"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/attributeview"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/security/validation"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			newCategoryHandler,
			provideProcedurePermissions,
		),
		fx.Invoke(registerConnectRoutes),
	)
}

func newCategoryHandler(
	getAllActiveCategoriesHandler categoryview.GetAllActiveCategoriesQueryHandler,
	getCategoryByIDHandler categoryview.GetCategoryByIDQueryHandler,
	attributeRepo attributeview.Repository,
) *categoryHandler {
	return &categoryHandler{
		getAllActiveCategoriesHandler: getAllActiveCategoriesHandler,
		getCategoryByIDHandler:        getCategoryByIDHandler,
		attributeRepo:                 attributeRepo,
	}
}

func registerConnectRoutes(
	mux *http.ServeMux,
	handler *categoryHandler,
	interceptors []connect.Interceptor,
) {
	path, h := categoryqueryv1connect.NewCategoryQueryServiceHandler(handler, connect.WithInterceptors(interceptors...))
	mux.Handle(path, h)
}

func provideProcedurePermissions() validation.ProcedurePermissions {
	return validation.ProcedurePermissions{
		categoryqueryv1connect.CategoryQueryServiceGetCategoryByIdProcedure:        nil,
		categoryqueryv1connect.CategoryQueryServiceGetAllActiveCategoriesProcedure: nil,
	}
}
