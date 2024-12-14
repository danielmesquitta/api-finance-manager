package entity

type PaginatedList[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"total_items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}
