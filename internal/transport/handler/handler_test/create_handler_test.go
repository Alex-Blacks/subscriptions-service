package test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/transport/handler"
	"github.com/go-chi/chi/v5"
)

func TestCreateSubscriptionHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockFn     func(context.Context, domain.CreateSubscriptionInput) (int, error)
		wantStatus int
	}{
		{
			name: "success",
			body: `{
				"service_name":"yandex",
				"price":100,
				"user_id":"550e8400-e29b-41d4-a716-446655440000",
				"start_date":"05-2026"
			}`,
			mockFn: func(ctx context.Context, input domain.CreateSubscriptionInput) (int, error) {
				return 1, nil
			},
			wantStatus: 201,
		},
		{
			name: "invalid json",
			body: "{",
			mockFn: func(ctx context.Context, input domain.CreateSubscriptionInput) (int, error) {
				return 0, nil
			},
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := &MockService{
				CreateFn: tt.mockFn,
			}

			r := chi.NewRouter()
			r.Post("/subscriptions", handler.CreateSubscriptionHandler(svc))

			req := httptest.NewRequest("POST", "/subscriptions", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected %d got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
