package storage

import (
	"fmt"

	"github.com/Alex-Blacks/subscriptions/internal/domain"
)

func CheckFilter(filter domain.ListFilter) ([]string, []any, int) {
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
