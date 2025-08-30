package postgres

import (
	"AggregationService/internal/domain/models/entity"
	"AggregationService/internal/domain/ports/repository"
	errors_custom "AggregationService/internal/errors"
	"AggregationService/internal/infrastructure/database/go_postgres"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"time"
)

const _tableSubscriptions = "subscriptions"

type subscriptionsRepository struct {
	client *go_postgres.PostgresClient
}

func NewSubscriptionsRepository(client *go_postgres.PostgresClient) repository.ISubscriptionRepository {
	return &subscriptionsRepository{client: client}
}

func (s *subscriptionsRepository) Create(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error) {
	const op = "repository.postgres.Create"
	sq := s.client.Builder.
		Insert(_tableSubscriptions).
		Columns(
			"id",
			"service_name",
			"price",
			"user_id",
			"start_date",
			"end_date",
			"created_at",
			"updated_at",
		).
		Values(
			subscription.ID,
			subscription.ServiceName,
			subscription.Price,
			subscription.UserID,
			subscription.StartDate,
			subscription.EndDate,
			subscription.CreatedAt,
			subscription.UpdatedAt,
		).
		Suffix(`RETURNING id, created_at`)
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: to sql: %w", op, err)
	}

	var id int
	var createdAt time.Time
	if err = s.client.DB.QueryRowxContext(ctx, query, args...).Scan(&id, &createdAt); err != nil {
		return nil, fmt.Errorf("%s: to scan: %w", op, err)
	}
	subscription.ID = id
	subscription.CreatedAt = createdAt
	return subscription, nil
}

func (s *subscriptionsRepository) GetByID(ctx context.Context, id int) (*entity.Subscription, error) {
	const op = "repository.postgres.GetByID"
	var sub *entity.Subscription
	sq := s.client.Builder.
		Select("*").
		From(_tableSubscriptions).
		Where(squirrel.Eq{"id": id})
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: to sql: %w", op, err)
	}

	if err = s.client.DB.GetContext(ctx, &sub, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors_custom.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	return sub, nil
}

func (s *subscriptionsRepository) GetAll(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]*entity.Subscription, error) {
	const op = "repository.postgres.GetAll"
	sq := s.client.Builder.
		Select("*").
		From(_tableSubscriptions)
	if userID != nil {
		sq = sq.Where(squirrel.Eq{"user_id": userID})
	}
	if serviceName != nil {
		sq = sq.Where(squirrel.ILike{"service_name": "%" + *serviceName + "%"})
	}
	sq = sq.Limit(uint64(limit)).Offset(uint64(offset))
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: to sql: %w", op, err)
	}
	var subs []*entity.Subscription
	if err = s.client.DB.SelectContext(ctx, &subs, query, args...); err != nil {
		return nil, fmt.Errorf("%s: select: %w", op, err)
	}
	if len(subs) == 0 {
		return nil, errors_custom.ErrNoSubscriptionsFound
	}
	return subs, nil
}

func (s *subscriptionsRepository) Update(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error) {
	const op = "repository.postgres.Update"
	sq := s.client.Builder.
		Update(_tableSubscriptions).
		Set("service_name", subscription.ServiceName).
		Set("price", subscription.Price).
		Set("end_date", subscription.EndDate).
		Set("updated_at", subscription.UpdatedAt).
		Where(squirrel.Eq{"id": subscription.ID}).
		Suffix("RETURNING updated_at")

	query, args, err := sq.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: to sql: %w", op, err)
	}
	var updatedAt time.Time
	if err = s.client.DB.QueryRowxContext(ctx, query, args...).Scan(&updatedAt); err != nil {
		return nil, fmt.Errorf("%s: to scan: %w", op, err)
	}
	return subscription, nil
}

func (s *subscriptionsRepository) Delete(ctx context.Context, id int) error {
	const op = "repository.postgres.Delete"
	sq := s.client.Builder.
		Delete(_tableSubscriptions).
		Where(squirrel.Eq{"id": id})
	query, args, err := sq.ToSql()
	if err != nil {
		return fmt.Errorf("%s: to sql: %w", op, err)
	}

	res, err := s.client.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: to delete: %w", op, err)
	}

	affectedRows, _ := res.RowsAffected()
	if affectedRows == 0 {
		return errors_custom.ErrSubscriptionNotFound
	}
	return nil
}

func (s *subscriptionsRepository) CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate string) (int, error) {
	const op = "repository.postgres.CalculateCost"
	sq := s.client.Builder.
		Select("COALESCE(SUM(price),0)").
		From(_tableSubscriptions).
		Where(squirrel.LtOrEq{"start_date": endDate}).
		Where(squirrel.Or{
			squirrel.Eq{"end_date": nil},
			squirrel.GtOrEq{"end_date": startDate},
		})
	if userID != nil {
		sq = sq.Where(squirrel.Eq{"user_id": *userID})
	}
	if serviceName != nil {
		sq = sq.Where(squirrel.ILike{"service_name": "%" + *serviceName + "%"})
	}
	query, args, err := sq.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: to sql: %w", op, err)
	}

	var totalCost int
	if err = s.client.DB.GetContext(ctx, &totalCost, query, args...); err != nil {
		return 0, fmt.Errorf("%s: to extract total cost: %w", op, err)
	}
	return totalCost, nil
}
