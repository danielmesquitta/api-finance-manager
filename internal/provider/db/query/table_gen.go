// Code generated by prisma-to-go. DO NOT EDIT.

package query

import "fmt"

type TableAccount string

func (t TableAccount) String() string {
	return string(t)
}

func (t TableAccount) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t TableAccount) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableAccount) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

func (t TableAccount) ColumnInstitutionID() string {
	return fmt.Sprintf("%s.institution_id", t)
}

func (t TableAccount) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableAccount) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TableAccount) ColumnType() string {
	return fmt.Sprintf("%s.type", t)
}

func (t TableAccount) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableAccount) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

const tableAccount = TableAccount("accounts")

type TableAccountBalance string

func (t TableAccountBalance) String() string {
	return string(t)
}

func (t TableAccountBalance) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableAccountBalance) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t TableAccountBalance) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableAccountBalance) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableAccountBalance) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableAccountBalance) ColumnAccountID() string {
	return fmt.Sprintf("%s.account_id", t)
}

func (t TableAccountBalance) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

const tableAccountBalance = TableAccountBalance("account_balances")

type TableBudget string

func (t TableBudget) String() string {
	return string(t)
}

func (t TableBudget) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableBudget) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableBudget) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

func (t TableBudget) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableBudget) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t TableBudget) ColumnDate() string {
	return fmt.Sprintf("%s.date", t)
}

func (t TableBudget) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

const tableBudget = TableBudget("budgets")

type TableBudgetCategory string

func (t TableBudgetCategory) String() string {
	return string(t)
}

func (t TableBudgetCategory) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableBudgetCategory) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t TableBudgetCategory) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableBudgetCategory) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableBudgetCategory) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableBudgetCategory) ColumnBudgetID() string {
	return fmt.Sprintf("%s.budget_id", t)
}

func (t TableBudgetCategory) ColumnCategoryID() string {
	return fmt.Sprintf("%s.category_id", t)
}

const tableBudgetCategory = TableBudgetCategory("budget_categories")

type TableInstitution string

func (t TableInstitution) String() string {
	return string(t)
}

func (t TableInstitution) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableInstitution) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableInstitution) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableInstitution) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableInstitution) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t TableInstitution) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TableInstitution) ColumnLogo() string {
	return fmt.Sprintf("%s.logo", t)
}

const tableInstitution = TableInstitution("institutions")

type TableInvestment string

func (t TableInvestment) String() string {
	return string(t)
}

func (t TableInvestment) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableInvestment) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TableInvestment) ColumnRateType() string {
	return fmt.Sprintf("%s.rateType", t)
}

func (t TableInvestment) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableInvestment) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableInvestment) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableInvestment) ColumnCategoryID() string {
	return fmt.Sprintf("%s.category_id", t)
}

func (t TableInvestment) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t TableInvestment) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t TableInvestment) ColumnRate() string {
	return fmt.Sprintf("%s.rate", t)
}

func (t TableInvestment) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

const tableInvestment = TableInvestment("investments")

type TableInvestmentCategory string

func (t TableInvestmentCategory) String() string {
	return string(t)
}

func (t TableInvestmentCategory) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TableInvestmentCategory) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableInvestmentCategory) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableInvestmentCategory) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableInvestmentCategory) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableInvestmentCategory) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

const tableInvestmentCategory = TableInvestmentCategory("investment_categories")

type TablePaymentMethod string

func (t TablePaymentMethod) String() string {
	return string(t)
}

func (t TablePaymentMethod) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TablePaymentMethod) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t TablePaymentMethod) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TablePaymentMethod) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TablePaymentMethod) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TablePaymentMethod) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

const tablePaymentMethod = TablePaymentMethod("payment_methods")

type TableTransaction string

func (t TableTransaction) String() string {
	return string(t)
}

func (t TableTransaction) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableTransaction) ColumnIsIgnored() string {
	return fmt.Sprintf("%s.is_ignored", t)
}

func (t TableTransaction) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableTransaction) ColumnAccountID() string {
	return fmt.Sprintf("%s.account_id", t)
}

func (t TableTransaction) ColumnPaymentMethodID() string {
	return fmt.Sprintf("%s.payment_method_id", t)
}

func (t TableTransaction) ColumnInstitutionID() string {
	return fmt.Sprintf("%s.institution_id", t)
}

func (t TableTransaction) ColumnCategoryID() string {
	return fmt.Sprintf("%s.category_id", t)
}

func (t TableTransaction) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t TableTransaction) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TableTransaction) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableTransaction) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableTransaction) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t TableTransaction) ColumnDate() string {
	return fmt.Sprintf("%s.date", t)
}

func (t TableTransaction) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

const tableTransaction = TableTransaction("transactions")

type TableTransactionCategory string

func (t TableTransactionCategory) String() string {
	return string(t)
}

func (t TableTransactionCategory) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableTransactionCategory) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t TableTransactionCategory) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TableTransactionCategory) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableTransactionCategory) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableTransactionCategory) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

const tableTransactionCategory = TableTransactionCategory("transaction_categories")

type TableUser string

func (t TableUser) String() string {
	return string(t)
}

func (t TableUser) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t TableUser) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t TableUser) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t TableUser) ColumnAvatar() string {
	return fmt.Sprintf("%s.avatar", t)
}

func (t TableUser) ColumnSynchronizedAt() string {
	return fmt.Sprintf("%s.synchronized_at", t)
}

func (t TableUser) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t TableUser) ColumnOpenFinanceID() string {
	return fmt.Sprintf("%s.open_finance_id", t)
}

func (t TableUser) ColumnEmail() string {
	return fmt.Sprintf("%s.email", t)
}

func (t TableUser) ColumnVerifiedEmail() string {
	return fmt.Sprintf("%s.verified_email", t)
}

func (t TableUser) ColumnProvider() string {
	return fmt.Sprintf("%s.provider", t)
}

func (t TableUser) ColumnTier() string {
	return fmt.Sprintf("%s.tier", t)
}

func (t TableUser) ColumnSubscriptionExpiresAt() string {
	return fmt.Sprintf("%s.subscription_expires_at", t)
}

func (t TableUser) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t TableUser) ColumnAuthID() string {
	return fmt.Sprintf("%s.auth_id", t)
}

const tableUser = TableUser("users")
