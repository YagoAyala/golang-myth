package categories

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoriesRepository interface {
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) error
}

type categoriesResponse struct {
	Categories []categoryItem `json:"categories"`
}

type categoryItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type createCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesHandler struct {
	repo CategoriesRepository
}

func NewCategoriesHandler(r CategoriesRepository) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

func (h *CategoriesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories(r.Context())
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	items := make([]categoryItem, len(categories))
	for i, c := range categories {
		items[i] = categoryItem{
			Code: c.Code,
			Name: c.Name,
		}
	}

	resp := categoriesResponse{
		Categories: items,
	}

	api.OKResponse(w, resp)
}

func (h *CategoriesHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Code == "" || req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Code and name are required")
		return
	}

	category := &models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.CreateCategory(r.Context(), category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := categoryItem{
		Code: category.Code,
		Name: category.Name,
	}

	api.OKResponse(w, resp)
}
