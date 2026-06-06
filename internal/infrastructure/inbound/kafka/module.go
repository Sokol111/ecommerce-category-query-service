package kafka

import (
	"github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Options(
		events.Module(),
		fx.Provide(newCategoryHandler, newAttributeHandler),
		consumer.RegisterHandlerAndConsumer("category-events", newCategoryRouter),
		consumer.RegisterHandlerAndConsumer("attribute-events", newAttributeRouter),
	)
}

func newCategoryRouter(h *categoryHandler, log *zap.Logger) consumer.Handler {
	r := consumer.NewRouter(log)
	consumer.Register(r, h.HandleCategoryUpdated)
	return r
}

func newAttributeRouter(h *attributeHandler, log *zap.Logger) consumer.Handler {
	r := consumer.NewRouter(log)
	consumer.Register(r, h.HandleAttributeUpdated)
	return r
}
