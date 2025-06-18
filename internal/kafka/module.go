package kafka

import (
	"github.com/Sokol111/ecommerce-commons/pkg/event/payload"
	"github.com/Sokol111/ecommerce-commons/pkg/kafka/consumer"
	"go.uber.org/fx"
)

func NewKafkaHandlerModule() fx.Option {
	return fx.Options(
		consumer.RegisterHandlerAndConsumer[payload.CategoryCreated]("categoryCreatedHandler", newCategoryCreatedHandler),
		consumer.RegisterHandlerAndConsumer[payload.CategoryUpdated]("categoryUpdatedHandler", newCategoryUpdatedHandler),
	)
}
