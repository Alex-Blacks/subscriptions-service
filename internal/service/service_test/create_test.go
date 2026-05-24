package test

import (
	"context"
	"testing"
	"time"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/google/uuid"
)

func TestService_CreateSubscription(t *testing.T) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	startDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		seed    map[int]domain.Subscription
		wantErr bool
		wantID  int
	}{
		{
			name:    "success",
			seed:    map[int]domain.Subscription{},
			wantErr: false,
			wantID:  1,
		},
		{
			name: "already exists",
			seed: map[int]domain.Subscription{
				1: {
					ServiceName: "yandex",
					Price:       100,
					UserID:      userID,
					StartDate:   startDate,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc, _ := NewTestService(tt.seed)

			input := domain.CreateSubscriptionInput{
				ServiceName: "yandex",
				Price:       100,
				UserID:      userID,
				StartDate:   startDate,
			}

			id, err := svc.CreateSubscription(context.Background(), input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id != tt.wantID {
				t.Fatalf("expected %d, got %d", tt.wantID, id)
			}
		})
	}
}
