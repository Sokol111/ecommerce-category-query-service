package kafka

import (
	"context"
	"reflect"

	"github.com/Sokol111/ecommerce-category-service-api/events"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/config"
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideDeserializer),
		consumer.RegisterHandlerAndConsumer("category-events", newCategoryHandler),
	)
}

func provideDeserializer(lc fx.Lifecycle, kafkaConf config.Config, log *zap.Logger) (consumer.Deserializer, error) {
	// Map Avro schema full names to Go types
	typeMap := consumer.TypeMapping{
		"com.ecommerce.events.category.CategoryCreatedEvent": reflect.TypeOf(events.CategoryCreatedEvent{}),
		"com.ecommerce.events.category.CategoryUpdatedEvent": reflect.TypeOf(events.CategoryUpdatedEvent{}),
	}

	deserializer, err := consumer.NewAvroDeserializer(kafkaConf.SchemaRegistry, typeMap)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("closing schema registry deserializer")
			return deserializer.Close()
		},
	})

	return deserializer, nil
}
