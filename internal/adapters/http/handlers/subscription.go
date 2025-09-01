package handlers

import (
	"AggregationService/internal/domain/models/dto"
	"AggregationService/internal/pkg/logger"
	"AggregationService/internal/pkg/utils"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ISubscriptionUseCase interface {
	Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	GetByID(ctx context.Context, id int) (*dto.SubscriptionResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, userID *uuid.UUID, service *string, limit, offset int) ([]*dto.SubscriptionResponse, error)
	CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate *time.Time) (int, error)
}

type SubscriptionHandler struct {
	useCase ISubscriptionUseCase
}

func NewSubscriptionHandler(useCase ISubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{useCase: useCase}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", slog.Any("err", err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	sub, err := h.useCase.Create(ctx, &req)
	if err != nil {
		log.Error("failed to create subscription", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Debug("success create subscription", slog.Int("id", sub.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id", slog.String("id", idStr), slog.Any("err", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := h.useCase.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to get subscription", slog.Int("id", id), slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Debug("success get subscription", slog.Int("id", id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id", slog.String("id", idStr), slog.Any("err", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", slog.Any("err", err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	sub, err := h.useCase.Update(ctx, id, &req)
	if err != nil {
		log.Error("failed to update subscription", slog.Int("id", id), slog.Any("err", err))
		if err.Error() == "subscription not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	log.Debug("success update subscription", slog.Int("id", id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id", slog.String("id", idStr), slog.Any("err", err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.useCase.Delete(ctx, id); err != nil {
		log.Error("failed to delete subscription", slog.Int("id", id), slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Debug("success delete subscription", slog.Int("id", id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var userID *uuid.UUID
	var serviceName *string
	var limit, offset int

	if v := r.URL.Query().Get("user_id"); v != "" {
		uid, err := uuid.Parse(v)
		if err == nil {
			userID = &uid
		} else {
			log.Error("invalid user_id", slog.String("user_id", v), slog.Any("err", err))
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
	}
	if v := r.URL.Query().Get("service_name"); v != "" {
		serviceName = &v
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		limit, _ = strconv.Atoi(v)
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		offset, _ = strconv.Atoi(v)
	}

	subs, err := h.useCase.GetAll(ctx, userID, serviceName, limit, offset)
	if err != nil {
		log.Error("failed to get subscriptions", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug("success get subscriptions", slog.Int("count", len(subs)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

func (h *SubscriptionHandler) CalculateCost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var userID *uuid.UUID
	var serviceName *string
	var startDate, endDate *time.Time

	if v := r.URL.Query().Get("user_id"); v != "" {
		uid, err := uuid.Parse(v)
		if err == nil {
			userID = &uid
		} else {
			log.Error("invalid user_id", slog.String("user_id", v), slog.Any("err", err))
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
	}
	if v := r.URL.Query().Get("service_name"); v != "" {
		serviceName = &v
	}
	if v := r.URL.Query().Get("start_date"); v != "" {
		t, err := utils.ParseMonthYearToTime(v)
		if err == nil {
			startDate = &t
		} else {
			log.Error("invalid start_date", slog.String("start_date", v), slog.Any("err", err))
			http.Error(w, "invalid start_date", http.StatusBadRequest)
			return
		}
	}
	if v := r.URL.Query().Get("end_date"); v != "" {
		t, err := utils.ParseMonthYearToTime(v)
		if err == nil {
			endDate = &t
		} else {
			log.Error("invalid end_date", slog.String("end_date", v), slog.Any("err", err))
			http.Error(w, "invalid end_date", http.StatusBadRequest)
			return
		}
	}

	cost, err := h.useCase.CalculateCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		log.Error("failed to calculate cost", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug("success calculate cost", slog.Int("cost", cost))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"cost": cost})
}
