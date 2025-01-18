// Code generated by prisma-to-go. DO NOT EDIT.

package query

type Table = string

const (
	TableAccount        Table = "accounts"
	TableBudget         Table = "budgets"
	TableBudgetCategory Table = "budget_categories"
	TableCategory       Table = "categories"
	TableInstitution    Table = "institutions"
	TableInvestment     Table = "investments"
	TablePaymentMethod  Table = "payment_methods"
	TableTransaction    Table = "transactions"
	TableUser           Table = "users"
)

type Column = string

// Columns for table Account
const (
	ColumnAccountCreatedAt     Column = "accounts.created_at"
	ColumnAccountDeletedAt     Column = "accounts.deleted_at"
	ColumnAccountExternalID    Column = "accounts.external_id"
	ColumnAccountID            Column = "accounts.id"
	ColumnAccountInstitutionID Column = "accounts.institution_id"
	ColumnAccountName          Column = "accounts.name"
	ColumnAccountType          Column = "accounts.type"
	ColumnAccountUpdatedAt     Column = "accounts.updated_at"
	ColumnAccountUserID        Column = "accounts.user_id"
)

// Columns for table Budget
const (
	ColumnBudgetAmount    Column = "budgets.amount"
	ColumnBudgetCreatedAt Column = "budgets.created_at"
	ColumnBudgetDate      Column = "budgets.date"
	ColumnBudgetDeletedAt Column = "budgets.deleted_at"
	ColumnBudgetID        Column = "budgets.id"
	ColumnBudgetUpdatedAt Column = "budgets.updated_at"
	ColumnBudgetUserID    Column = "budgets.user_id"
)

// Columns for table BudgetCategory
const (
	ColumnBudgetCategoryAmount     Column = "budget_categories.amount"
	ColumnBudgetCategoryBudgetID   Column = "budget_categories.budget_id"
	ColumnBudgetCategoryCategoryID Column = "budget_categories.category_id"
	ColumnBudgetCategoryCreatedAt  Column = "budget_categories.created_at"
	ColumnBudgetCategoryDeletedAt  Column = "budget_categories.deleted_at"
	ColumnBudgetCategoryID         Column = "budget_categories.id"
	ColumnBudgetCategoryUpdatedAt  Column = "budget_categories.updated_at"
)

// Columns for table Category
const (
	ColumnCategoryCreatedAt  Column = "categories.created_at"
	ColumnCategoryDeletedAt  Column = "categories.deleted_at"
	ColumnCategoryExternalID Column = "categories.external_id"
	ColumnCategoryID         Column = "categories.id"
	ColumnCategoryName       Column = "categories.name"
	ColumnCategoryUpdatedAt  Column = "categories.updated_at"
)

// Columns for table Institution
const (
	ColumnInstitutionCreatedAt  Column = "institutions.created_at"
	ColumnInstitutionDeletedAt  Column = "institutions.deleted_at"
	ColumnInstitutionExternalID Column = "institutions.external_id"
	ColumnInstitutionID         Column = "institutions.id"
	ColumnInstitutionLogo       Column = "institutions.logo"
	ColumnInstitutionName       Column = "institutions.name"
	ColumnInstitutionUpdatedAt  Column = "institutions.updated_at"
)

// Columns for table Investment
const (
	ColumnInvestmentAmount     Column = "investments.amount"
	ColumnInvestmentCreatedAt  Column = "investments.created_at"
	ColumnInvestmentDeletedAt  Column = "investments.deleted_at"
	ColumnInvestmentExternalID Column = "investments.external_id"
	ColumnInvestmentID         Column = "investments.id"
	ColumnInvestmentName       Column = "investments.name"
	ColumnInvestmentRate       Column = "investments.rate"
	ColumnInvestmentRateType   Column = "investments.rateType"
	ColumnInvestmentType       Column = "investments.type"
	ColumnInvestmentUpdatedAt  Column = "investments.updated_at"
	ColumnInvestmentUserID     Column = "investments.user_id"
)

// Columns for table PaymentMethod
const (
	ColumnPaymentMethodCreatedAt  Column = "payment_methods.created_at"
	ColumnPaymentMethodDeletedAt  Column = "payment_methods.deleted_at"
	ColumnPaymentMethodExternalID Column = "payment_methods.external_id"
	ColumnPaymentMethodID         Column = "payment_methods.id"
	ColumnPaymentMethodName       Column = "payment_methods.name"
	ColumnPaymentMethodUpdatedAt  Column = "payment_methods.updated_at"
)

// Columns for table Transaction
const (
	ColumnTransactionAccountID       Column = "transactions.account_id"
	ColumnTransactionAmount          Column = "transactions.amount"
	ColumnTransactionCategoryID      Column = "transactions.category_id"
	ColumnTransactionCreatedAt       Column = "transactions.created_at"
	ColumnTransactionDate            Column = "transactions.date"
	ColumnTransactionDeletedAt       Column = "transactions.deleted_at"
	ColumnTransactionExternalID      Column = "transactions.external_id"
	ColumnTransactionID              Column = "transactions.id"
	ColumnTransactionInstitutionID   Column = "transactions.institution_id"
	ColumnTransactionIsIgnored       Column = "transactions.is_ignored"
	ColumnTransactionName            Column = "transactions.name"
	ColumnTransactionPaymentMethodID Column = "transactions.payment_method_id"
	ColumnTransactionUpdatedAt       Column = "transactions.updated_at"
	ColumnTransactionUserID          Column = "transactions.user_id"
)

// Columns for table User
const (
	ColumnUserAvatar                Column = "users.avatar"
	ColumnUserCreatedAt             Column = "users.created_at"
	ColumnUserDeletedAt             Column = "users.deleted_at"
	ColumnUserEmail                 Column = "users.email"
	ColumnUserExternalID            Column = "users.external_id"
	ColumnUserID                    Column = "users.id"
	ColumnUserName                  Column = "users.name"
	ColumnUserProvider              Column = "users.provider"
	ColumnUserSubscriptionExpiresAt Column = "users.subscription_expires_at"
	ColumnUserSynchronizedAt        Column = "users.synchronized_at"
	ColumnUserTier                  Column = "users.tier"
	ColumnUserUpdatedAt             Column = "users.updated_at"
	ColumnUserVerifiedEmail         Column = "users.verified_email"
)
