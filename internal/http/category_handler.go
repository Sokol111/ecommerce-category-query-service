package http

import (
	"context"
	"errors"

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
) httpapi.StrictServerInterface {
	return &categoryHandler{
		getAllActiveCategoriesHandler: getAllActiveCategoriesHandler,
		getCategoryByIDHandler:        getCategoryByIDHandler,
	}
}

func (h *categoryHandler) GetAllActiveCategories(c context.Context, _ httpapi.GetAllActiveCategoriesRequestObject) (httpapi.GetAllActiveCategoriesResponseObject, error) {
	q := query.GetAllActiveCategoriesQuery{}

	categories, err := h.getAllActiveCategoriesHandler.Handle(c, q)
	if err != nil {
		return nil, err
	}

	response := make(httpapi.GetAllActiveCategories200JSONResponse, 0, len(categories))
	for _, cat := range categories {
		response = append(response, httpapi.CategoryResponse{
			Id:         cat.ID,
			Name:       cat.Name,
			Attributes: mapAttributes(cat.Attributes),
		})
	}
	return response, nil
}

func (h *categoryHandler) GetCategoryById(c context.Context, request httpapi.GetCategoryByIdRequestObject) (httpapi.GetCategoryByIdResponseObject, error) {
	q := query.GetCategoryByIDQuery{ID: request.Id}

	category, err := h.getCategoryByIDHandler.Handle(c, q)
	if err != nil {
		if errors.Is(err, persistence.ErrEntityNotFound) {
			detail := "Category with ID " + request.Id + " was not found"
			return httpapi.GetCategoryById404ApplicationProblemPlusJSONResponse{
				Type:   "about:blank",
				Title:  "Category not found",
				Status: 404,
				Detail: &detail,
			}, nil
		}
		return nil, err
	}

	return httpapi.GetCategoryById200JSONResponse{
		Id:         category.ID,
		Name:       category.Name,
		Attributes: mapAttributes(category.Attributes),
	}, nil
}

func mapAttributes(attrs []categoryview.CategoryAttribute) []httpapi.CategoryAttribute {
	result := make([]httpapi.CategoryAttribute, 0, len(attrs))
	for _, attr := range attrs {
		result = append(result, httpapi.CategoryAttribute{
			AttributeId: attr.AttributeID,
			Role:        httpapi.AttributeRole(attr.Role),
			Required:    attr.Required,
			SortOrder:   attr.SortOrder,
			Filterable:  attr.Filterable,
			Searchable:  attr.Searchable,
			Enabled:     attr.Enabled,
		})
	}
	return result
}
