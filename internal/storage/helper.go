package storage

import (
	"fmt"
	"strings"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
)

func CheckSumFilter(filter domain.SumFilter) ([]string, []any, int) {
	var (
		whereParts []string
		args       []any
		argPos     = 1
	)
	if filter.UserID != nil {
		whereParts = append(whereParts, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, *filter.UserID)
		argPos++
	}
	if filter.ServiceName != nil {
		whereParts = append(whereParts, fmt.Sprintf("service_name = $%d", argPos))
		args = append(args, *filter.ServiceName)
		argPos++
	}
	if filter.From != nil {
		whereParts = append(whereParts, fmt.Sprintf("start_date >= $%d", argPos))
		args = append(args, *filter.From)
		argPos++
	}
	if filter.To != nil {
		whereParts = append(whereParts, fmt.Sprintf("start_date <= $%d", argPos))
		args = append(args, *filter.To)
		argPos++
	}
	return whereParts, args, argPos
}

func CheckSListFilter(filter domain.ListFilter) (string, []any) {
	var (
		whereParts []string
		args       []any
		argPos     = 1
	)
	if filter.UserID != nil {
		whereParts = append(whereParts, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, *filter.UserID)
		argPos++
	}
	if filter.ServiceName != nil {
		whereParts = append(whereParts, fmt.Sprintf("service_name = $%d", argPos))
		args = append(args, *filter.ServiceName)
		argPos++
	}
	if filter.From != nil {
		whereParts = append(whereParts, fmt.Sprintf("start_date >= $%d", argPos))
		args = append(args, *filter.From)
		argPos++
	}
	if filter.To != nil {
		whereParts = append(whereParts, fmt.Sprintf("start_date <= $%d", argPos))
		args = append(args, *filter.To)
		argPos++
	}

	query := `
		SELECT id,service_name,price,user_id,start_date,end_date,created_at
		FROM subscriptions
	`
	if len(whereParts) != 0 {
		query += " WHERE " + strings.Join(whereParts, " AND ")
	}

	query += " ORDER BY start_date DESC, id DESC"

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}

	query += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, filter.Limit)
	argPos++

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, filter.Offset)
		argPos++
	}

	return query, args
}
