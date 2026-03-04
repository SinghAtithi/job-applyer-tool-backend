package dto

import (
	"strings"
	"time"
)

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginationParams represents pagination request parameters
type PaginationParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data      interface{} `json:"data"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	Total     int64       `json:"total"`
	TotalPage int64       `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, page, pageSize int, total int64) *PaginatedResponse {
	if pageSize <= 0 {
		pageSize = 1
	}
	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Data:      data,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPages,
	}
}

// SuccessResponse creates a successful response
func SuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code, message, details string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// TimestampFields represents common timestamp fields
type TimestampFields struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// QueryParams represents common query parameters for filtering, sorting, and pagination
type QueryParams struct {
	// Pagination
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	Offset   int `form:"offset" binding:"omitempty,min=0"`
	Limit    int `form:"limit" binding:"omitempty,min=1,max=100"`

	// Sorting - format: field_name or -field_name for descending
	// Examples: sort=created_at or sort=-created_at
	Sort string `form:"sort" binding:"-"`

	// Field selection - comma-separated list of fields to return
	// Example: fields=id,name,email
	Fields string `form:"fields" binding:"-"`

	// Search - generic search query
	Search string `form:"search" binding:"-"`
}

// GetOffset calculates the offset for pagination
// If Offset is explicitly provided, it takes precedence
// Otherwise, calculates from Page and PageSize
func (q *QueryParams) GetOffset() int {
	if q.Offset > 0 {
		return q.Offset
	}
	if q.Page > 0 && q.PageSize > 0 {
		return (q.Page - 1) * q.PageSize
	}
	return 0
}

// GetLimit returns the limit for query
// If Limit is explicitly provided, it takes precedence
// Otherwise, uses PageSize
func (q *QueryParams) GetLimit() int {
	if q.Limit > 0 {
		return q.Limit
	}
	if q.PageSize > 0 {
		return q.PageSize
	}
	return 10 // default limit
}

// GetSortFieldAndOrder parses the sort parameter into field and order
// Returns field name and true for descending, false for ascending
func (q *QueryParams) GetSortFieldAndOrder() (string, bool) {
	if q.Sort == "" {
		return "", false
	}

	// Check if descending (starts with -)
	if q.Sort[0] == '-' {
		return q.Sort[1:], true
	}
	return q.Sort, false
}

// GetFieldsSlice parses the fields parameter into a slice
func (q *QueryParams) GetFieldsSlice() []string {
	if q.Fields == "" {
		return nil
	}

	fields := []string{}
	for _, f := range strings.Split(q.Fields, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			fields = append(fields, f)
		}
	}
	return fields
}
