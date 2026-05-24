package test

import (
	"context"
	"errors"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/google/uuid"
)

func TestService_GetSubscriptionByID(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		id        int
		seed      map[int]domain.Subscription
		wantErr   bool
		wantPrice int
	}{
		{
			name: "success",
			id:   1,
			seed: map[int]domain.Subscription{
				1: {ID: 1, Price: 100, UserID: userID},
			},
			wantErr:   false,
			wantPrice: 100,
		},
		{
			name:    "not found",
			id:      999,
			seed:    map[int]domain.Subscription{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc, _ := NewTestService(tt.seed)

			sub, err := svc.GetSubscriptionByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if !errors.Is(err, domain.ErrNotFound) {
					t.Fatalf("expected ErrNotFound, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if sub.Price != tt.wantPrice {
				t.Fatalf("expected %d, got %d", tt.wantPrice, sub.Price)
			}
		})
	}
}
