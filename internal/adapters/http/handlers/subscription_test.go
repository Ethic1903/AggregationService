package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"AggregationService/internal/domain/models/dto"
	custom_err "AggregationService/internal/errors"
)

// Мок usecase с правильной сигнатурой CalculateCost
type mockUseCase struct{ mock.Mock }

func (m *mockUseCase) Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.SubscriptionResponse), args.Error(1)
}
func (m *mockUseCase) GetByID(ctx context.Context, id int) (*dto.SubscriptionResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.SubscriptionResponse), args.Error(1)
}
func (m *mockUseCase) Update(ctx context.Context, id int, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*dto.SubscriptionResponse), args.Error(1)
}
func (m *mockUseCase) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockUseCase) GetAll(ctx context.Context, userID *uuid.UUID, service *string, limit, offset int) ([]*dto.SubscriptionResponse, error) {
	args := m.Called(ctx, userID, service, limit, offset)
	return args.Get(0).([]*dto.SubscriptionResponse), args.Error(1)
}
func (m *mockUseCase) CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate *time.Time) (int, error) {
	args := m.Called(ctx, userID, serviceName, startDate, endDate)
	return args.Int(0), args.Error(1)
}

// Конструктор хэндлера
func newTestHandler(useCase *mockUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{useCase: useCase}
}

func TestSubscriptionHandler_Create(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	reqBody := dto.CreateSubscriptionRequest{
		UserID:      uuid.New(),
		ServiceName: "yandex",
		Price:       299,
		StartDate:   "09-2025",
	}
	sub := &dto.SubscriptionResponse{ID: 1, ServiceName: "yandex"}
	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.CreateSubscriptionRequest")).Return(sub, nil)

	body, _ := json.Marshal(reqBody)
	r := chi.NewRouter()
	r.Post("/subscriptions", handler.Create)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp dto.SubscriptionResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, sub.ID, resp.ID)
}

func TestSubscriptionHandler_Create_Invalid(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.CreateSubscriptionRequest")).Return((*dto.SubscriptionResponse)(nil), custom_err.ErrInvalidRequest)

	reqBody := dto.CreateSubscriptionRequest{}
	body, _ := json.Marshal(reqBody)
	r := chi.NewRouter()
	r.Post("/subscriptions", handler.Create)

	req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSubscriptionHandler_GetByID(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	sub := &dto.SubscriptionResponse{ID: 1, ServiceName: "yandex"}
	mockUC.On("GetByID", mock.Anything, 1).Return(sub, nil)

	r := chi.NewRouter()
	r.Get("/subscriptions/{id}", handler.GetByID)

	req := httptest.NewRequest("GET", "/subscriptions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.SubscriptionResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, sub.ID, resp.ID)
}

func TestSubscriptionHandler_GetByID_NotFound(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	mockUC.On("GetByID", mock.Anything, 999).Return((*dto.SubscriptionResponse)(nil), custom_err.ErrSubscriptionNotFound)

	r := chi.NewRouter()
	r.Get("/subscriptions/{id}", handler.GetByID)

	req := httptest.NewRequest("GET", "/subscriptions/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubscriptionHandler_Update(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	serviceName := "yandex plus"
	price := 399
	reqBody := dto.UpdateSubscriptionRequest{
		ServiceName: &serviceName,
		Price:       &price,
	}
	sub := &dto.SubscriptionResponse{ID: 1, ServiceName: "yandex plus"}
	mockUC.On("Update", mock.Anything, 1, mock.AnythingOfType("*dto.UpdateSubscriptionRequest")).Return(sub, nil)

	body, _ := json.Marshal(reqBody)
	r := chi.NewRouter()
	r.Put("/subscriptions/{id}", handler.Update)

	req := httptest.NewRequest("PUT", "/subscriptions/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.SubscriptionResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, sub.ID, resp.ID)
}

func TestSubscriptionHandler_Update_NotFound(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	mockUC.On("Update", mock.Anything, 999, mock.AnythingOfType("*dto.UpdateSubscriptionRequest")).Return((*dto.SubscriptionResponse)(nil), custom_err.ErrSubscriptionNotFound)

	reqBody := dto.UpdateSubscriptionRequest{}
	body, _ := json.Marshal(reqBody)
	r := chi.NewRouter()
	r.Put("/subscriptions/{id}", handler.Update)

	req := httptest.NewRequest("PUT", "/subscriptions/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubscriptionHandler_Delete(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	mockUC.On("Delete", mock.Anything, 1).Return(nil)

	r := chi.NewRouter()
	r.Delete("/subscriptions/{id}", handler.Delete)

	req := httptest.NewRequest("DELETE", "/subscriptions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestSubscriptionHandler_Delete_NotFound(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	mockUC.On("Delete", mock.Anything, 999).Return(custom_err.ErrSubscriptionNotFound)

	r := chi.NewRouter()
	r.Delete("/subscriptions/{id}", handler.Delete)

	req := httptest.NewRequest("DELETE", "/subscriptions/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubscriptionHandler_GetAll(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	validUUID := uuid.New()
	serviceName := "yandex"
	subs := []*dto.SubscriptionResponse{
		{ID: 1, ServiceName: "yandex"},
		{ID: 2, ServiceName: "yandex plus"},
	}
	mockUC.On("GetAll", mock.Anything, &validUUID, &serviceName, 10, 0).Return(subs, nil)

	r := chi.NewRouter()
	r.Get("/subscriptions", handler.GetAll)

	req := httptest.NewRequest("GET", "/subscriptions?user_id="+validUUID.String()+"&service_name="+serviceName+"&limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.SubscriptionResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
}

func TestSubscriptionHandler_CalculateCost(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	validUUID := uuid.New()
	serviceName := "yandex"
	startDate := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
	mockUC.On("CalculateCost", mock.Anything, &validUUID, &serviceName, &startDate, &endDate).Return(451, nil)

	r := chi.NewRouter()
	r.Get("/subscriptions/cost", handler.CalculateCost)

	req := httptest.NewRequest("GET", "/subscriptions/cost?user_id="+validUUID.String()+"&service_name="+serviceName+"&start_date=09-2025&end_date=12-2025", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Cost int `json:"cost"`
	}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 451, resp.Cost)
}

func TestSubscriptionHandler_CalculateCost_Error(t *testing.T) {
	mockUC := new(mockUseCase)
	handler := newTestHandler(mockUC)

	validUUID := uuid.New()
	serviceName := "yandex"
	startDate := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
	mockUC.On("CalculateCost", mock.Anything, &validUUID, &serviceName, &startDate, &endDate).Return(0, custom_err.ErrInternalServer)

	r := chi.NewRouter()
	r.Get("/subscriptions/cost", handler.CalculateCost)

	req := httptest.NewRequest("GET", "/subscriptions/cost?user_id="+validUUID.String()+"&service_name="+serviceName+"&start_date=09-2025&end_date=12-2025", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
