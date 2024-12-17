package entity

type PaginatedList[T any] struct {
	Items      []T `json:"items,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}
