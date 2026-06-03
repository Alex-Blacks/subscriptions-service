package handler

import (
	"net/http"

	"github.com/Alex-Blacks/subscriptions/internal/logging"
	"github.com/Alex-Blacks/subscriptions/internal/service"
	"github.com/Alex-Blacks/subscriptions/internal/transport/dto"
)

// CreateSubscriptionHandler godoc
//
// @Summary Create subscription
// @Description Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionRequest true "subscription payload"
// @Success 201 {object} dto.SubscriptionIDResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [post]
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

		subReq, err := dto.SubscriptionToDomain(req)
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}

		id, err := svc.CreateSubscription(ctx, subReq)
		if err != nil {
			WriteDomainError(w, logger, err, req)
			return
		}

		WriteJSON(w, logger, http.StatusCreated, dto.SubscriptionIDResponse{ID: id})
	}
}

// GetSubscriptionByIDHandler godoc
//
// @Summary Get subscription by ID
// @Description Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "subscription ID"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [get]
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

// DeleteSubscriptionHandler godoc
//
// @Summary Delete subscription by ID
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path int true "subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [delete]
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

// UpdateSubscriptionHandler godoc
//
// @Summary Update subscription
// @Description Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "subscription ID"
// @Param request body dto.UpdateSubscriptionRequest true "subscription payload"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [patch]
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

		subReq, err := dto.UpdateSubscriptionToDomain(req)
		if err != nil {
			WriteError(w, logger, http.StatusBadRequest, err.Error())
			return
		}
		sub, err := svc.UpdateSubscription(ctx, id, subReq)
		if err != nil {
			WriteDomainError(w, logger, err, map[string]any{
				"id": id,
			})
			return
		}

		WriteJSON(w, logger, http.StatusOK, dto.SubscriptionToResponse(sub))
	}
}

// ListSubscriptionHandler godoc
//
// @Summary List subscriptions
// @Description List subscriptions with filters
// @Tags subscriptions
// @Produce json
//
// @Param service_name query string false "service_name"
// @Param user_id query string false "user_id (uuid)"
// @Param from query string false "from (MM-YYYY)"
// @Param to query string false "to (MM-YYYY)"
// @Param limit query int false "limit (default 50)"
// @Param offset query int false "offset (default 0)"
//
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
//
// @Router /subscriptions [get]
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

// SumSubscriptionPriceHandler godoc
//
// @Summary Sum price subscriptions
// @Description Sum price subscriptions
// @Tags subscriptions
// @Produce json
// @Param service_name query string false "service_name"
// @Param user_id query string false "user_id (uuid)"
// @Param from query string false "from (MM-YYYY)"
// @Param to query string false "to (MM-YYYY)"
// @Success 200 {object} dto.SubscriptionSumPriceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/sum [get]
func SumSubscriptionPriceHandler(svc service.SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		filter, err := ParseSumFilter(r)
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
