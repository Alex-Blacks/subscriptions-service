package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateSubscription(ctx context.Context, q domain.Querier, input domain.CreateSubscriptionInput) (int, error) {
	var id int
	if err := q.QueryRow(ctx, `
		INSERT INTO subscriptions(
			service_name,
			price,
			user_id,
			start_date,
			end_date
		) VALUES ($1,$2,$3,$4,$5) RETURNING id
		`,
		input.ServiceName,
		input.Price,
		input.UserID,
		input.StartDate,
		input.EndDate,
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return 0, domain.ErrAlreadyExists
		}
		return 0, fmt.Errorf("create store: %w", err)
	}
	return id, nil
}

func (s *Storage) GetSubscriptionByID(ctx context.Context, q domain.Querier, id int) (domain.Subscription, error) {
	var subscription domain.Subscription
	if err := q.QueryRow(ctx, `
		SELECT id,service_name,price,user_id,start_date,end_date,created_at
		FROM subscriptions
		WHERE id = $1	
	`, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrNotFound
		}
		return domain.Subscription{}, fmt.Errorf("get subscription by id: %w", err)
	}
	return subscription, nil
}

func (s *Storage) DeleteSubscription(ctx context.Context, q domain.Querier, id int) error {
	tag, err := q.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Storage) UpdateSubscription(ctx context.Context, q domain.Querier, id int, update domain.UpdateSubscriptionInput) (domain.Subscription, error) {
	var (
		setParts []string
		args     []any
		argPos   = 1
	)
	if update.ServiceName != nil {
		setParts = append(setParts, fmt.Sprintf("service_name = $%d", argPos))
		args = append(args, *update.ServiceName)
		argPos++
	}
	if update.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argPos))
		args = append(args, *update.Price)
		argPos++
	}
	if update.UserID != nil {
		setParts = append(setParts, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, *update.UserID)
		argPos++
	}
	if update.StartDate != nil {
		setParts = append(setParts, fmt.Sprintf("start_date = $%d", argPos))
		args = append(args, *update.StartDate)
		argPos++
	}
	if update.EndDate != nil {
		setParts = append(setParts, fmt.Sprintf("end_date = $%d", argPos))
		args = append(args, *update.EndDate)
		argPos++
	}

	if len(setParts) == 0 {
		return domain.Subscription{}, domain.ErrNoFieldsToUpdate
	}

	query := fmt.Sprintf(`
		UPDATE subscriptions
		SET %s 
		WHERE id = $%d
		RETURNING id,service_name,price,user_id,start_date,end_date,created_at
	`, strings.Join(setParts, ", "), argPos)

	args = append(args, id)

	var sub domain.Subscription

	if err := q.QueryRow(ctx, query, args...).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrNotFound
		}
		return domain.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}
	return sub, nil
}

func (s *Storage) ListSubscription(ctx context.Context, q domain.Querier, filter domain.ListFilter) ([]domain.Subscription, error) {
	query, args := CheckSListFilter(filter)

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query list subscription: %w", err)
	}
	defer rows.Close()

	var list []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription

		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan list subscription: %w", err)
		}

		list = append(list, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration failed: %w", err)
	}
	return list, nil

}

func (s *Storage) SumSubscriptionPrice(ctx context.Context, q domain.Querier, filter domain.SumFilter) (int, error) {
	var total int
	whereParts, args, _ := CheckSumFilter(filter)
	query := `
		SELECT COALESCE(SUM(price),0)
		FROM subscriptions
	`
	if len(whereParts) != 0 {
		query += " WHERE " + strings.Join(whereParts, " AND ")
	}

	if err := q.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("query sum failed: %w", err)
	}
	return total, nil
}
