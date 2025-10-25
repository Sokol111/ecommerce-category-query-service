package http

import (
	"context"
	"net/http"

	"github.com/Sokol111/ecommerce-category-query-service-api/api"
	"github.com/Sokol111/ecommerce-category-query-service/internal/application/query"
)

type categoryHandler struct {
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler
}

func newCategoryHandler(
	getAllActiveCategoriesHandler query.GetAllActiveCategoriesQueryHandler,
) api.StrictServerInterface {
	return &categoryHandler{
		getAllActiveCategoriesHandler: getAllActiveCategoriesHandler,
	}
}

func (h *categoryHandler) GetAllActiveCategories(c context.Context, _ api.GetAllActiveCategoriesRequestObject) (api.GetAllActiveCategoriesResponseObject, error) {
	q := query.GetAllActiveCategoriesQuery{}

	categories, err := h.getAllActiveCategoriesHandler.Handle(c, q)
	if err != nil {
		return api.GetAllActiveCategories500JSONResponse{Code: 500, Message: http.StatusText(500)}, err
	}

	response := make(api.GetAllActiveCategories200JSONResponse, 0, len(categories))
	for _, cat := range categories {
		response = append(response, api.CategoryResponse{
			Id:   cat.ID,
			Name: cat.Name,
		})
	}
	return response, nil
}
