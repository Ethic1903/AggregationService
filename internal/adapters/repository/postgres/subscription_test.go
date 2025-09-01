package postgres

import (
	"AggregationService/internal/domain/ports/repository"
	"context"
	"testing"
	"time"

	"AggregationService/internal/domain/models/entity"
	"AggregationService/internal/infrastructure/database/go_postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestRepo(t *testing.T) repository.ISubscriptionRepository {
	client, err := go_postgres.NewTestClient()
	assert.NoError(t, err)
	return NewSubscriptionsRepository(client)
}

func TestSubscriptionRepository_Create(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()
	sub := &entity.Subscription{
		ServiceName: "yandex",
		Price:       299,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	created, err := repo.Create(ctx, sub)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, sub.ServiceName, created.ServiceName)
	assert.NotZero(t, created.ID)
}

func TestSubscriptionRepository_GetByID(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()
	sub := &entity.Subscription{
		ServiceName: "yandex",
		Price:       299,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}
	created, _ := repo.Create(ctx, sub)

	found, err := repo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.ServiceName, found.ServiceName)
}

func TestSubscriptionRepository_Update(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()
	sub := &entity.Subscription{
		ServiceName: "yandex",
		Price:       299,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}
	created, _ := repo.Create(ctx, sub)

	created.ServiceName = "yandex plus"
	created.Price = 399

	updated, err := repo.Update(ctx, created)
	assert.NoError(t, err)
	assert.Equal(t, "yandex plus", updated.ServiceName)
	assert.Equal(t, 399, updated.Price)
}

func TestSubscriptionRepository_Delete(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()
	sub := &entity.Subscription{
		ServiceName: "yandex",
		Price:       299,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}
	created, _ := repo.Create(ctx, sub)

	err := repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, created.ID)
	assert.Error(t, err)
}

func TestSubscriptionRepository_GetAll(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()
	userID := uuid.New()

	sub1 := &entity.Subscription{
		ServiceName: "yandex",
		Price:       299,
		UserID:      userID,
		StartDate:   time.Now(),
	}
	sub2 := &entity.Subscription{
		ServiceName: "yandex plus",
		Price:       399,
		UserID:      userID,
		StartDate:   time.Now(),
	}
	repo.Create(ctx, sub1)
	repo.Create(ctx, sub2)

	subs, err := repo.GetAll(ctx, &userID, nil, 10, 0)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(subs), 2)
}

func TestSubscriptionRepository_CalculateCost(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()
	userID := uuid.New()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 1, 0)

	sub := &entity.Subscription{
		ServiceName: "yandex",
		Price:       451,
		UserID:      userID,
		StartDate:   time.Now(),
	}
	repo.Create(ctx, sub)

	cost, err := repo.CalculateCost(ctx, &userID, nil, &startDate, &endDate)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, cost, 451)
}
