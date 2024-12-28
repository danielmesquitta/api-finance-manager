package query

type Table string
type Column string

const (
	TableCategory       Table = "categories"
	TableAccount        Table = "accounts"
	TableCreditCard     Table = "credit_cards"
	TableInstitution    Table = "institutions"
	TableBudgetCategory Table = "budget_categories"
	TableTransaction    Table = "transactions"
	TableInvestment     Table = "investments"
	TableBudget         Table = "budgets"
	TableUser           Table = "users"
)

// Columns for table User
const (
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserUpdatedAt             Column = "updated_at"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserTier                  Column = "tier"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserID                    Column = "id"
	ColumnUserProvider              Column = "provider"
	ColumnUserName                  Column = "name"
	ColumnUserEmail                 Column = "email"
)

// Columns for table Transaction
const (
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionCategoryID    Column = "category_id"
)

// Columns for table CreditCard
const (
	ColumnCreditCardLevel          Column = "level"
	ColumnCreditCardBrand          Column = "brand"
	ColumnCreditCardLimit          Column = "limit"
	ColumnCreditCardAvailableLimit Column = "available_limit"
	ColumnCreditCardCreatedAt      Column = "created_at"
	ColumnCreditCardUpdatedAt      Column = "updated_at"
	ColumnCreditCardAccountID      Column = "account_id"
	ColumnCreditCardID             Column = "id"
)

// Columns for table Budget
const (
	ColumnBudgetID        Column = "id"
	ColumnBudgetAmount    Column = "amount"
	ColumnBudgetDate      Column = "date"
	ColumnBudgetCreatedAt Column = "created_at"
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
)

// Columns for table BudgetCategory
const (
	ColumnBudgetCategoryID         Column = "id"
	ColumnBudgetCategoryAmount     Column = "amount"
	ColumnBudgetCategoryCreatedAt  Column = "created_at"
	ColumnBudgetCategoryUpdatedAt  Column = "updated_at"
	ColumnBudgetCategoryBudgetID   Column = "budget_id"
	ColumnBudgetCategoryCategoryID Column = "category_id"
)

// Columns for table Investment
const (
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentName       Column = "name"
)

// Columns for table Category
const (
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
)

// Columns for table Account
const (
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountID            Column = "id"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountType          Column = "type"
	ColumnAccountInstitutionID Column = "institution_id"
)

// Columns for table Institution
const (
	ColumnInstitutionName       Column = "name"
	ColumnInstitutionLogo       Column = "logo"
	ColumnInstitutionCreatedAt  Column = "created_at"
	ColumnInstitutionUpdatedAt  Column = "updated_at"
	ColumnInstitutionID         Column = "id"
	ColumnInstitutionExternalID Column = "external_id"
)
