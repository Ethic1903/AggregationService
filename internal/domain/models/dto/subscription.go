package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" validate:"required,min=1,max=255"`
	Price       int       `json:"price" validate:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
	StartDate   string    `json:"start_date" validate:"required,mmYYYY"`
	EndDate     *string   `json:"end_date,omitempty" validate:"omitempty,mmYYYY"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,min=1"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,mmYYYY"`
}

type CreateSubscriptionResponse struct {
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,min=1"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,mmYYYY"`
}

type SubscriptionResponse struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CalculateCostRequest struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	ServiceName *string    `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
	StartDate   string     `json:"start_date" validate:"required,mmYYYY"`
	EndDate     string     `json:"end_date" validate:"required,mmYYYY"`
}

type CalculateCostResponse struct {
	TotalCost int `json:"total_cost"`
}
