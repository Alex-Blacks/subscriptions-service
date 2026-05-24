package test

import (
	"context"
	"errors"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/service"
)

type MockService struct {
	CreateFn func(ctx context.Context, input domain.CreateSubscriptionInput) (int, error)
	GetFn    func(ctx context.Context, id int) (domain.Subscription, error)
	DeleteFn func(ctx context.Context, id int) error
	UpdateFn func(ctx context.Context, id int, input domain.UpdateSubscriptionInput) (domain.Subscription, error)
	ListFn   func(ctx context.Context, filter domain.ListFilter) ([]domain.Subscription, error)
	SumFn    func(ctx context.Context, filter domain.SumFilter) (int, error)
}

func (m *MockService) CreateSubscription(ctx context.Context, input domain.CreateSubscriptionInput) (int, error) {
	return m.CreateFn(ctx, input)
}

func (m *MockService) GetSubscriptionByID(ctx context.Context, id int) (domain.Subscription, error) {
	return m.GetFn(ctx, id)
}

func (m *MockService) DeleteSubscription(ctx context.Context, id int) error {
	return m.DeleteFn(ctx, id)
}

func (m *MockService) UpdateSubscription(ctx context.Context, id int, input domain.UpdateSubscriptionInput) (domain.Subscription, error) {
	return m.UpdateFn(ctx, id, input)
}

func (m *MockService) ListSubscription(ctx context.Context, filter domain.ListFilter) ([]domain.Subscription, error) {
	return m.ListFn(ctx, filter)
}

func (m *MockService) SumSubscriptionPrice(ctx context.Context, filter domain.SumFilter) (int, error) {
	return m.SumFn(ctx, filter)
}

var _ service.SubscriptionService = (*MockService)(nil)

var ErrMock = errors.New("mock error")
