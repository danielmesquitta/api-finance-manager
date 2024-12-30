// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: copyfrom.go

package sqlc

import (
	"context"
)

// iteratorForCreateAccounts implements pgx.CopyFromSource.
type iteratorForCreateAccounts struct {
	rows                 []CreateAccountsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateAccounts) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateAccounts) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ExternalID,
		r.rows[0].Name,
		r.rows[0].Type,
		r.rows[0].UserID,
		r.rows[0].InstitutionID,
	}, nil
}

func (r iteratorForCreateAccounts) Err() error {
	return nil
}

func (q *Queries) CreateAccounts(ctx context.Context, arg []CreateAccountsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"accounts"}, []string{"external_id", "name", "type", "user_id", "institution_id"}, &iteratorForCreateAccounts{rows: arg})
}

// iteratorForCreateBudgetCategories implements pgx.CopyFromSource.
type iteratorForCreateBudgetCategories struct {
	rows                 []CreateBudgetCategoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateBudgetCategories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateBudgetCategories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Amount,
		r.rows[0].BudgetID,
		r.rows[0].CategoryID,
	}, nil
}

func (r iteratorForCreateBudgetCategories) Err() error {
	return nil
}

func (q *Queries) CreateBudgetCategories(ctx context.Context, arg []CreateBudgetCategoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"budget_categories"}, []string{"amount", "budget_id", "category_id"}, &iteratorForCreateBudgetCategories{rows: arg})
}

// iteratorForCreateCategories implements pgx.CopyFromSource.
type iteratorForCreateCategories struct {
	rows                 []CreateCategoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateCategories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateCategories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ExternalID,
		r.rows[0].Name,
	}, nil
}

func (r iteratorForCreateCategories) Err() error {
	return nil
}

func (q *Queries) CreateCategories(ctx context.Context, arg []CreateCategoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"categories"}, []string{"external_id", "name"}, &iteratorForCreateCategories{rows: arg})
}

// iteratorForCreateInstitutions implements pgx.CopyFromSource.
type iteratorForCreateInstitutions struct {
	rows                 []CreateInstitutionsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateInstitutions) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateInstitutions) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ExternalID,
		r.rows[0].Name,
		r.rows[0].Logo,
	}, nil
}

func (r iteratorForCreateInstitutions) Err() error {
	return nil
}

func (q *Queries) CreateInstitutions(ctx context.Context, arg []CreateInstitutionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"institutions"}, []string{"external_id", "name", "logo"}, &iteratorForCreateInstitutions{rows: arg})
}

// iteratorForCreateTransactions implements pgx.CopyFromSource.
type iteratorForCreateTransactions struct {
	rows                 []CreateTransactionsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateTransactions) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateTransactions) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ExternalID,
		r.rows[0].Name,
		r.rows[0].Amount,
		r.rows[0].PaymentMethod,
		r.rows[0].Date,
		r.rows[0].UserID,
		r.rows[0].AccountID,
		r.rows[0].InstitutionID,
		r.rows[0].CategoryID,
	}, nil
}

func (r iteratorForCreateTransactions) Err() error {
	return nil
}

func (q *Queries) CreateTransactions(ctx context.Context, arg []CreateTransactionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"transactions"}, []string{"external_id", "name", "amount", "payment_method", "date", "user_id", "account_id", "institution_id", "category_id"}, &iteratorForCreateTransactions{rows: arg})
}
