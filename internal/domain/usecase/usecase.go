package usecase

import (
	"math"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type PaginationInput struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
}

func PreparePaginationInput(in *PaginationInput) (offset uint) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 {
		in.PageSize = 20
	}
	return (in.Page - 1) * in.PageSize
}

func PreparePaginationOutput[T any](
	out *entity.PaginatedList[T],
	in PaginationInput,
	count int64,
) {
	out.Page = in.Page
	out.PageSize = in.PageSize
	out.TotalItems = uint(count)
	out.TotalPages = uint(math.Ceil(float64(count) / float64(in.PageSize)))
}
