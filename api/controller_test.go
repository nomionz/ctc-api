package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomionz/ctc-api/models"
	"github.com/nomionz/ctc-api/repositories"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockStore struct {
	prods []models.Product
}

var _ repositories.Repository = &mockStore{}

func (m *mockStore) List() ([]models.Product, error) {
	return m.prods, nil
}

func (m *mockStore) Get(id int) (*models.Product, error) {
	for _, p := range m.prods {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockStore) Create(prod *models.Product) error {
	m.prods = append(m.prods, *prod)
	return nil
}

func (m *mockStore) Update(prod *models.Product) error {
	for i, p := range m.prods {
		if p.ID == prod.ID {
			m.prods[i] = *prod
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockStore) Delete(id int) error {
	for i, p := range m.prods {
		if p.ID == id {
			m.prods = append(m.prods[:i], m.prods[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func NewMockStore() *mockStore {
	return &mockStore{
		prods: make([]models.Product, 0),
	}
}

var (
	validProd = models.Product{
		Name:   "test",
		Price:  decimal.NewFromFloat(1.0),
		Amount: 1,
	}
	invalidProd = models.Product{
		Name:   "a",
		Price:  decimal.NewFromFloat(1.0),
		Amount: -1,
	}
	validJson, _   = json.Marshal(validProd)
	invalidJson, _ = json.Marshal(invalidProd)
)

func sendRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreate(t *testing.T) {
	store := NewMockStore()
	controller := NewController(store)
	r := controller.router
	w := sendRequest(r, "POST", "/products/", string(validJson))

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, 1, len(store.prods))
	assert.Contains(t, store.prods, validProd)

	w = sendRequest(r, "POST", "/products/", string(invalidJson))
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, 1, len(store.prods))
	assert.NotContains(t, store.prods, invalidProd)
}

func TestList(t *testing.T) {
	store := NewMockStore()
	controller := NewController(store)
	r := controller.router
	w := sendRequest(r, "GET", "/products/", "")

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 0, len(store.prods))

	store.prods = append(store.prods, validProd)
	w = sendRequest(r, "GET", "/products/", "")
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 1, len(store.prods))
}

func TestGet(t *testing.T) {
	store := NewMockStore()
	controller := NewController(store)
	r := controller.router
	w := sendRequest(r, "GET", "/products/1", "")

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, 0, len(store.prods))

	copyValid := validProd
	copyValid.ID = 1
	store.prods = append(store.prods, copyValid)
	w = sendRequest(r, "GET", "/products/1", "")
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 1, len(store.prods))
}

func TestUpdate(t *testing.T) {
	store := NewMockStore()
	controller := NewController(store)
	r := controller.router
	w := sendRequest(r, "PATCH", "/products/1", string(validJson))

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, 0, len(store.prods))

	copyValid := validProd
	copyValid.ID = 1
	store.prods = append(store.prods, copyValid)
	w = sendRequest(r, "PATCH", "/products/1", string(validJson))
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 1, len(store.prods))
}

func TestDelete(t *testing.T) {
	store := NewMockStore()
	controller := NewController(store)
	r := controller.router
	w := sendRequest(r, "DELETE", "/products/1", "")

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, 0, len(store.prods))

	copyValid := validProd
	copyValid.ID = 1
	store.prods = append(store.prods, copyValid)
	w = sendRequest(r, "DELETE", "/products/1", "")
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 0, len(store.prods))
	w = sendRequest(r, "DELETE", "/products/1", "")
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, 0, len(store.prods))
	w = sendRequest(r, "DELETE", "/products/1", "")
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, 0, len(store.prods))
}
