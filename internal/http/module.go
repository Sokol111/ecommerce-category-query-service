package http //nolint:revive // intentional package name to group HTTP handlers

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Sokol111/ecommerce-category-query-service-api/gen/httpapi"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			newCategoryHandler,
			httpapi.ProvideServer,
		),
		fx.Invoke(registerOgenRoutes),
	)
}

func registerOgenRoutes(mux *http.ServeMux, server *httpapi.Server) {
	mux.Handle("/", server)
}
