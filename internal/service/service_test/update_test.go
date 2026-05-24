package test

import (
	"context"
	"errors"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
)

func TestService_UpdateSubscription(t *testing.T) {
	price := 200

	tests := []struct {
		name      string
		id        int
		seed      map[int]domain.Subscription
		update    domain.UpdateSubscriptionInput
		wantErr   bool
		wantPrice int
	}{
		{
			name: "success partial update",
			id:   1,
			seed: map[int]domain.Subscription{
				1: {ID: 1, Price: 100},
			},
			update: domain.UpdateSubscriptionInput{
				Price: &price,
			},
			wantErr:   false,
			wantPrice: 200,
		},
		{
			name:    "not found",
			id:      999,
			seed:    map[int]domain.Subscription{},
			update:  domain.UpdateSubscriptionInput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc, _ := NewTestService(tt.seed)

			sub, err := svc.UpdateSubscription(context.Background(), tt.id, tt.update)

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
