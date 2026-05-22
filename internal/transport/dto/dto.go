package dto

import (
	"time"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func SubscriptionToDomain(input CreateSubscriptionRequest) domain.CreateSubscriptionInput {
	return domain.CreateSubscriptionInput{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}
}

type SubscriptionIDResponse struct {
	ID int `json:"id"`
}

type SubscriptionResponse struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
}

func SubscriptionToResponse(input domain.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     *input.EndDate,
		CreatedAt:   input.CreatedAt,
	}
}
