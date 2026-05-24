package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/transport/handler"
	"github.com/go-chi/chi/v5"
)

func TestListHandler(t *testing.T) {
	svc := &MockService{
		ListFn: func(ctx context.Context, filter domain.ListFilter) ([]domain.Subscription, error) {
			return []domain.Subscription{
				{ID: 1, Price: 100},
				{ID: 2, Price: 200},
			}, nil
		},
	}

	r := chi.NewRouter()
	r.Get("/subscriptions", handler.ListSubscriptionHandler(svc))

	req := httptest.NewRequest("GET", "/subscriptions", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
