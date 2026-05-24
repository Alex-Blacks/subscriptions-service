package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/transport/handler"
	"github.com/go-chi/chi/v5"
)

func TestSumHandler(t *testing.T) {
	svc := &MockService{
		SumFn: func(ctx context.Context, filter domain.SumFilter) (int, error) {
			return 350, nil
		},
	}

	r := chi.NewRouter()
	r.Get("/subscriptions/sum", handler.SumSubscriptionPriceHandler(svc))

	req := httptest.NewRequest("GET", "/subscriptions/sum", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
