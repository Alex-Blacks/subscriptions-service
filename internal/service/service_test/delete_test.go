package test

import (
	"context"
	"errors"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
)

func TestService_DeleteSubscription(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		seed    map[int]domain.Subscription
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			seed: map[int]domain.Subscription{
				1: {ID: 1, Price: 100},
			},
			wantErr: false,
		},
		{
			name:    "not found",
			id:      1,
			seed:    map[int]domain.Subscription{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc, storage := NewTestService(tt.seed)

			err := svc.DeleteSubscription(context.Background(), tt.id)

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

			if _, ok := storage.data[tt.id]; ok {
				t.Fatal("subscription was not deleted")
			}
		})
	}
}
