package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" validate:"required, min=1, max=255"`
	Price       int       `json:"price" validate:"required, min=1"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     time.Time `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type CreateSubscriptionResponse struct {
	ServiceName *string   `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
	Price       *int      `json:"price,omitempty" validate:"omitempty,min=1"`
	EndDate     time.Time `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
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
	StartDate   time.Time  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     time.Time  `json:"end_date" validate:"required,datetime=01-2006"`
}

type CalculateCostResponse struct {
	TotalCost int `json:"total_cost"`
}
