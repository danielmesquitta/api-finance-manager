package pgrepo

type Table string
type Column string

const (
	TableTransaction    Table = "transactions"
	TableCreditCard     Table = "credit_cards"
	TableInstitution    Table = "institutions"
	TableAccount        Table = "accounts"
	TableBudget         Table = "budgets"
	TableBudgetCategory Table = "budget_categories"
	TableUser           Table = "users"
	TableInvestment     Table = "investments"
	TableCategory       Table = "categories"
)

// Columns for table Institution
const (
	ColumnInstitutionID         Column = "id"
	ColumnInstitutionExternalID Column = "external_id"
	ColumnInstitutionName       Column = "name"
	ColumnInstitutionLogo       Column = "logo"
	ColumnInstitutionCreatedAt  Column = "created_at"
	ColumnInstitutionUpdatedAt  Column = "updated_at"
)

// Columns for table Budget
const (
	ColumnBudgetAmount    Column = "amount"
	ColumnBudgetCreatedAt Column = "created_at"
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
	ColumnBudgetID        Column = "id"
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

// Columns for table Transaction
const (
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionID            Column = "id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionDate          Column = "date"
)

// Columns for table Category
const (
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
)

// Columns for table Account
const (
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountType          Column = "type"
	ColumnAccountID            Column = "id"
	ColumnAccountExternalID    Column = "external_id"
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

// Columns for table User
const (
	ColumnUserTier                  Column = "tier"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserUpdatedAt             Column = "updated_at"
	ColumnUserID                    Column = "id"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserName                  Column = "name"
	ColumnUserEmail                 Column = "email"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
)

// Columns for table Investment
const (
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
)
