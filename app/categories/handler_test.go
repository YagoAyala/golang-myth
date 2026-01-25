package categories

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type mockCategoriesRepository struct {
	getAllCategoriesFunc func() ([]models.Category, error)
	createCategoryFunc   func(category *models.Category) error
}

func (m *mockCategoriesRepository) GetAllCategories() ([]models.Category, error) {
	if m.getAllCategoriesFunc != nil {
		return m.getAllCategoriesFunc()
	}
	return nil, nil
}

func (m *mockCategoriesRepository) CreateCategory(category *models.Category) error {
	if m.createCategoryFunc != nil {
		return m.createCategoryFunc(category)
	}
	return nil
}

func TestCategoriesHandler_HandleGet(t *testing.T) {
	t.Run("successful categories retrieval", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{
			getAllCategoriesFunc: func() ([]models.Category, error) {
				return []models.Category{
					{
						ID:   1,
						Code: "clothing",
						Name: "Clothing",
					},
					{
						ID:   2,
						Code: "shoes",
						Name: "Shoes",
					},
					{
						ID:   3,
						Code: "accessories",
						Name: "Accessories",
					},
				}, nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		req := httptest.NewRequest("GET", "/categories", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "clothing")
		assert.Contains(t, recorder.Body.String(), "Clothing")
		assert.Contains(t, recorder.Body.String(), "shoes")
		assert.Contains(t, recorder.Body.String(), "Shoes")
		assert.Contains(t, recorder.Body.String(), "accessories")
		assert.Contains(t, recorder.Body.String(), "Accessories")
	})

	t.Run("empty categories list", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{
			getAllCategoriesFunc: func() ([]models.Category, error) {
				return []models.Category{}, nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		req := httptest.NewRequest("GET", "/categories", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"categories":[]`)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{
			getAllCategoriesFunc: func() ([]models.Category, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewCategoriesHandler(mockRepo)
		req := httptest.NewRequest("GET", "/categories", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "database error")
	})
}

func TestCategoriesHandler_HandlePost(t *testing.T) {
	t.Run("successful category creation", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{
			createCategoryFunc: func(category *models.Category) error {
				assert.Equal(t, "electronics", category.Code)
				assert.Equal(t, "Electronics", category.Name)
				category.ID = 4
				return nil
			},
		}

		handler := NewCategoriesHandler(mockRepo)

		reqBody := CreateCategoryRequest{
			Code: "electronics",
			Name: "Electronics",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "electronics")
		assert.Contains(t, recorder.Body.String(), "Electronics")
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Invalid request body")
	})

	t.Run("missing code field", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)

		reqBody := CreateCategoryRequest{
			Code: "",
			Name: "Electronics",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Code and name are required")
	})

	t.Run("missing name field", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{}
		handler := NewCategoriesHandler(mockRepo)

		reqBody := CreateCategoryRequest{
			Code: "electronics",
			Name: "",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Code and name are required")
	})

	t.Run("repository error on create", func(t *testing.T) {
		mockRepo := &mockCategoriesRepository{
			createCategoryFunc: func(category *models.Category) error {
				return errors.New("duplicate key violation")
			},
		}

		handler := NewCategoriesHandler(mockRepo)

		reqBody := CreateCategoryRequest{
			Code: "clothing",
			Name: "Clothing",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.HandlePost(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "duplicate key violation")
	})
}
