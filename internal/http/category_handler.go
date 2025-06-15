package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/Sokol111/ecommerce-category-query-service-api/api"
	"github.com/Sokol111/ecommerce-category-query-service/pkg/category"
)

type categoryHandler struct {
	service category.Service
}

func newCategoryHandler(service category.Service) api.StrictServerInterface {
	return &categoryHandler{service}
}

func (h *categoryHandler) GetCategoryById(c context.Context, request api.GetCategoryByIdRequestObject) (api.GetCategoryByIdResponseObject, error) {
	found, err := h.service.GetById(c, request.Id)
	if errors.Is(err, category.NotFoundError) {
		return api.GetCategoryById404JSONResponse{Code: 404, Message: "Category not found"}, nil
	}
	if err != nil {
		return api.GetCategoryById500JSONResponse{Code: 500, Message: http.StatusText(500)}, err
	}
	return api.GetCategoryById200JSONResponse{
		Id:   found.ID,
		Name: found.Name,
	}, nil
}

func (h *categoryHandler) GetAll(c context.Context, _ api.GetAllRequestObject) (api.GetAllResponseObject, error) {
	categories, err := h.service.GetAll(c)
	if err != nil {
		return api.GetAll500JSONResponse{Code: 500, Message: http.StatusText(500)}, err
	}
	response := make(api.GetAll200JSONResponse, 0, len(categories))
	for _, u := range categories {
		response = append(response, api.CategoryResponse{
			Id:   u.ID,
			Name: u.Name,
		})
	}
	return response, nil
}
