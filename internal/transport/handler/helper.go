package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/transport/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, logger *slog.Logger, dest any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dest); err != nil {
		logger.Warn("decode failed", "error", err)
		return fmt.Errorf("invalid json")
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		logger.Warn("multiple json objects", "error", err)
		return fmt.Errorf("body must contain single json object")
	}

	return nil
}

func ParsePositiveIntParam(r *http.Request, name string) (int, error) {
	valStr := chi.URLParam(r, name)
	if strings.TrimSpace(valStr) == "" {
		return 0, fmt.Errorf("%s must not be empty", name)
	}
	val, err := strconv.Atoi(valStr)
	if err != nil || val <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return val, nil
}
func ParseStrParam(r *http.Request, name string) (string, error) {
	valStr := chi.URLParam(r, name)
	if strings.TrimSpace(valStr) == "" {
		return "", fmt.Errorf("%s must not be empty", name)
	}
	return valStr, nil
}

func WriteJSON(w http.ResponseWriter, logger *slog.Logger, status int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("encoding response failed", "error", err)
	}
}

func WriteError(w http.ResponseWriter, logger *slog.Logger, status int, msg string) {
	logger.Warn("request failed",
		"status", status,
		"error", msg,
	)
	WriteJSON(w, logger, status, ErrorResponse{Error: msg})
}
func WriteInternalError(w http.ResponseWriter, logger *slog.Logger, err error, req any) {
	logger.Error("request failed",
		"error", err,
		"request", req,
	)
	WriteJSON(w, logger, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
}

func ParseListFilter(r *http.Request) (domain.ListFilter, error) {
	q := r.URL.Query()
	var filter domain.ListFilter

	if raw := q.Get("user_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return filter, fmt.Errorf("invalid user_id: %w", err)
		}
		filter.UserID = &id
	}

	if raw := q.Get("service_name"); raw != "" {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return filter, fmt.Errorf("service_name cannot be empty")
		}
		filter.ServiceName = &raw
	}

	if raw := q.Get("from"); raw != "" {
		t, err := time.Parse("2006-01", raw)
		if err != nil {
			return filter, fmt.Errorf("invalid from date (expected YYYY-MM): %w", err)
		}
		filter.From = &t
	}

	if raw := q.Get("to"); raw != "" {
		t, err := time.Parse("2006-01", raw)
		if err != nil {
			return filter, fmt.Errorf("invalid to date (expected YYYY-MM): %w", err)
		}
		filter.To = &t
	}

	if filter.From != nil && filter.To != nil {
		if filter.To.Before(*filter.From) {
			return filter, fmt.Errorf("to must be >= from")
		}
	}

	if raw := q.Get("limit"); raw != "" {
		l, err := strconv.Atoi(raw)
		if err != nil {
			return filter, fmt.Errorf("invalid limit")
		}
		if l <= 0 {
			return filter, fmt.Errorf("limit must be > 0")
		}
		if l > 100 {
			l = 100
		}
		filter.Limit = l
	} else {
		filter.Limit = 50
	}

	if raw := q.Get("offset"); raw != "" {
		o, err := strconv.Atoi(raw)
		if err != nil {
			return filter, fmt.Errorf("invalid offset")
		}
		if o < 0 {
			return filter, fmt.Errorf("offset must be >= 0")
		}
		filter.Offset = o
	}

	return filter, nil
}

func ValidateCreateSubscription(input dto.CreateSubscriptionRequest) error {
	if strings.TrimSpace(input.ServiceName) == "" {
		return fmt.Errorf("service_name must not be empty")
	}
	if input.Price <= 0 {
		return fmt.Errorf("invalid input: price must be > 0")
	}
	if input.UserID == uuid.Nil {
		return fmt.Errorf("user_id must not be empty")
	}
	if input.StartDate.IsZero() {
		return fmt.Errorf("start_date must not be empty")
	}
	if input.EndDate != nil && input.EndDate.Before(input.StartDate) {
		return fmt.Errorf("end_date must be >= start_date")
	}
	return nil
}

func ValidateUpdateSubscription(input dto.UpdateSubscriptionRequest) error {
	if input.ServiceName == nil &&
		input.Price == nil &&
		input.UserID == nil &&
		input.StartDate == nil &&
		input.EndDate == nil {
		return fmt.Errorf("at least one field must be provided")
	}

	if input.ServiceName != nil &&
		strings.TrimSpace(*input.ServiceName) == "" {
		return fmt.Errorf("service_name must not be empty")
	}

	if input.Price != nil &&
		*input.Price <= 0 {
		return fmt.Errorf("price must be > 0")
	}

	if input.UserID != nil &&
		*input.UserID == uuid.Nil {
		return fmt.Errorf("user_id must not be empty")
	}

	if input.StartDate != nil &&
		input.StartDate.IsZero() {
		return fmt.Errorf("start_date must not be empty")
	}

	if input.EndDate != nil &&
		input.EndDate.IsZero() {
		return fmt.Errorf("end_date must not be empty")
	}

	if input.StartDate != nil &&
		input.EndDate != nil &&
		input.EndDate.Before(*input.StartDate) {
		return fmt.Errorf("end_date must be >= start_date")
	}

	return nil
}

func WriteDomainError(w http.ResponseWriter, logger *slog.Logger, err error, details any) {
	type errData struct {
		Code int
		Msg  string
	}

	errorMap := map[error]errData{
		domain.ErrConflict:         {http.StatusConflict, "conflict"},
		domain.ErrAlreadyExists:    {http.StatusConflict, "already exists"},
		domain.ErrNotFound:         {http.StatusNotFound, "not found"},
		domain.ErrNoFieldsToUpdate: {http.StatusBadRequest, "no fields to update"},
	}

	for domainErr, data := range errorMap {
		if errors.Is(err, domainErr) {
			WriteError(w, logger, data.Code, data.Msg)
			return
		}
	}

	// default
	WriteInternalError(w, logger, err, details)
}
