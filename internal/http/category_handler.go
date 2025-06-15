package http

import (
	"context"
	"net/http"

	"github.com/Sokol111/ecommerce-category-query-service-api/api"
	"github.com/Sokol111/ecommerce-category-query-service/pkg/model"
)

type categoryHandler struct {
	categoryListService model.CategoryListService
}

func newCategoryHandler(service model.CategoryListService) api.StrictServerInterface {
	return &categoryHandler{service}
}

func (h *categoryHandler) GetAllActiveCategories(c context.Context, _ api.GetAllActiveCategoriesRequestObject) (api.GetAllActiveCategoriesResponseObject, error) {
	dto, err := h.categoryListService.GetActiveCategories(c)
	if err != nil {
		return api.GetAllActiveCategories500JSONResponse{Code: 500, Message: http.StatusText(500)}, err
	}
	response := make(api.GetAllActiveCategories200JSONResponse, 0, len(dto.Categories))
	for _, u := range dto.Categories {
		response = append(response, api.CategoryResponse{
			Id:   u.ID,
			Name: u.Name,
		})
	}
	return response, nil
}
