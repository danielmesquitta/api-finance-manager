package transaction

import "github.com/danielmesquitta/api-finance-manager/internal/provider/repo"

func PrepareTransactionOptions(
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
