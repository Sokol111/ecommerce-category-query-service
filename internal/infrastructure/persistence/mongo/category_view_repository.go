package mongo

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

type categoryViewRepository struct {
	*commonsmongo.GenericRepository[categoryview.CategoryView, categoryViewEntity]
	mapper *categoryViewMapper
}

func newCategoryViewRepository(mongo commonsmongo.Mongo, mapper *categoryViewMapper) (categoryview.Repository, error) {
	genericRepo, err := commonsmongo.NewGenericRepository(mongo.GetCollection("category_view"), mapper)
	if err != nil {
		return nil, err
	}

	return &categoryViewRepository{
		GenericRepository: genericRepo,
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
