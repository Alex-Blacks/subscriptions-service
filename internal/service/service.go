package service

import (
	"context"
	"fmt"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/storage"
)

type Service struct {
	storage storage.Storage
	sub     domain.SubscriptionRepository
}

func NewService(st storage.Storage, sub domain.SubscriptionRepository) *Service {
	return &Service{
		storage: st,
		sub:     sub,
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
				err = fmt.Errorf("tx err: %w, rollback err: %w", err, rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(ctx); commitErr != nil {
			err = fmt.Errorf("commit err: %w", commitErr)
		}
	}()

	return fn(tx)
}

func (s *Service) CreateSubscription(ctx context.Context, input domain.CreateSubscriptionInput) (int, error) {
	var id int
	if err := s.WithTx(ctx, func(q domain.Querier) error {
		var err error
		id, err = s.sub.CreateSubscription(ctx, q, input)
		return err
	}); err != nil {
		return 0, err
	}
	return id, nil
}
func (s *Service) GetSubscriptionByID(ctx context.Context, id int) (domain.Subscription, error)
func (s *Service) DeleteSubscription(ctx context.Context, id int) error
func (s *Service) UpdateSubscription(ctx context.Context, id int, update domain.UpdateSubscriptionInput) (domain.Subscription, error)
func (s *Service) ListSubscription(ctx context.Context, filter domain.ListFilter) ([]domain.Subscription, error)
func (s *Service) SumSubscriptionPrice(ctx context.Context, filter domain.ListFilter) (int, error)
