package handlers

import (
	"math"
	"oriongo/internal/infrastructure"
)

type (
	Response struct {
		Data            interface{} `json:"data"`
		ResponseMessage string      `json:"response_message"`
		Errors          []string    `json:"errors"`
		Success         bool        `json:"success"`
	}

	PaginatedData[T any] struct {
		Items []T `json:"items"`
	}
	PaginatedMeta struct {
		HasNextPage     bool  `json:"has_next_page"`
		HasPreviousPage bool  `json:"has_previous_page"`
		TotalItems      int64 `json:"total_items"`
		TotalPages      int   `json:"total_pages"`
		CurrentPage     int   `json:"current_page"`
		PageSize        int   `json:"page_size"`
	}
	PaginatedResponse[T any] struct {
		Data PaginatedData[T] `json:"data"`
		Meta PaginatedMeta    `json:"meta"`
	}

	BaseHandler struct {
		_dbContext infrastructure.DbContext
	}
)

func NewBaseHandler(_dbContext infrastructure.DbContext) *BaseHandler {
	return &BaseHandler{
		_dbContext: _dbContext,
	}
}
func MakePagination[T any](items []T, count int64, page, perPage int) PaginatedResponse[T] {
	// 10, 2, 30
	totalPages := 0
	if count > 0 {
		totalPages = int(math.Ceil(float64(count) / float64(perPage)))
	}

	hasNextPage := totalPages > 1 && count > 0 && page < totalPages
	hasPreviousPage := count > 0 && page > 1 && page >= totalPages

	return PaginatedResponse[T]{
		Data: PaginatedData[T]{
			Items: items,
		},
		Meta: PaginatedMeta{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			TotalItems:      count,
			TotalPages:      totalPages,
			CurrentPage:     page,
			PageSize:        perPage,
		},
	}
}

func SuccessResponse(data interface{}, message string) *Response {
	if message == "" {
		message = "Request Successful"
	}
	return &Response{
		Success:         true,
		Data:            data,
		ResponseMessage: message,
		Errors:          []string{},
	}
}

func ErrorResponse(err string) *Response {
	return &Response{
		Success:         false,
		Errors:          []string{err},
		ResponseMessage: "Request failed",
	}
}
