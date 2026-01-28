package catalog

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockProductsRepository struct {
	getAllProductsFunc   func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error)
	getProductByCodeFunc func(ctx context.Context, code string) (*models.Product, error)
}

func (m *mockProductsRepository) GetAllProducts(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
	if m.getAllProductsFunc != nil {
		return m.getAllProductsFunc(ctx, offset, limit, categoryCode, priceLessThan)
	}
	return nil, 0, nil
}

func (m *mockProductsRepository) GetProductByCode(ctx context.Context, code string) (*models.Product, error) {
	if m.getProductByCodeFunc != nil {
		return m.getProductByCodeFunc(ctx, code)
	}
	return nil, nil
}

func TestCatalogHandler_HandleGet(t *testing.T) {
	t.Run("successful catalog retrieval with default pagination", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				assert.Equal(t, 0, offset)
				assert.Equal(t, 10, limit)
				assert.Equal(t, "", categoryCode)
				assert.Nil(t, priceLessThan)

				return []models.Product{
					{
						Code:  "PROD001",
						Price: decimal.NewFromFloat(10.99),
						Category: models.Category{
							Code: "clothing",
							Name: "Clothing",
						},
					},
				}, 1, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "PROD001")
		assert.Contains(t, recorder.Body.String(), "Clothing")
		assert.Contains(t, recorder.Body.String(), `"total":1`)
	})

	t.Run("catalog retrieval with custom pagination", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				assert.Equal(t, 5, offset)
				assert.Equal(t, 20, limit)
				return []models.Product{}, 0, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog?offset=5&limit=20", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("catalog retrieval with limit validation", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				assert.Equal(t, 100, limit)
				return []models.Product{}, 0, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog?limit=150", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("catalog retrieval with minimum limit validation", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				assert.Equal(t, 1, limit)
				return []models.Product{}, 0, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog?limit=0", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("catalog retrieval with category filter", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				assert.Equal(t, "shoes", categoryCode)
				return []models.Product{
					{
						Code:  "PROD002",
						Price: decimal.NewFromFloat(12.49),
						Category: models.Category{
							Code: "shoes",
							Name: "Shoes",
						},
					},
				}, 1, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog?category=shoes", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "PROD002")
		assert.Contains(t, recorder.Body.String(), "Shoes")
	})

	t.Run("catalog retrieval with price filter", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				assert.NotNil(t, priceLessThan)
				expected := decimal.NewFromFloat(15.00)
				assert.True(t, priceLessThan.Equal(expected))
				return []models.Product{
					{
						Code:  "PROD001",
						Price: decimal.NewFromFloat(10.99),
						Category: models.Category{
							Code: "clothing",
							Name: "Clothing",
						},
					},
				}, 1, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog?priceLessThan=15.00", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("catalog retrieval with all filters combined", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				assert.Equal(t, 2, offset)
				assert.Equal(t, 5, limit)
				assert.Equal(t, "clothing", categoryCode)
				assert.NotNil(t, priceLessThan)
				return []models.Product{}, 0, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog?offset=2&limit=5&category=clothing&priceLessThan=20.00", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("repository error handling", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getAllProductsFunc: func(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
				return nil, 0, errors.New("database error")
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "database error")
	})
}

func TestCatalogHandler_HandleGetByCode(t *testing.T) {
	t.Run("successful product retrieval", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
				assert.Equal(t, "PROD001", code)
				return &models.Product{
					Code:  "PROD001",
					Price: decimal.NewFromFloat(10.99),
					Category: models.Category{
						Code: "clothing",
						Name: "Clothing",
					},
					Variants: []models.Variant{
						{
							Name:  "Variant A",
							SKU:   "SKU001A",
							Price: decimal.NewFromFloat(11.99),
						},
						{
							Name:  "Variant B",
							SKU:   "SKU001B",
							Price: decimal.NewFromFloat(0),
						},
					},
				}, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "PROD001")
		assert.Contains(t, recorder.Body.String(), "Clothing")
		assert.Contains(t, recorder.Body.String(), "Variant A")
		assert.Contains(t, recorder.Body.String(), "SKU001A")
		assert.Contains(t, recorder.Body.String(), "11.99")
		assert.Contains(t, recorder.Body.String(), "10.99")
	})

	t.Run("product not found", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
				return nil, errors.New("record not found")
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog/INVALID", nil)
		req.SetPathValue("code", "INVALID")
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Product not found")
	})

	t.Run("empty product code", func(t *testing.T) {
		mockRepo := &mockProductsRepository{}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog/", nil)
		req.SetPathValue("code", "")
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Product code is required")
	})

	t.Run("product with no variants", func(t *testing.T) {
		mockRepo := &mockProductsRepository{
			getProductByCodeFunc: func(ctx context.Context, code string) (*models.Product, error) {
				return &models.Product{
					Code:  "PROD006",
					Price: decimal.NewFromFloat(5.50),
					Category: models.Category{
						Code: "shoes",
						Name: "Shoes",
					},
					Variants: []models.Variant{},
				}, nil
			},
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest("GET", "/catalog/PROD006", nil)
		req.SetPathValue("code", "PROD006")
		recorder := httptest.NewRecorder()

		handler.HandleGetByCode(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "PROD006")
		assert.Contains(t, recorder.Body.String(), `"variants":[]`)
	})
}
