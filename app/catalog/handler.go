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

type response struct {
	Products []product `json:"products"`
	Total    int64     `json:"total"`
}

type product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type productDetailResponse struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category string           `json:"category"`
	Variants []productVariant `json:"variants"`
}

type productVariant struct {
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
		val, err := strconv.Atoi(offsetStr)
		if err != nil {
			api.ErrorResponse(w, http.StatusBadRequest, "Invalid offset parameter")
			return
		}
		if val < 0 {
			api.ErrorResponse(w, http.StatusBadRequest, "Offset must be non-negative")
			return
		}
		offset = val
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil {
			api.ErrorResponse(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
		if val < 1 {
			limit = 1
		} else if val > 100 {
			limit = 100
		} else {
			limit = val
		}
	}

	categoryCode := r.URL.Query().Get("category")
	var priceLessThan *decimal.Decimal
	if priceStr := r.URL.Query().Get("priceLessThan"); priceStr != "" {
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			api.ErrorResponse(w, http.StatusBadRequest, "Invalid priceLessThan parameter")
			return
		}
		priceLessThan = &price
	}

	res, total, err := h.repo.GetAllProducts(r.Context(), offset, limit, categoryCode, priceLessThan)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	products := make([]product, len(res))
	for i, p := range res {
		products[i] = product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	resp := response{
		Products: products,
		Total:    total,
	}

	api.OKResponse(w, resp)
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	prod, err := h.repo.GetProductByCode(r.Context(), code)
	if err != nil {
		if isNotFoundError(err) {
			api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		} else {
			api.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve product")
		}
		return
	}

	variants := buildVariantsResponse(prod)

	resp := productDetailResponse{
		Code:     prod.Code,
		Price:    prod.Price.InexactFloat64(),
		Category: prod.Category.Name,
		Variants: variants,
	}

	api.OKResponse(w, resp)
}

func buildVariantsResponse(prod *models.Product) []productVariant {
	variants := make([]productVariant, len(prod.Variants))
	for i, v := range prod.Variants {
		price := v.Price
		if price.IsZero() {
			price = prod.Price
		}
		variants[i] = productVariant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.InexactFloat64(),
		}
	}
	return variants
}

func isNotFoundError(err error) bool {
	return err.Error() == "record not found"
}
