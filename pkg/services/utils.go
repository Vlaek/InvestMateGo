package services

import (
	"context"
)

// Функция для получения данных с пагинацией
func GetWithPagination[T any](
	ctx context.Context,
	getFunc func(context.Context) ([]T, error),
	page, limit int,
) ([]T, int64, error) {
	items, err := getFunc(ctx)

	if err != nil {
		return nil, 0, err
	}

	result, total := paginate(items, page, limit)

	return result, total, nil
}

// Функция пагинации
func paginate[T any](items []T, page, limit int) ([]T, int64) {
	if len(items) == 0 {
		return []T{}, 0
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = len(items)
	}

	start := (page - 1) * limit

	if start >= len(items) {
		return []T{}, int64(len(items))
	}

	end := start + limit

	if end > len(items) {
		end = len(items)
	}

	return items[start:end], int64(len(items))
}
