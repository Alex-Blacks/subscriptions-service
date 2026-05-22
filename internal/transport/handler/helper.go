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
