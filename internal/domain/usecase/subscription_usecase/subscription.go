package subscription_usecase

import (
	"AggregationService/internal/domain/models/dto"
	custom_err "AggregationService/internal/errors"
	"AggregationService/internal/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (u *subscriptionUseCase) Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log := logger.FromContext(ctx)
	log.Debug(fmt.Sprintf("trying to create subscription: %+v", req))

	if err := u.validator.Validate(req); err != nil {
		log.Error(fmt.Sprintf("invalid input: %v", custom_err.ErrInvalidRequest))
		return nil, custom_err.ErrInvalidRequest
	}

	entitySub := u.converter.ToSubscriptionEntity(req)
	entitySub.CreatedAt = time.Now()
	entitySub.UpdatedAt = entitySub.CreatedAt

	createdSub, err := u.subscriptionRepository.Create(ctx, entitySub)
	if err != nil {
		if errors.Is(err, custom_err.ErrSubscriptionAlreadyFound) {
			log.Error(fmt.Sprintf("duplicate subscription: %v", err))
			return nil, custom_err.ErrSubscriptionAlreadyFound
		}
		log.Error(fmt.Sprintf("failed to create subscription: %v", err))
		return nil, custom_err.ErrInvalidRequest
	}

	log.Debug(fmt.Sprintf("success create subscription: %+v", createdSub))
	return u.converter.ToSubscriptionDTO(createdSub), nil
}

func (u *subscriptionUseCase) Update(ctx context.Context, id int, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log := logger.FromContext(ctx)
	log.Debug(fmt.Sprintf("trying to update subscription: id=%d", req))

	if err := u.validator.Validate(req); err != nil {
		log.Error(fmt.Sprintf("invalid input: %v", custom_err.ErrInvalidRequest))
		return nil, custom_err.ErrInvalidRequest
	}

	sub, err := u.subscriptionRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_err.ErrSubscriptionNotFound) {
			log.Error(fmt.Sprintf("failed to get subscription for update: %v", err))
			return nil, custom_err.ErrSubscriptionNotFound
		}
		log.Error(fmt.Sprintf("failed to get subscription for update: %v", err))
		return nil, custom_err.ErrInvalidRequest
	}

	u.converter.ApplyUpdateToEntity(sub, req)
	sub.UpdatedAt = time.Now()

	updatedSub, err := u.subscriptionRepository.Update(ctx, sub)
	if err != nil {
		log.Error(fmt.Sprintf("failed to update subscription: %v", err))
		return nil, custom_err.ErrInvalidRequest
	}

	log.Debug(fmt.Sprintf("success update subscription: id=%d", id))
	return u.converter.ToSubscriptionDTO(updatedSub), nil
}

func (u *subscriptionUseCase) GetByID(ctx context.Context, id int) (*dto.SubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log := logger.FromContext(ctx)
	log.Debug(fmt.Sprintf("trying to get subscription by id: %d", id))

	sub, err := u.subscriptionRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_err.ErrSubscriptionNotFound) {
			log.Error(fmt.Sprintf("failed to get subscription by id: %v", err))
			return nil, custom_err.ErrSubscriptionNotFound
		}
		log.Error(fmt.Sprintf("failed to get subscription by id: %v", err))
		return nil, custom_err.ErrSubscriptionNotFound
	}

	log.Debug(fmt.Sprintf("success get subscription by id: %d", id))
	return u.converter.ToSubscriptionDTO(sub), nil
}

func (u *subscriptionUseCase) GetAll(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]*dto.SubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log := logger.FromContext(ctx)

	subs, err := u.subscriptionRepository.GetAll(ctx, userID, serviceName, limit, offset)
	if err != nil {
		if errors.Is(err, custom_err.ErrNoSubscriptionsFound) {
			log.Error(fmt.Sprintf("no subscriptions found: %v", err))
			return nil, custom_err.ErrNoSubscriptionsFound
		}
		log.Error(fmt.Sprintf("failed to get subscriptions: %v", err))
		return nil, custom_err.ErrInternalServer
	}

	result := u.converter.ToSubscriptionDTOs(subs)
	log.Debug(fmt.Sprintf("success getting subscriptions"))
	return result, nil
}

func (u *subscriptionUseCase) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log := logger.FromContext(ctx)
	log.Debug(fmt.Sprintf("trying to delete subscription: id=%d", id))

	if err := u.subscriptionRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, custom_err.ErrSubscriptionNotFound) {
			log.Error(fmt.Sprintf("failed to delete subscription: %v", err))
			return custom_err.ErrSubscriptionNotFound
		}
		log.Error(fmt.Sprintf("failed to delete subscription: %v", err))
		return custom_err.ErrSubscriptionNotFound
	}

	log.Debug(fmt.Sprintf("success delete subscription: id=%d", id))
	return nil
}

func (u *subscriptionUseCase) CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate *time.Time) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log := logger.FromContext(ctx)
	log.Debug(fmt.Sprintf("trying to calculate cost"))

	cost, err := u.subscriptionRepository.CalculateCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		log.Error(fmt.Sprintf("failed to calculate cost: %v", err))
		return 0, custom_err.ErrInternalServer
	}

	log.Debug(fmt.Sprintf("success calculating cost: %d", cost))
	return cost, nil
}
