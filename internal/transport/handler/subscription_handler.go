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

		id, err := svc.CreateSubscription(ctx, dto.SubscriptionToDomain(req))
		if err != nil {
			WriteDomainError(w, logger, err, req)
			return
		}

		WriteJSON(w, logger, http.StatusCreated, dto.SubscriptionIDResponse{ID: id})
	}
}

func GetSubscriptionByIDHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		id, err := ParsePositiveIntParam(r, "id")
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		sub, err := svc.GetSubscriptionByID(ctx, id)
		if err != nil {
			WriteDomainError(w, logger, err, map[string]any{"id": id})
			return
		}

		WriteJSON(w, logger, http.StatusOK, dto.SubscriptionToResponse(sub))
	}
}

func DeleteSubscriptionHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		id, err := ParsePositiveIntParam(r, "id")
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		if err := svc.DeleteSubscription(ctx, id); err != nil {
			WriteDomainError(w, logger, err, map[string]any{"id": id})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateSubscriptionHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		id, err := ParsePositiveIntParam(r, "id")
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		var req dto.UpdateSubscriptionRequest

		if err := DecodeJSON(w, r, logger, &req); err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		if err := ValidateUpdateSubscription(req); err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		sub, err := svc.UpdateSubscription(ctx, id, dto.UpdateSubscriptionToDomain(req))
		if err != nil {
			WriteDomainError(w, logger, err, map[string]any{
				"id": id,
			})
			return
		}

		WriteJSON(w, logger, http.StatusOK, dto.SubscriptionToResponse(sub))
	}
}

func ListSubscriptionHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		filter, err := ParseListFilter(r)
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		list, err := svc.ListSubscription(ctx, filter)
		if err != nil {
			WriteDomainError(w, logger, err, map[string]any{
				"filter": filter,
			})
			return
		}

		WriteJSON(w, logger, http.StatusOK, dto.ToListResponse(list))
	}
}

func SumSubscriptionPriceHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		filter, err := ParseListFilter(r)
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		sum, err := svc.SumSubscriptionPrice(ctx, filter)
		if err != nil {
			WriteDomainError(w, logger, err, map[string]any{
				"filter": filter,
			})
			return
		}

		WriteJSON(w, logger, http.StatusOK, dto.SubscriptionSumPriceResponse{
			SumPrice: sum,
		})
	}
}
