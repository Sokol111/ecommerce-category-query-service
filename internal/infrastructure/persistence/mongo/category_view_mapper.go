package mongo

import (
	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

type categoryViewMapper struct{}

func newCategoryViewMapper() *categoryViewMapper {
	return &categoryViewMapper{}
}

func (m *categoryViewMapper) ToEntity(domain *categoryview.CategoryView) *categoryViewEntity {
	return &categoryViewEntity{
		ID:         domain.ID,
		Version:    domain.Version,
		Name:       domain.Name,
		Enabled:    domain.Enabled,
		CreatedAt:  domain.CreatedAt,
		ModifiedAt: domain.ModifiedAt,
	}
}

func (m *categoryViewMapper) ToDomain(entity *categoryViewEntity) *categoryview.CategoryView {
	return categoryview.Reconstruct(
		entity.ID,
		entity.Version,
		entity.Name,
		entity.Enabled,
		entity.CreatedAt,
		entity.ModifiedAt,
	)
}

func (m *categoryViewMapper) GetID(entity *categoryViewEntity) string {
	return entity.ID
}

func (m *categoryViewMapper) GetVersion(entity *categoryViewEntity) int {
	return entity.Version
}

func (m *categoryViewMapper) SetVersion(entity *categoryViewEntity, version int) {
	entity.Version = version
}

// Ensure mapper implements the interface
var _ commonsmongo.EntityMapper[categoryview.CategoryView, categoryViewEntity] = (*categoryViewMapper)(nil)
