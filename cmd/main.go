package main

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/categorylist"
	"github.com/Sokol111/ecommerce-category-query-service/internal/http"
	"github.com/Sokol111/ecommerce-commons/pkg/module"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	module.NewInfraModule(),
	module.NewKafkaModule(),
	categorylist.NewCategoryListViewModule(),
	http.NewHttpHandlerModule(),
)

func main() {
	app := fx.New(
		AppModules,
		fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Info("Application starting...")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("Application stopping...")
					return nil
				},
			})
		}),
	)
	app.Run()
}
