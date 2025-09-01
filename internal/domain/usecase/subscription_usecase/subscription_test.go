package subscription_usecase

import (
	"AggregationService/internal/converters"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"AggregationService/internal/domain/models/dto"
	"AggregationService/internal/domain/models/entity"
	"AggregationService/internal/domain/ports/repository/mocks"
	custom_err "AggregationService/internal/errors"
	"AggregationService/internal/pkg/validation"
	"errors"
)

func Test_CreateSubscription(t *testing.T) {
	t.Parallel()

	validUUID := uuid.New()
	endDate := "12-2025"

	tests := []struct {
		name       string
		input      dto.CreateSubscriptionRequest
		setupMocks func(repo *mocks.ISubscriptionRepository)
		wantErr    error
	}{
		{
			name: "Valid subscription",
			input: dto.CreateSubscriptionRequest{
				UserID:      validUUID,
				ServiceName: "yandex",
				Price:       299,
				StartDate:   "09-2025",
				EndDate:     &endDate,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("Create", mock.Anything, mock.Anything).
					Return(&entity.Subscription{ID: 1}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Invalid price",
			input: dto.CreateSubscriptionRequest{
				UserID:      validUUID,
				ServiceName: "yandex",
				Price:       0,
				StartDate:   "09-2025",
				EndDate:     &endDate,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {},
			wantErr:    custom_err.ErrInvalidRequest,
		},
		{
			name: "Invalid service name",
			input: dto.CreateSubscriptionRequest{
				UserID:      validUUID,
				ServiceName: "",
				Price:       299,
				StartDate:   "09-2025",
				EndDate:     &endDate,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {},
			wantErr:    custom_err.ErrInvalidRequest,
		},
		{
			name: "Duplicate subscription",
			input: dto.CreateSubscriptionRequest{
				UserID:      validUUID,
				ServiceName: "yandex",
				Price:       299,
				StartDate:   "09-2025",
				EndDate:     &endDate,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("Create", mock.Anything, mock.Anything).
					Return(nil, custom_err.ErrSubscriptionAlreadyFound)
			},
			wantErr: custom_err.ErrSubscriptionAlreadyFound,
		},
		{
			name: "Repository error",
			input: dto.CreateSubscriptionRequest{
				UserID:      validUUID,
				ServiceName: "yandex",
				Price:       299,
				StartDate:   "09-2025",
				EndDate:     &endDate,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("Create", mock.Anything, mock.Anything).
					Return(nil, custom_err.ErrInternalServer)
			},
			wantErr: custom_err.ErrInternalServer,
		},
		{
			name: "Invalid UUID",
			input: dto.CreateSubscriptionRequest{
				UserID:      uuid.Nil,
				ServiceName: "yandex",
				Price:       299,
				StartDate:   "09-2025",
				EndDate:     &endDate,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {},
			wantErr:    custom_err.ErrInvalidRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewISubscriptionRepository(t)
			validator, _ := validation.New()
			converter := converters.New()
			useCase := New(mockRepo, validator, converter)

			tt.setupMocks(mockRepo)

			ctx := context.Background()
			_, err := useCase.Create(ctx, &tt.input)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error: %v, got: %v", tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_GetByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		id         int
		setupMocks func(repo *mocks.ISubscriptionRepository)
		wantErr    error
	}{
		{
			name: "Found",
			id:   1,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetByID", mock.Anything, 1).
					Return(&entity.Subscription{ID: 1, ServiceName: "yandex"}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Not found",
			id:   999,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetByID", mock.Anything, 999).
					Return(nil, custom_err.ErrSubscriptionNotFound)
			},
			wantErr: custom_err.ErrSubscriptionNotFound,
		},
		{
			name: "Repository error",
			id:   2,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetByID", mock.Anything, 2).
					Return(nil, custom_err.ErrInternalServer)
			},
			wantErr: custom_err.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewISubscriptionRepository(t)
			validator, _ := validation.New()
			converter := converters.New()
			useCase := New(mockRepo, validator, converter)

			tt.setupMocks(mockRepo)

			ctx := context.Background()
			_, err := useCase.GetByID(ctx, tt.id)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error: %v, got: %v", tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_UpdateSubscription(t *testing.T) {
	t.Parallel()
	serviceName := "yandex plus"
	price := 399
	endDate := "12-2025"

	tests := []struct {
		name       string
		id         int
		input      dto.UpdateSubscriptionRequest
		setupMocks func(repo *mocks.ISubscriptionRepository)
		wantErr    error
	}{
		{
			name: "Valid update",
			id:   1,
			input: dto.UpdateSubscriptionRequest{
				ServiceName: &serviceName,
				Price:       &price,
				EndDate:     &endDate,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetByID", mock.Anything, 1).
					Return(&entity.Subscription{ID: 1, ServiceName: "yandex"}, nil)
				repo.On("Update", mock.Anything, mock.Anything).
					Return(&entity.Subscription{ID: 1, ServiceName: "yandex plus"}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Not found",
			id:   999,
			input: dto.UpdateSubscriptionRequest{
				ServiceName: &serviceName,
				Price:       &price,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetByID", mock.Anything, 999).
					Return(nil, custom_err.ErrSubscriptionNotFound)
			},
			wantErr: custom_err.ErrSubscriptionNotFound,
		},
		{
			name: "Repository error",
			id:   1,
			input: dto.UpdateSubscriptionRequest{
				ServiceName: &serviceName,
				Price:       &price,
			},
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetByID", mock.Anything, 1).
					Return(&entity.Subscription{ID: 1, ServiceName: "yandex"}, nil)
				repo.On("Update", mock.Anything, mock.Anything).
					Return(nil, custom_err.ErrInternalServer)
			},
			wantErr: custom_err.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewISubscriptionRepository(t)
			validator, _ := validation.New()
			converter := converters.New()
			useCase := New(mockRepo, validator, converter)

			tt.setupMocks(mockRepo)

			ctx := context.Background()
			_, err := useCase.Update(ctx, tt.id, &tt.input)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error: %v, got: %v", tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_DeleteSubscription(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		id         int
		setupMocks func(repo *mocks.ISubscriptionRepository)
		wantErr    error
	}{
		{
			name: "Valid delete",
			id:   1,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("Delete", mock.Anything, 1).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Not found",
			id:   999,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("Delete", mock.Anything, 999).
					Return(custom_err.ErrSubscriptionNotFound)
			},
			wantErr: custom_err.ErrSubscriptionNotFound,
		},
		{
			name: "Repository error",
			id:   2,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("Delete", mock.Anything, 2).
					Return(custom_err.ErrInternalServer)
			},
			wantErr: custom_err.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewISubscriptionRepository(t)
			validator, _ := validation.New()
			converter := converters.New()
			useCase := New(mockRepo, validator, converter)

			tt.setupMocks(mockRepo)

			ctx := context.Background()
			err := useCase.Delete(ctx, tt.id)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error: %v, got: %v", tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_GetAllSubscriptions(t *testing.T) {
	t.Parallel()
	validUUID := uuid.New()
	serviceName := "yandex"
	tests := []struct {
		name       string
		userID     *uuid.UUID
		service    *string
		limit      int
		offset     int
		setupMocks func(repo *mocks.ISubscriptionRepository)
		wantErr    error
	}{
		{
			name:    "Found",
			userID:  &validUUID,
			service: &serviceName,
			limit:   10,
			offset:  0,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetAll", mock.Anything, &validUUID, &serviceName, 10, 0).
					Return([]*entity.Subscription{
						{ID: 1, ServiceName: "yandex"},
						{ID: 2, ServiceName: "yandex plus"},
					}, nil)
			},
			wantErr: nil,
		},
		{
			name:    "Empty result",
			userID:  &validUUID,
			service: &serviceName,
			limit:   10,
			offset:  0,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetAll", mock.Anything, &validUUID, &serviceName, 10, 0).
					Return([]*entity.Subscription{}, nil)
			},
			wantErr: nil,
		},
		{
			name:    "Repository error",
			userID:  &validUUID,
			service: &serviceName,
			limit:   10,
			offset:  0,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("GetAll", mock.Anything, &validUUID, &serviceName, 10, 0).
					Return(nil, custom_err.ErrInternalServer)
			},
			wantErr: custom_err.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewISubscriptionRepository(t)
			validator, _ := validation.New()
			converter := converters.New()
			useCase := New(mockRepo, validator, converter)

			tt.setupMocks(mockRepo)

			ctx := context.Background()
			_, err := useCase.GetAll(ctx, tt.userID, tt.service, tt.limit, tt.offset)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error: %v, got: %v", tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_CalculateCost(t *testing.T) {
	t.Parallel()
	validUUID := uuid.New()
	serviceName := "yandex"
	startDate := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		userID     *uuid.UUID
		service    *string
		startDate  *time.Time
		endDate    *time.Time
		setupMocks func(repo *mocks.ISubscriptionRepository)
		wantCost   int
		wantErr    error
	}{
		{
			name:      "Valid cost",
			userID:    &validUUID,
			service:   &serviceName,
			startDate: &startDate,
			endDate:   &endDate,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("CalculateCost", mock.Anything, &validUUID, &serviceName, &startDate, &endDate).
					Return(451, nil)
			},
			wantCost: 451,
			wantErr:  nil,
		},
		{
			name:      "Repository error",
			userID:    &validUUID,
			service:   &serviceName,
			startDate: &startDate,
			endDate:   &endDate,
			setupMocks: func(repo *mocks.ISubscriptionRepository) {
				repo.On("CalculateCost", mock.Anything, &validUUID, &serviceName, &startDate, &endDate).
					Return(0, custom_err.ErrInternalServer)
			},
			wantCost: 0,
			wantErr:  custom_err.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewISubscriptionRepository(t)
			validator, _ := validation.New()
			converter := converters.New()
			useCase := New(mockRepo, validator, converter)

			tt.setupMocks(mockRepo)

			ctx := context.Background()
			cost, err := useCase.CalculateCost(ctx, tt.userID, tt.service, tt.startDate, tt.endDate)
			assert.Equal(t, tt.wantCost, cost)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected error: %v, got: %v", tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
