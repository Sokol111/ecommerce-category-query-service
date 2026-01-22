package mongo

import (
	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

type categoryViewMapper struct{}

func newCategoryViewMapper() *categoryViewMapper {
	return &categoryViewMapper{}
}

func (m *categoryViewMapper) ToEntity(domain *categoryview.CategoryView) *categoryViewEntity {
	attributes := lo.Map(domain.Attributes, func(attr categoryview.CategoryAttribute, _ int) categoryAttributeEntity {
		return categoryAttributeEntity{
			AttributeID: attr.AttributeID,
			Slug:        attr.Slug,
			Role:        attr.Role,
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		}
	})

	return &categoryViewEntity{
		ID:         domain.ID,
		Version:    domain.Version,
		Name:       domain.Name,
		Enabled:    domain.Enabled,
		Attributes: attributes,
		CreatedAt:  domain.CreatedAt,
		ModifiedAt: domain.ModifiedAt,
	}
}

func (m *categoryViewMapper) ToDomain(entity *categoryViewEntity) *categoryview.CategoryView {
	attributes := lo.Map(entity.Attributes, func(attr categoryAttributeEntity, _ int) categoryview.CategoryAttribute {
		return categoryview.CategoryAttribute{
			AttributeID: attr.AttributeID,
			Slug:        attr.Slug,
			Role:        attr.Role,
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		}
	})

	return categoryview.Reconstruct(
		entity.ID,
		entity.Version,
		entity.Name,
		entity.Enabled,
		attributes,
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
