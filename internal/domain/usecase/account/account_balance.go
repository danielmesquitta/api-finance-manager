package account

import "github.com/danielmesquitta/api-finance-manager/internal/provider/repo"

func prepareBalanceOptions(
	in repo.TransactionOptions,
) []repo.AccountBalanceOption {
	balanceOpts := []repo.AccountBalanceOption{}

	if len(in.InstitutionIDs) > 0 {
		balanceOpts = append(
			balanceOpts,
			repo.WithAAccountBalanceInstitutions(in.InstitutionIDs),
		)
	}

	return balanceOpts
}
