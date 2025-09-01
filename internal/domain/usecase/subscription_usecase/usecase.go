package subscription_usecase

import (
	"AggregationService/internal/converters"
	"AggregationService/internal/domain/models/dto"
	"AggregationService/internal/domain/ports/repository"
	"AggregationService/internal/pkg/validation"
	"context"
	"github.com/google/uuid"
	"time"
)

type ISubscriptionUseCase interface {
	Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	GetByID(ctx context.Context, id int) (*dto.SubscriptionResponse, error)
	GetAll(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]*dto.SubscriptionResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	Delete(ctx context.Context, id int) error
	CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate *time.Time) (int, error)
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
