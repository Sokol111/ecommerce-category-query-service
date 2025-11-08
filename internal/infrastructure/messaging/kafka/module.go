package kafka

import (
	"reflect"

	"github.com/Sokol111/ecommerce-category-service-api/events"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/fx"
)

func Module() fx.Option {
	typeMapping := consumer.TypeMapping{
		events.EventTypeCategoryCreated: reflect.TypeOf(&events.CategoryCreatedEvent{}),
		events.EventTypeCategoryUpdated: reflect.TypeOf(&events.CategoryUpdatedEvent{}),
	}

	return fx.Provide(
		consumer.RegisterHandlerAndConsumer("category-events", newCategoryHandler),
		fx.Annotate(
			func() consumer.TypeMapping { return typeMapping },
			fx.ResultTags(`name:"category-events"`),
		),
	)
}
