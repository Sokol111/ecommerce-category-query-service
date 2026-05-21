package kafka

import (
	"github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		events.Module(),
		consumer.RegisterHandlerAndConsumer("category-events", newCategoryHandler),
		consumer.RegisterHandlerAndConsumer("attribute-events", newAttributeHandler),
	)
}
