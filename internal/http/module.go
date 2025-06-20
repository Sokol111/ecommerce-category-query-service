package http

import (
	"github.com/Sokol111/ecommerce-category-query-service-api/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewHttpHandlerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			newCategoryHandler,
			func(ssi api.StrictServerInterface) api.ServerInterface {
				return api.NewStrictHandler(ssi, nil)
			},
		),
		fx.Invoke(registerRoutes),
	)
}

func registerRoutes(engine *gin.Engine, serverInterface api.ServerInterface) {
	api.RegisterHandlers(engine, serverInterface)
}
