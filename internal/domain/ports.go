package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage interface {
	Querier
	BeginTx(ctx context.Context) (Tx, error)
}

type Tx interface {
	Querier
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type CreateSubscriptionInput struct {
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

type UpdateSubscriptionInput struct {
	ServiceName *string
	Price       *int
	UserID      *uuid.UUID
	StartDate   *time.Time
	EndDate     *time.Time
}

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
}

type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        *time.Time
	To          *time.Time
	Limit       int
	Offset      int
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, q Querier, input CreateSubscriptionInput) (int, error)
	GetSubscriptionByID(ctx context.Context, q Querier, id int) (Subscription, error)
	DeleteSubscription(ctx context.Context, q Querier, id int) error
	UpdateSubscription(ctx context.Context, q Querier, id int, update UpdateSubscriptionInput) (Subscription, error)
	ListSubscription(ctx context.Context, q Querier, filter ListFilter) ([]Subscription, error)
	SumSubscriptionPrice(ctx context.Context, q Querier, filter ListFilter) (int, error)
}
