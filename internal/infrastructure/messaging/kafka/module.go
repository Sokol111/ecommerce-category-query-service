package kafka

import (
	"github.com/Sokol111/ecommerce-catalog-service-api/gen/events"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/avro/mapping"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		consumer.RegisterHandlerAndConsumer("category-events", newCategoryHandler),
		fx.Invoke(registerSchemas),
	)
}

func registerSchemas(tm *mapping.TypeMapping) error {
	return tm.RegisterBindings(events.SchemaBindings)
}
