package subscription_usecase

import (
	"AggregationService/internal/converters"
	"AggregationService/internal/domain/models/dto"
	"AggregationService/internal/domain/models/entity"
	"AggregationService/internal/domain/ports/repository"
	"AggregationService/internal/pkg/validation"
	"context"
	"github.com/google/uuid"
)

type ISubscriptionUseCase interface {
	Create(ctx context.Context, req dto.CreateSubscriptionRequest) (*entity.Subscription, error)
	GetByID(ctx context.Context, id int) (*entity.Subscription, error)
	GetAll(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]*entity.Subscription, error)
	Update(ctx context.Context, req dto.UpdateSubscriptionRequest) (*entity.Subscription, error)
	Delete(ctx context.Context, id int) error
	CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate string) (int, error)
}

type subscriptionUseCase struct {
	subscriptionRepository repository.ISubscriptionRepository
	validator              *validation.Validator
	converter              *converters.SubscriptionConverter
}

func New(
	subscriptionRepository repository.ISubscriptionRepository,
	validator *validation.Validator,
	converter *converters.SubscriptionConverter,
) ISubscriptionUseCase {
	return &subscriptionUseCase{
		subscriptionRepository: subscriptionRepository,
		validator:              validator,
		converter:              converter,
	}
}
