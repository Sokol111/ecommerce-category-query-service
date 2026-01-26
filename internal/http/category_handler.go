package http

import (
	"context"
	"errors"
	"net/url"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-category-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/query"
	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/attributeview"
	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type categoryHandler struct {
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler
	getCategoryByIDHandler        query.GetCategoryByIDQueryHandler
	attributeRepo                 attributeview.Repository
}

func newCategoryHandler(
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler,
	getCategoryByIDHandler query.GetCategoryByIDQueryHandler,
	attributeRepo attributeview.Repository,
) httpapi.Handler {
	return &categoryHandler{
		getAllActiveCategoriesHandler: getAllActiveCategoriesHandler,
		getCategoryByIDHandler:        getCategoryByIDHandler,
		attributeRepo:                 attributeRepo,
	}
}

var aboutBlankURL, _ = url.Parse("about:blank")

func mapOption(opt attributeview.AttributeOption, _ int) httpapi.AttributeOption {
	return httpapi.AttributeOption{
		Name:      opt.Name,
		Slug:      opt.Slug,
		ColorCode: httpapi.NewOptString(lo.FromPtr(opt.ColorCode)),
		SortOrder: opt.SortOrder,
	}
}

// fetchMasterAttributes collects all unique attribute IDs from categories and fetches them in a single batch query
func (h *categoryHandler) fetchMasterAttributes(ctx context.Context, categories []*categoryview.CategoryView) (map[string]*attributeview.AttributeView, error) {
	// Collect all unique attribute IDs from all categories
	allAttrIDs := make(map[string]struct{})
	for _, cat := range categories {
		for _, attr := range cat.Attributes {
			allAttrIDs[attr.AttributeID] = struct{}{}
		}
	}

	if len(allAttrIDs) == 0 {
		return make(map[string]*attributeview.AttributeView), nil
	}

	// Convert to slice
	attrIDs := make([]string, 0, len(allAttrIDs))
	for id := range allAttrIDs {
		attrIDs = append(attrIDs, id)
	}

	// Single batch query for all attributes
	masterAttrs, err := h.attributeRepo.FindByIDs(ctx, attrIDs)
	if err != nil {
		return nil, err
	}

	return lo.KeyBy(masterAttrs, func(attr *attributeview.AttributeView) string { return attr.ID }), nil
}

// toCategoryAttributesWithMasterData joins category attribute assignments with pre-fetched attribute master data
// Returns enriched attributes, filtering out disabled attributes
func toCategoryAttributesWithMasterData(attrs []categoryview.CategoryAttribute, attrByID map[string]*attributeview.AttributeView) []httpapi.CategoryAttribute {
	if len(attrs) == 0 {
		return []httpapi.CategoryAttribute{}
	}

	return lo.FilterMap(attrs, func(attr categoryview.CategoryAttribute, _ int) (httpapi.CategoryAttribute, bool) {
		master, ok := attrByID[attr.AttributeID]
		// Skip if master data not found or attribute is disabled
		if !ok || !master.Enabled {
			return httpapi.CategoryAttribute{}, false
		}

		return httpapi.CategoryAttribute{
			AttributeId: attr.AttributeID,
			Name:        master.Name,
			Slug:        attr.Slug,
			Type:        httpapi.AttributeType(master.Type),
			Unit:        httpapi.NewOptString(lo.FromPtr(master.Unit)),
			Options:     lo.Map(master.Options, mapOption),
			Role:        httpapi.AttributeRole(attr.Role),
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
		}, true
	})
}

// toCategoryResponseWithMasterData converts a category view to HTTP response using pre-fetched master data
func toCategoryResponseWithMasterData(cat *categoryview.CategoryView, attrByID map[string]*attributeview.AttributeView) httpapi.CategoryResponse {
	return httpapi.CategoryResponse{
		ID:         cat.ID,
		Name:       cat.Name,
		Attributes: toCategoryAttributesWithMasterData(cat.Attributes, attrByID),
	}
}

// toCategoryAttributes joins category attribute assignments with attribute master data
// Returns enriched attributes, filtering out disabled attributes
func (h *categoryHandler) toCategoryAttributes(ctx context.Context, attrs []categoryview.CategoryAttribute) ([]httpapi.CategoryAttribute, error) {
	if len(attrs) == 0 {
		return []httpapi.CategoryAttribute{}, nil
	}

	// Collect unique attribute IDs
	attrIDs := lo.Map(attrs, func(attr categoryview.CategoryAttribute, _ int) string {
		return attr.AttributeID
	})

	// Fetch attribute master data
	masterAttrs, err := h.attributeRepo.FindByIDs(ctx, attrIDs)
	if err != nil {
		return nil, err
	}

	attrByID := lo.KeyBy(masterAttrs, func(attr *attributeview.AttributeView) string { return attr.ID })

	return toCategoryAttributesWithMasterData(attrs, attrByID), nil
}

func (h *categoryHandler) toCategoryResponse(ctx context.Context, cat *categoryview.CategoryView) (httpapi.CategoryResponse, error) {
	attrs, err := h.toCategoryAttributes(ctx, cat.Attributes)
	if err != nil {
		return httpapi.CategoryResponse{}, err
	}

	return httpapi.CategoryResponse{
		ID:         cat.ID,
		Name:       cat.Name,
		Attributes: attrs,
	}, nil
}

func (h *categoryHandler) GetAllActiveCategories(ctx context.Context) (httpapi.GetAllActiveCategoriesRes, error) {
	q := query.GetAllActiveCategoriesQuery{}

	categories, err := h.getAllActiveCategoriesHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	// Fetch all attribute master data in a single batch query (N+1 fix)
	attrByID, err := h.fetchMasterAttributes(ctx, categories)
	if err != nil {
		return nil, err
	}

	// Convert to response using pre-fetched master data
	response := make([]httpapi.CategoryResponse, 0, len(categories))
	for _, cat := range categories {
		response = append(response, toCategoryResponseWithMasterData(cat, attrByID))
	}

	return (*httpapi.GetAllActiveCategoriesOKApplicationJSON)(&response), nil
}

func (h *categoryHandler) GetCategoryById(ctx context.Context, params httpapi.GetCategoryByIdParams) (httpapi.GetCategoryByIdRes, error) {
	q := query.GetCategoryByIDQuery{ID: params.ID}

	category, err := h.getCategoryByIDHandler.Handle(ctx, q)
	if err != nil {
		if errors.Is(err, persistence.ErrEntityNotFound) {
			return &httpapi.GetCategoryByIdNotFound{
				Type:   *aboutBlankURL,
				Title:  "Category not found",
				Status: 404,
				Detail: httpapi.NewOptString("Category with ID " + params.ID + " was not found"),
			}, nil
		}
		return nil, err
	}

	response, err := h.toCategoryResponse(ctx, category)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
