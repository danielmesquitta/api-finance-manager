package pgrepo

type Table string
type Column string

const (
	TableInvestment     Table = "investments"
	TableAccount        Table = "accounts"
	TableBudget         Table = "budgets"
	TableBudgetCategory Table = "budget_categories"
	TableUser           Table = "users"
	TableTransaction    Table = "transactions"
	TableCategory       Table = "categories"
	TableCreditCard     Table = "credit_cards"
	TableInstitution    Table = "institutions"
)

// Columns for table Transaction
const (
	ColumnTransactionDate          Column = "date"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionID            Column = "id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionUpdatedAt     Column = "updated_at"
)

// Columns for table Investment
const (
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
)

// Columns for table BudgetCategory
const (
	ColumnBudgetCategoryCategoryID Column = "category_id"
	ColumnBudgetCategoryID         Column = "id"
	ColumnBudgetCategoryAmount     Column = "amount"
	ColumnBudgetCategoryCreatedAt  Column = "created_at"
	ColumnBudgetCategoryUpdatedAt  Column = "updated_at"
	ColumnBudgetCategoryBudgetID   Column = "budget_id"
)

// Columns for table Budget
const (
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
	ColumnBudgetID        Column = "id"
	ColumnBudgetAmount    Column = "amount"
	ColumnBudgetDate      Column = "date"
	ColumnBudgetCreatedAt Column = "created_at"
)

// Columns for table User
const (
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserUpdatedAt             Column = "updated_at"
	ColumnUserID                    Column = "id"
	ColumnUserName                  Column = "name"
	ColumnUserEmail                 Column = "email"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
	ColumnUserTier                  Column = "tier"
	ColumnUserAvatar                Column = "avatar"
)

// Columns for table Category
const (
	ColumnCategoryUpdatedAt  Column = "updated_at"
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
)

// Columns for table Account
const (
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountID            Column = "id"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountType          Column = "type"
	ColumnAccountCreatedAt     Column = "created_at"
)

// Columns for table CreditCard
const (
	ColumnCreditCardCreatedAt      Column = "created_at"
	ColumnCreditCardUpdatedAt      Column = "updated_at"
	ColumnCreditCardAccountID      Column = "account_id"
	ColumnCreditCardID             Column = "id"
	ColumnCreditCardLevel          Column = "level"
	ColumnCreditCardBrand          Column = "brand"
	ColumnCreditCardLimit          Column = "limit"
	ColumnCreditCardAvailableLimit Column = "available_limit"
)

// Columns for table Institution
const (
	ColumnInstitutionExternalID Column = "external_id"
	ColumnInstitutionName       Column = "name"
	ColumnInstitutionLogo       Column = "logo"
	ColumnInstitutionCreatedAt  Column = "created_at"
	ColumnInstitutionUpdatedAt  Column = "updated_at"
	ColumnInstitutionID         Column = "id"
)
