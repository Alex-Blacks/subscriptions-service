package test

import (
	"context"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
)

func TestService_SumSubscriptionPrice(t *testing.T) {
	tests := []struct {
		name    string
		seed    map[int]domain.Subscription
		wantSum int
	}{
		{
			name:    "empty",
			seed:    map[int]domain.Subscription{},
			wantSum: 0,
		},
		{
			name: "sum all",
			seed: map[int]domain.Subscription{
				1: {Price: 100},
				2: {Price: 250},
			},
			wantSum: 350,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc, _ := NewTestService(tt.seed)

			sum, err := svc.SumSubscriptionPrice(context.Background(), domain.SumFilter{})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if sum != tt.wantSum {
				t.Fatalf("expected %d, got %d", tt.wantSum, sum)
			}
		})
	}
}
