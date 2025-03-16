package entity

type PaginatedList[T any] struct {
	Items      []T  `json:"items"`
	TotalItems uint `json:"total_items"`
	Page       uint `json:"page,omitempty"`
	PageSize   uint `json:"page_size,omitempty"`
	TotalPages uint `json:"total_pages,omitempty"`
}
