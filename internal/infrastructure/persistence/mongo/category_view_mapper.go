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
	attributes := make([]categoryAttributeEntity, 0, len(domain.Attributes))
	for _, attr := range domain.Attributes {
		attributes = append(attributes, categoryAttributeEntity{
			AttributeID: attr.AttributeID,
			Role:        string(attr.Role),
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
			Enabled:     attr.Enabled,
		})
	}

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
	attributes := make([]categoryview.CategoryAttribute, 0, len(entity.Attributes))
	for _, attr := range entity.Attributes {
		attributes = append(attributes, categoryview.CategoryAttribute{
			AttributeID: attr.AttributeID,
			Role:        categoryview.AttributeRole(attr.Role),
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
			Enabled:     attr.Enabled,
		})
	}

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
