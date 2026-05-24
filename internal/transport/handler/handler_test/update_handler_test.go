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

func TestUpdateHandler(t *testing.T) {
	price := 200

	svc := &MockService{
		UpdateFn: func(ctx context.Context, id int, input domain.UpdateSubscriptionInput) (domain.Subscription, error) {
			return domain.Subscription{
				ID:    id,
				Price: *input.Price,
			}, nil
		},
	}

	r := chi.NewRouter()
	r.Patch("/subscriptions/{id}", handler.UpdateSubscriptionHandler(svc))

	body := `{"price":200}`

	req := httptest.NewRequest("PATCH", "/subscriptions/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	_ = price
}
