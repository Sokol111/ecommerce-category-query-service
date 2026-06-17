package connect

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	categoryqueryv1 "github.com/Sokol111/ecommerce-category-query-service-api/gen/connect/category_query/v1"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/attributeview"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/samber/lo"
)

type categoryHandler struct {
	getAllActiveCategoriesHandler categoryview.GetAllActiveCategoriesQueryHandler
	getCategoryByIDHandler        categoryview.GetCategoryByIDQueryHandler
	attributeRepo                 attributeview.Repository
}

func (h *categoryHandler) GetCategoryById(ctx context.Context, req *connect.Request[categoryqueryv1.GetCategoryByIdRequest]) (*connect.Response[categoryqueryv1.GetCategoryByIdResponse], error) { //nolint:revive
	q := categoryview.GetCategoryByIDQuery{ID: req.Msg.GetId()}

	cat, err := h.getCategoryByIDHandler.Handle(ctx, q)
	if err != nil {
		if errors.Is(err, mongo.ErrEntityNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	attrByID, err := h.fetchMasterAttributes(ctx, []*categoryview.CategoryView{cat})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&categoryqueryv1.GetCategoryByIdResponse{
		Category: toProtoCategory(cat, attrByID),
	}), nil
}

func (h *categoryHandler) GetAllActiveCategories(ctx context.Context, _ *connect.Request[categoryqueryv1.GetAllActiveCategoriesRequest]) (*connect.Response[categoryqueryv1.GetAllActiveCategoriesResponse], error) {
	q := categoryview.GetAllActiveCategoriesQuery{}

	categories, err := h.getAllActiveCategoriesHandler.Handle(ctx, q)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	attrByID, err := h.fetchMasterAttributes(ctx, categories)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoCategories := make([]*categoryqueryv1.Category, len(categories))
	for i, cat := range categories {
		protoCategories[i] = toProtoCategory(cat, attrByID)
	}

	return connect.NewResponse(&categoryqueryv1.GetAllActiveCategoriesResponse{
		Categories: protoCategories,
	}), nil
}

// ==================== Helpers ====================

func (h *categoryHandler) fetchMasterAttributes(ctx context.Context, categories []*categoryview.CategoryView) (map[string]*attributeview.AttributeView, error) {
	allAttrIDs := make(map[string]struct{})
	for _, cat := range categories {
		for _, attr := range cat.Attributes {
			allAttrIDs[attr.AttributeID] = struct{}{}
		}
	}

	if len(allAttrIDs) == 0 {
		return make(map[string]*attributeview.AttributeView), nil
	}

	attrIDs := make([]string, 0, len(allAttrIDs))
	for id := range allAttrIDs {
		attrIDs = append(attrIDs, id)
	}

	masterAttrs, err := h.attributeRepo.FindByIDs(ctx, attrIDs)
	if err != nil {
		return nil, err
	}

	return lo.KeyBy(masterAttrs, func(a *attributeview.AttributeView) string { return a.ID }), nil
}

func toProtoCategory(cat *categoryview.CategoryView, attrByID map[string]*attributeview.AttributeView) *categoryqueryv1.Category {
	attrs := lo.FilterMap(cat.Attributes, func(attr categoryview.CategoryAttribute, _ int) (*categoryqueryv1.CategoryAttribute, bool) {
		master, ok := attrByID[attr.AttributeID]
		if !ok || !master.Enabled {
			return nil, false
		}

		opts := make([]*categoryqueryv1.AttributeOption, len(master.Options))
		for i, o := range master.Options {
			opts[i] = &categoryqueryv1.AttributeOption{
				Name:      o.Name,
				Slug:      o.Slug,
				ColorCode: o.ColorCode,
				SortOrder: int32(o.SortOrder),
			}
		}

		return &categoryqueryv1.CategoryAttribute{
			AttributeId: attr.AttributeID,
			Name:        master.Name,
			Slug:        attr.Slug,
			Type:        stringToProtoAttributeType(string(master.Type)),
			Unit:        master.Unit,
			Options:     opts,
			Role:        stringToProtoCategoryAttributeRole(attr.Role),
			SortOrder:   int32(attr.SortOrder),
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		}, true
	})

	return &categoryqueryv1.Category{
		Id:         cat.ID,
		Name:       cat.Name,
		Attributes: attrs,
	}
}

func stringToProtoAttributeType(s string) categoryqueryv1.AttributeType {
	switch s {
	case "single":
		return categoryqueryv1.AttributeType_ATTRIBUTE_TYPE_SINGLE
	case "multiple":
		return categoryqueryv1.AttributeType_ATTRIBUTE_TYPE_MULTIPLE
	case "range":
		return categoryqueryv1.AttributeType_ATTRIBUTE_TYPE_RANGE
	case "boolean":
		return categoryqueryv1.AttributeType_ATTRIBUTE_TYPE_BOOLEAN
	case "text":
		return categoryqueryv1.AttributeType_ATTRIBUTE_TYPE_TEXT
	default:
		return categoryqueryv1.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED
	}
}

func stringToProtoCategoryAttributeRole(s string) categoryqueryv1.CategoryAttributeRole {
	switch s {
	case "variant":
		return categoryqueryv1.CategoryAttributeRole_CATEGORY_ATTRIBUTE_ROLE_VARIANT
	case "specification":
		return categoryqueryv1.CategoryAttributeRole_CATEGORY_ATTRIBUTE_ROLE_SPECIFICATION
	default:
		return categoryqueryv1.CategoryAttributeRole_CATEGORY_ATTRIBUTE_ROLE_UNSPECIFIED
	}
}
