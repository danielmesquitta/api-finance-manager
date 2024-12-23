package query

type Table string
type Column string

const (
	TableUser           Table = "users"
	TableAccount        Table = "accounts"
	TableCreditCard     Table = "credit_cards"
	TableInstitution    Table = "institutions"
	TableTransaction    Table = "transactions"
	TableInvestment     Table = "investments"
	TableCategory       Table = "categories"
	TableBudget         Table = "budgets"
	TableBudgetCategory Table = "budget_categories"
)

// Columns for table CreditCard
const (
	ColumnCreditCardLimit          Column = "limit"
	ColumnCreditCardAvailableLimit Column = "available_limit"
	ColumnCreditCardCreatedAt      Column = "created_at"
	ColumnCreditCardUpdatedAt      Column = "updated_at"
	ColumnCreditCardAccountID      Column = "account_id"
	ColumnCreditCardID             Column = "id"
	ColumnCreditCardLevel          Column = "level"
	ColumnCreditCardBrand          Column = "brand"
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

// Columns for table Transaction
const (
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionName          Column = "name"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionPaymentMethod Column = "payment_method"
)

// Columns for table Investment
const (
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentID         Column = "id"
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

// Columns for table Account
const (
	ColumnAccountBalance       Column = "balance"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountID            Column = "id"
	ColumnAccountName          Column = "name"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountType          Column = "type"
)

// Columns for table Budget
const (
	ColumnBudgetCreatedAt Column = "created_at"
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
	ColumnBudgetID        Column = "id"
	ColumnBudgetAmount    Column = "amount"
	ColumnBudgetDate      Column = "date"
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

// Columns for table User
const (
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserUpdatedAt             Column = "updated_at"
	ColumnUserProvider              Column = "provider"
	ColumnUserName                  Column = "name"
	ColumnUserEmail                 Column = "email"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserID                    Column = "id"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserTier                  Column = "tier"
)
