package usecase

type PaginationInput struct {
	Search   string `json:"search,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

func preparePaginationInput(in *PaginationInput) (offset int) {
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 20
	}
	return (in.Page - 1) * in.PageSize
}
