package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type categoryViewRepository struct {
	collection commonsmongo.Collection
	mapper     *categoryViewMapper
}

func newCategoryViewRepository(mongo commonsmongo.Mongo, mapper *categoryViewMapper) categoryview.Repository {
	return &categoryViewRepository{
		collection: mongo.GetCollection("category_list"),
		mapper:     mapper,
	}
}

func (r *categoryViewRepository) Upsert(ctx context.Context, category *categoryview.CategoryView) error {
	entity := r.mapper.ToEntity(category)

	filter := bson.M{
		"_id":     entity.ID,
		"version": bson.M{"$lt": entity.Version},
	}

	update := bson.M{
		"$set": bson.M{
			"name":       entity.Name,
			"enabled":    entity.Enabled,
			"version":    entity.Version,
			"createdAt":  entity.CreatedAt,
			"modifiedAt": entity.ModifiedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert category view: %w", err)
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		logger.FromContext(ctx).Debug("version conflict during upsert", zap.String("id", category.ID))
	}

	return nil
}

func (r *categoryViewRepository) FindByID(ctx context.Context, id string) (*categoryview.CategoryView, error) {
	var entity categoryViewEntity

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, persistence.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to find category view by id: %w", err)
	}

	return r.mapper.ToDomain(&entity), nil
}

func (r *categoryViewRepository) FindAllEnabled(ctx context.Context) ([]*categoryview.CategoryView, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"enabled": true})
	if err != nil {
		return nil, fmt.Errorf("failed to find enabled categories: %w", err)
	}
	defer cursor.Close(ctx)

	var entities []categoryViewEntity
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, fmt.Errorf("failed to decode categories: %w", err)
	}

	views := make([]*categoryview.CategoryView, 0, len(entities))
	for i := range entities {
		views = append(views, r.mapper.ToDomain(&entities[i]))
	}

	return views, nil
}
