package entity

type PaginatedList[T any] struct {
	Items      []T  `json:"items"`
	TotalItems uint `json:"total_items"`
	Page       uint `json:"page"`
	PageSize   uint `json:"page_size"`
	TotalPages uint `json:"total_pages"`
}
