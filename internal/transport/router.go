package transport

import (
	"github.com/Alex-Blacks/subscriptions/internal/logging"
	"github.com/Alex-Blacks/subscriptions/internal/service"
	"github.com/Alex-Blacks/subscriptions/internal/transport/handler"
	"github.com/Alex-Blacks/subscriptions/internal/transport/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	logger := logging.NewLogger()

	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.LoggerMiddleware(logger))

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", handler.CreateSubscriptionHandler(svc))
		r.Get("/{id}", handler.GetSubscriptionByIDHandler(svc))
		r.Delete("/{id}", handler.DeleteSubscriptionHandler(svc))
		r.Patch("/{id}", handler.UpdateSubscriptionHandler(svc))
		r.Get("/", handler.ListSubscriptionHandler(svc))
		r.Get("/sum", handler.SumSubscriptionPriceHandler(svc))
	})

	return router
}
