package test

import (
	"context"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
)

func TestService_ListSubscription(t *testing.T) {
	tests := []struct {
		name      string
		seed      map[int]domain.Subscription
		wantCount int
	}{
		{
			name:      "empty",
			seed:      map[int]domain.Subscription{},
			wantCount: 0,
		},
		{
			name: "two items",
			seed: map[int]domain.Subscription{
				1: {ID: 1, Price: 100},
				2: {ID: 2, Price: 200},
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc, _ := NewTestService(tt.seed)

			res, err := svc.ListSubscription(context.Background(), domain.ListFilter{})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(res) != tt.wantCount {
				t.Fatalf("expected %d, got %d", tt.wantCount, len(res))
			}
		})
	}
}
