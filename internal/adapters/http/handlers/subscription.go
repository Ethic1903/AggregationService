package handlers

import (
	"AggregationService/internal/domain/models/dto"
	"AggregationService/internal/domain/usecase/subscription_usecase"
	"AggregationService/internal/pkg/logger"
	"AggregationService/internal/pkg/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type SubscriptionHandler struct {
	useCase subscription_usecase.ISubscriptionUseCase
}

func NewSubscriptionHandler(useCase subscription_usecase.ISubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{useCase: useCase}
}

// @Summary Создать подписку
// @Description Создает новую подписку для пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400,422 {string} string "Ошибка"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", "err", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sub, err := h.useCase.Create(ctx, &req)
	if err != nil {
		log.Error("failed to create subscription", "err", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Error("failed to encode")
	}
}

// @Summary Получить подписку по ID
// @Description Получить подписку по её идентификатору
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 404,500 {string} string "Ошибка"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id", "err", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := h.useCase.GetByID(ctx, id)
	if err != nil {
		log.Error("not found", "err", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Error("failed to encode")
	}
}

// @Summary Получить все подписки
// @Description Получить список всех подписок с фильтрами
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string false "Дата начала (mm-yyyy)"
// @Param end_date query string false "Дата окончания (mm-yyyy)"
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 400,500 {string} string "Ошибка"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var (
		userID      *uuid.UUID
		serviceName *string
		limit       int
		offset      int
	)

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			userID = &uid
		}
	}
	if sn := r.URL.Query().Get("service_name"); sn != "" {
		serviceName = &sn
	}
	limit = 100
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	subs, err := h.useCase.GetAll(ctx, userID, serviceName, limit, offset)
	if err != nil {
		log.Error("failed to get subscriptions", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		log.Error("failed to encode")
	}
}

// @Summary Обновить подписку
// @Description Обновить данные подписки
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body dto.UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400,404,422 {string} string "Ошибка"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id", id, "err", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", "err", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sub, err := h.useCase.Update(ctx, id, &req)
	if err != nil {
		log.Error("failed to update subscription", "err", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(sub); err != nil {
		log.Error("failed to encode")
	}
}

// @Summary Удалить подписку
// @Description Удалить подписку по ID
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Success 204 {string} string "Удалено"
// @Failure 404,500 {string} string "Ошибка"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id", "err", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.useCase.Delete(ctx, id); err != nil {
		log.Error("failed to delete subscription", "err", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Посчитать суммарную стоимость подписок
// @Description Посчитать суммарную стоимость подписок пользователя/сервиса за период
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string true "Дата начала периода (mm-yyyy)"
// @Param end_date query string true "Дата окончания периода (mm-yyyy)"
// @Success 200 {object} map[string]int
// @Failure 400,500 {string} string "Ошибка"
// @Router /subscriptions/cost [get]
func (h *SubscriptionHandler) CalculateCost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var (
		userID      *uuid.UUID
		serviceName *string
		startDate   *time.Time
		endDate     *time.Time
	)

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			userID = &uid
		}
	}
	if sn := r.URL.Query().Get("service_name"); sn != "" {
		serviceName = &sn
	}

	if sdStr := r.URL.Query().Get("start_date"); sdStr != "" {
		sd, err := utils.ParseMonthYearToTime(sdStr)
		if err != nil {
			http.Error(w, "invalid start_date", http.StatusBadRequest)
			return
		}
		startDate = &sd
	}
	if edStr := r.URL.Query().Get("end_date"); edStr != "" {
		ed, err := utils.ParseMonthYearToTime(edStr)
		if err != nil {
			http.Error(w, "invalid end_date", http.StatusBadRequest)
			return
		}
		endDate = &ed
	}

	cost, err := h.useCase.CalculateCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		log.Error("failed to calculate cost", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(map[string]int{"cost": cost}); err != nil {
		log.Error("failed to encode")
	}
}
