package handler

import (
	"net/http"

	"github.com/Alex-Blacks/subscriptions/internal/logging"
	"github.com/Alex-Blacks/subscriptions/internal/service"
	"github.com/Alex-Blacks/subscriptions/internal/transport/dto"
)

func CreateSubscriptionHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		var req dto.CreateSubscriptionRequest

		if err := DecodeJSON(w, r, logger, &req); err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		if err := ValidateCreateSubscription(req); err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		subID, err := svc.CreateSubscription(ctx, dto.SubscriptionToDomain(req))
		if err != nil {
			WriteDomainError(w, logger, err, req)
			return
		}
		resp := dto.SubscriptionIDResponse{ID: subID}
		WriteJSON(w, logger, http.StatusCreated, resp)
	}
}

func GetSubscriptionByIDHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		subID, err := ParsePositiveIntParam(r, "id")
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		sub, err := svc.GetSubscriptionByID(ctx, subID)
		if err != nil {
			WriteDomainError(w, logger, err, map[string]any{
				"subscription_id": subID,
			})
			return
		}

		WriteJSON(w, logger, http.StatusOK, dto.SubscriptionToResponse(sub))
	}
}

func DeleteSubscriptionHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		subID, err := ParsePositiveIntParam(r, "id")
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		if err := svc.DeleteSubscription(ctx, subID); err != nil {
			WriteDomainError(w, logger, err, map[string]any{
				"subscription_id": subID,
			})
			return
		}

		logger.Info("subscription deleted",
			"subscription_id", subID,
		)
		w.WriteHeader(http.StatusNoContent)
	}
}
