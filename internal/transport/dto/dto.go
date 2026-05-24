package dto

import (
	"fmt"
	"time"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

func SubscriptionToDomain(input CreateSubscriptionRequest) (domain.CreateSubscriptionInput, error) {
	start_date, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return domain.CreateSubscriptionInput{}, fmt.Errorf("invalid start_date from (MM-YYYY): %w", err)
	}
	var end_date *time.Time
	if input.EndDate != nil {
		parsed, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return domain.CreateSubscriptionInput{}, fmt.Errorf("invalid end_date from (MM-YYYY): %w", err)
		}
		end_date = &parsed
	}
	return domain.CreateSubscriptionInput{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   start_date,
		EndDate:     end_date,
	}, nil
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *int       `json:"price,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	StartDate   *string    `json:"start_date,omitempty"`
	EndDate     *string    `json:"end_date,omitempty"`
}

func UpdateSubscriptionToDomain(input UpdateSubscriptionRequest) (domain.UpdateSubscriptionInput, error) {
	var start_date *time.Time
	if input.StartDate != nil {
		parsed, err := time.Parse("01-2006", *input.StartDate)
		if err != nil {
			return domain.UpdateSubscriptionInput{}, fmt.Errorf("invalid update start_date from (MM-YYYY): %w", err)
		}

		start_date = &parsed
	}
	var end_date *time.Time
	if input.EndDate != nil {
		parsed, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return domain.UpdateSubscriptionInput{}, fmt.Errorf("invalid update end_date from (MM-YYYY): %w", err)
		}
		end_date = &parsed
	}
	return domain.UpdateSubscriptionInput{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   start_date,
		EndDate:     end_date,
	}, nil
}

func ToListResponse(input []domain.Subscription) []SubscriptionResponse {
	resp := make([]SubscriptionResponse, len(input))

	for i, p := range input {
		startDate := p.StartDate.Format("01-2006")

		var endDate *string

		if p.EndDate != nil {
			formatted := p.EndDate.Format("01-2006")
			endDate = &formatted
		}
		resp[i] = SubscriptionResponse{
			ID:          p.ID,
			ServiceName: p.ServiceName,
			Price:       p.Price,
			UserID:      p.UserID,
			StartDate:   startDate,
			EndDate:     endDate,
			CreatedAt:   p.CreatedAt,
		}
	}
	return resp
}

type SubscriptionSumPriceResponse struct {
	SumPrice int `json:"sum_price"`
}

type SubscriptionIDResponse struct {
	ID int `json:"id"`
}

type SubscriptionResponse struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func SubscriptionToResponse(input domain.Subscription) SubscriptionResponse {
	startDate := input.StartDate.Format("01-2006")

	var endDate *string

	if input.EndDate != nil {
		formatted := input.EndDate.Format("01-2006")
		endDate = &formatted
	}
	return SubscriptionResponse{
		ID:          input.ID,
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   input.CreatedAt,
	}
}
