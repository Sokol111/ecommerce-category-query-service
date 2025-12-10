package http

import (
	"context"

	"github.com/Sokol111/ecommerce-category-query-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/query"
)

type categoryHandler struct {
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler
}

func newCategoryHandler(
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler,
) httpapi.StrictServerInterface {
	return &categoryHandler{
		getAllActiveCategoriesHandler: getAllActiveCategoriesHandler,
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
			Id:   cat.ID,
			Name: cat.Name,
		})
	}
	return response, nil
}
