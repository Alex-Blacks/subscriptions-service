package test

import (
	"context"
	"time"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/Alex-Blacks/subscriptions/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockTx struct{}

func NewMockTx() *MockTx {
	return &MockTx{}
}

func (m *MockTx) Commit(ctx context.Context) error {
	return nil
}

func (m *MockTx) Rollback(ctx context.Context) error {
	return nil
}

func (m *MockTx) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	panic("unexpected call")
}

func (m *MockTx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (pgx.Rows, error) {
	panic("unexpected call")
}

func (m *MockTx) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) pgx.Row {
	panic("unexpected call")
}

func (ms *MockTx) BeginTx(ctx context.Context) (domain.Tx, error) {
	return &MockTx{}, nil
}

type MockStorage struct {
	data   map[int]domain.Subscription
	nextID int
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data:   make(map[int]domain.Subscription),
		nextID: 1,
	}
}

func (ms *MockStorage) CreateSubscription(
	ctx context.Context,
	q domain.Querier,
	input domain.CreateSubscriptionInput,
) (int, error) {

	for _, sub := range ms.data {
		if sub.UserID == input.UserID &&
			sub.ServiceName == input.ServiceName {

			return 0, domain.ErrAlreadyExists
		}
	}

	id := ms.nextID
	ms.nextID++

	sub := domain.Subscription{
		ID:          id,
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		CreatedAt:   time.Now(),
	}

	ms.data[id] = sub

	return id, nil
}

func (ms *MockStorage) GetSubscriptionByID(
	ctx context.Context,
	q domain.Querier,
	id int,
) (domain.Subscription, error) {

	sub, ok := ms.data[id]
	if !ok {
		return domain.Subscription{}, domain.ErrNotFound
	}

	return sub, nil
}

func (ms *MockStorage) DeleteSubscription(
	ctx context.Context,
	q domain.Querier,
	id int,
) error {

	_, ok := ms.data[id]
	if !ok {
		return domain.ErrNotFound
	}

	delete(ms.data, id)

	return nil
}

func (ms *MockStorage) ListSubscription(
	ctx context.Context,
	q domain.Querier,
	filter domain.ListFilter,
) ([]domain.Subscription, error) {

	result := make([]domain.Subscription, 0)

	for _, sub := range ms.data {
		result = append(result, sub)
	}

	return result, nil
}

func (ms *MockStorage) UpdateSubscription(
	ctx context.Context,
	q domain.Querier,
	id int,
	update domain.UpdateSubscriptionInput,
) (domain.Subscription, error) {

	sub, ok := ms.data[id]
	if !ok {
		return domain.Subscription{}, domain.ErrNotFound
	}

	if update.ServiceName != nil {
		sub.ServiceName = *update.ServiceName
	}

	if update.Price != nil {
		sub.Price = *update.Price
	}

	if update.UserID != nil {
		sub.UserID = *update.UserID
	}

	if update.StartDate != nil {
		sub.StartDate = *update.StartDate
	}

	if update.EndDate != nil {
		sub.EndDate = update.EndDate
	}

	ms.data[id] = sub

	return sub, nil
}

func (ms *MockStorage) SumSubscriptionPrice(
	ctx context.Context,
	q domain.Querier,
	filter domain.SumFilter,
) (int, error) {

	sum := 0

	for _, sub := range ms.data {
		sum += sub.Price
	}

	return sum, nil
}

func NewTestService(seed map[int]domain.Subscription) (*service.Service, *MockStorage) {
	storageTx := NewMockTx()
	storageRepo := NewMockStorage()
	storageRepo.data = seed
	storageRepo.nextID = 1

	svc := service.NewService(storageTx, storageRepo)

	return svc, storageRepo
}
