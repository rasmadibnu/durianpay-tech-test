package entity

type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewPaginatedResult[T any](data []T, page, limit, total int) PaginatedResult[T] {
	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}
	if data == nil {
		data = []T{}
	}
	return PaginatedResult[T]{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
