package usecase

import (
	"math"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type PaginationInput struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
}

func preparePaginationInput(in *PaginationInput) (offset uint) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 {
		in.PageSize = 20
	}
	return (in.Page - 1) * in.PageSize
}

func preparePaginationOutput[T any](
	out *entity.PaginatedList[T],
	in PaginationInput,
	count int64,
) {
	out.Page = in.Page
	out.PageSize = in.PageSize
	out.TotalItems = uint(count)
	out.TotalPages = uint(math.Ceil(float64(count) / float64(in.PageSize)))
}

func prepareTransactionOptions(
	in repo.TransactionOptions,
) []repo.TransactionOption {
	opts := []repo.TransactionOption{}

	if in.Search != "" {
		opts = append(opts, repo.WithTransactionSearch(in.Search))
	}

	if len(in.CategoryIDs) > 0 {
		opts = append(
			opts,
			repo.WithTransactionCategories(in.CategoryIDs...),
		)
	}

	if len(in.InstitutionIDs) > 0 {
		opts = append(
			opts,
			repo.WithTransactionInstitutions(in.InstitutionIDs...),
		)
	}

	if len(in.PaymentMethodIDs) > 0 {
		opts = append(
			opts,
			repo.WithTransactionPaymentMethods(in.PaymentMethodIDs...),
		)
	}

	if !in.StartDate.IsZero() {
		opts = append(
			opts,
			repo.WithTransactionDateAfter(in.StartDate),
		)
	}

	if !in.EndDate.IsZero() {
		opts = append(
			opts,
			repo.WithTransactionDateBefore(in.EndDate),
		)
	}

	if in.IsExpense {
		opts = append(
			opts,
			repo.WithTransactionIsExpense(in.IsExpense),
		)
	}

	if in.IsIncome {
		opts = append(
			opts,
			repo.WithTransactionIsIncome(in.IsIncome),
		)
	}

	if in.IsIgnored != nil {
		opts = append(
			opts,
			repo.WithTransactionIsIgnored(*in.IsIgnored),
		)
	}

	return opts
}
