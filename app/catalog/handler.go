package catalog

import (
	"context"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type ProductsRepository interface {
	GetAllProducts(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error)
	GetProductByCode(ctx context.Context, code string) (*models.Product, error)
}

type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type ProductDetailResponse struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category string           `json:"category"`
	Variants []ProductVariant `json:"variants"`
}

type ProductVariant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type CatalogHandler struct {
	repo ProductsRepository
}

func NewCatalogHandler(r ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 10

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			if val < 1 {
				limit = 1
			} else if val > 100 {
				limit = 100
			} else {
				limit = val
			}
		}
	}

	categoryCode := r.URL.Query().Get("category")
	var priceLessThan *decimal.Decimal
	if priceStr := r.URL.Query().Get("priceLessThan"); priceStr != "" {
		if price, err := decimal.NewFromString(priceStr); err == nil {
			priceLessThan = &price
		}
	}

	res, total, err := h.repo.GetAllProducts(r.Context(), offset, limit, categoryCode, priceLessThan)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	response := Response{
		Products: products,
		Total:    total,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	product, err := h.repo.GetProductByCode(r.Context(), code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	variants := make([]ProductVariant, len(product.Variants))
	for i, v := range product.Variants {
		price := v.Price
		if price.IsZero() {
			price = product.Price
		}
		variants[i] = ProductVariant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.InexactFloat64(),
		}
	}

	response := ProductDetailResponse{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: product.Category.Name,
		Variants: variants,
	}

	api.OKResponse(w, response)
}
