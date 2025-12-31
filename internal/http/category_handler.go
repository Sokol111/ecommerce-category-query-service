package http

import (
	"context"
	"errors"
	"net/url"

	"github.com/samber/lo"

	"github.com/Sokol111/ecommerce-category-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/query"
	"github.com/Sokol111/ecommerce-category-query-service/internal/domain/categoryview"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type categoryHandler struct {
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler
	getCategoryByIDHandler        query.GetCategoryByIDQueryHandler
}

func newCategoryHandler(
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler,
	getCategoryByIDHandler query.GetCategoryByIDQueryHandler,
) httpapi.Handler {
	return &categoryHandler{
		getAllActiveCategoriesHandler: getAllActiveCategoriesHandler,
		getCategoryByIDHandler:        getCategoryByIDHandler,
	}
}

var aboutBlankURL, _ = url.Parse("about:blank")

func mapOption(opt categoryview.AttributeOption, _ int) httpapi.AttributeOption {
	return httpapi.AttributeOption{
		Name:      opt.Name,
		Slug:      opt.Slug,
		ColorCode: httpapi.NewOptString(lo.FromPtr(opt.ColorCode)),
		SortOrder: opt.SortOrder,
		Enabled:   opt.Enabled,
	}
}

func mapAttribute(attr categoryview.CategoryAttribute, _ int) httpapi.CategoryAttribute {
	return httpapi.CategoryAttribute{
		AttributeId: attr.AttributeID,
		Name:        attr.Name,
		Slug:        attr.Slug,
		Type:        httpapi.AttributeType(attr.Type),
		Unit:        httpapi.NewOptString(lo.FromPtr(attr.Unit)),
		Options:     lo.Map(attr.Options, mapOption),
		Role:        httpapi.AttributeRole(attr.Role),
		Required:    attr.Required,
		SortOrder:   attr.SortOrder,
		Filterable:  attr.Filterable,
		Searchable:  attr.Searchable,
		Enabled:     attr.Enabled,
	}
}

func toCategoryResponse(cat *categoryview.CategoryView) httpapi.CategoryResponse {
	return httpapi.CategoryResponse{
		ID:         cat.ID,
		Name:       cat.Name,
		Attributes: lo.Map(cat.Attributes, mapAttribute),
	}
}

func (h *categoryHandler) GetAllActiveCategories(ctx context.Context) (httpapi.GetAllActiveCategoriesRes, error) {
	q := query.GetAllActiveCategoriesQuery{}

	categories, err := h.getAllActiveCategoriesHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	response := lo.Map(categories, func(cat *categoryview.CategoryView, _ int) httpapi.CategoryResponse {
		return toCategoryResponse(cat)
	})

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

	response := toCategoryResponse(category)
	return &response, nil
}
