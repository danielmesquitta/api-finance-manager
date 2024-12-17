package pgrepo

type Table string
type Column string

const (
	TableUser           Table = "users"
	TableCategory       Table = "categories"
	TableInstitution    Table = "institutions"
	TableTransaction    Table = "transactions"
	TableInvestment     Table = "investments"
	TableAccount        Table = "accounts"
	TableCreditCard     Table = "credit_cards"
	TableBudget         Table = "budgets"
	TableBudgetCategory Table = "budget_categories"
)

// Columns for table Account
const (
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountType          Column = "type"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountID            Column = "id"
)

// Columns for table CreditCard
const (
	ColumnCreditCardBrand          Column = "brand"
	ColumnCreditCardLimit          Column = "limit"
	ColumnCreditCardAvailableLimit Column = "available_limit"
	ColumnCreditCardCreatedAt      Column = "created_at"
	ColumnCreditCardUpdatedAt      Column = "updated_at"
	ColumnCreditCardAccountID      Column = "account_id"
	ColumnCreditCardID             Column = "id"
	ColumnCreditCardLevel          Column = "level"
)

// Columns for table BudgetCategory
const (
	ColumnBudgetCategoryCreatedAt  Column = "created_at"
	ColumnBudgetCategoryUpdatedAt  Column = "updated_at"
	ColumnBudgetCategoryBudgetID   Column = "budget_id"
	ColumnBudgetCategoryCategoryID Column = "category_id"
	ColumnBudgetCategoryID         Column = "id"
	ColumnBudgetCategoryAmount     Column = "amount"
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
	ColumnBudgetID        Column = "id"
	ColumnBudgetAmount    Column = "amount"
	ColumnBudgetCreatedAt Column = "created_at"
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
)

// Columns for table User
const (
	ColumnUserID                    Column = "id"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
	ColumnUserEmail                 Column = "email"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserName                  Column = "name"
	ColumnUserTier                  Column = "tier"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserUpdatedAt             Column = "updated_at"
)

// Columns for table Transaction
const (
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionID            Column = "id"
)

// Columns for table Investment
const (
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRateType   Column = "rateType"
)

// Columns for table Category
const (
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
)
