package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/transport/handler"
	"github.com/go-chi/chi/v5"
)

func TestGetSubscriptionHandler(t *testing.T) {
	svc := &MockService{
		GetFn: func(ctx context.Context, id int) (domain.Subscription, error) {
			return domain.Subscription{
				ID:    id,
				Price: 100,
			}, nil
		},
	}

	r := chi.NewRouter()
	r.Get("/subscriptions/{id}", handler.GetSubscriptionByIDHandler(svc))

	req := httptest.NewRequest("GET", "/subscriptions/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
