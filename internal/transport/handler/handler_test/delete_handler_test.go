package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alex-Blacks/subscriptions/internal/transport/handler"
	"github.com/go-chi/chi/v5"
)

func TestDeleteHandler(t *testing.T) {
	svc := &MockService{
		DeleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}

	r := chi.NewRouter()
	r.Delete("/subscriptions/{id}", handler.DeleteSubscriptionHandler(svc))

	req := httptest.NewRequest("DELETE", "/subscriptions/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status %d", w.Code)
	}
}
