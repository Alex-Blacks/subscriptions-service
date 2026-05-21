package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
)

type Service struct {
	sub domain.SubscriptionRepository
}

func NewService(sub domain.SubscriptionRepository) *Service {
	return &Service{
		sub: sub,
	}
}

func (s *Service) WithTx(ctx context.Context, fn func(q domain.Querier) error) (err error) {
	tx, err := s.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(ctx); commitErr != nil {
			err = fmt.Errorf("commit err: %w", commitErr)
		}
	}()

	err = fn(tx)
	return err
}

func (s *Service) CreateSubscription(ctx context.Context, input domain.CreateSubscriptionInput) (int, error) {
	var id int
	if err := s.WithTx(ctx, func(q domain.Querier) error {
		var createErr error
		id, createErr = s.sub.CreateSubscription(ctx, q, input)
		return createErr
	}); err != nil {
		return 0, err
	}
	return id, nil
}
func (s *Service) GetSubscriptionByID(ctx context.Context, id int) (domain.Subscription, error) {
	return s.sub.GetSubscriptionByID(ctx, s.storage, id)
}
func (s *Service) DeleteSubscription(ctx context.Context, id int) error {
	return s.WithTx(ctx, func(q domain.Querier) error {
		return s.sub.DeleteSubscription(ctx, q, id)
	})
}
func (s *Service) UpdateSubscription(ctx context.Context, id int, update domain.UpdateSubscriptionInput) (domain.Subscription, error) {
	var sub domain.Subscription
	if err := s.WithTx(ctx, func(q domain.Querier) error {
		var updateErr error
		sub, updateErr = s.sub.UpdateSubscription(ctx, q, id, update)
		return updateErr
	}); err != nil {
		return domain.Subscription{}, err
	}
	return sub, nil
}
func (s *Service) ListSubscription(ctx context.Context, filter domain.ListFilter) ([]domain.Subscription, error) {
	return s.sub.ListSubscription(ctx, s.storage, filter)
}
func (s *Service) SumSubscriptionPrice(ctx context.Context, filter domain.ListFilter) (int, error) {
	return s.sub.SumSubscriptionPrice(ctx, s.storage, filter)
}
