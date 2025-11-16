package kafka

import (
	"github.com/Sokol111/ecommerce-category-service-api/events"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		consumer.RegisterTypeMapping(events.DefaultTypeMapping),
		consumer.RegisterHandlerAndConsumer("category-events", newCategoryHandler),
	)
}
