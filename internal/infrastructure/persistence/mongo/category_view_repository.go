package mongo

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type categoryViewRepository struct {
	*commonsmongo.GenericRepository[categoryview.CategoryView, categoryViewEntity]
	collection commonsmongo.Collection
	mapper     *categoryViewMapper
}

func newCategoryViewRepository(mongo commonsmongo.Mongo, mapper *categoryViewMapper) (categoryview.Repository, error) {
	collection := mongo.GetCollection("category_view")
	genericRepo, err := commonsmongo.NewGenericRepository(collection, mapper)
	if err != nil {
		return nil, err
	}

	return &categoryViewRepository{
		GenericRepository: genericRepo,
		collection:        collection,
		mapper:            mapper,
	}, nil
}

func (r *categoryViewRepository) Upsert(ctx context.Context, category *categoryview.CategoryView) error {
	updated, err := r.UpsertIfNewer(ctx, category)
	if err != nil {
		return err
	}

	if !updated {
		logger.Get(ctx).Debug("version conflict during upsert", zap.String("id", category.ID))
	}

	return nil
}

func (r *categoryViewRepository) FindAllEnabled(ctx context.Context) ([]*categoryview.CategoryView, error) {
	return r.FindAllWithFilter(ctx, bson.D{{Key: "enabled", Value: true}}, nil)
}
