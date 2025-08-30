package repository

import (
	"AggregationService/internal/domain/models/entity"
	"context"
	"github.com/google/uuid"
)

type ISubscriptionRepository interface {
	Create(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error)
	GetByID(ctx context.Context, id int) (*entity.Subscription, error)
	GetAll(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]*entity.Subscription, error)
	Update(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error)
	Delete(ctx context.Context, id int) error
	CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate string) (int, error)
}
