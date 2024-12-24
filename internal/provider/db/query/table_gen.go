package query

type Table string
type Column string

const (
	TableInstitution    Table = "institutions"
	TableBudgetCategory Table = "budget_categories"
	TableUser           Table = "users"
	TableInvestment     Table = "investments"
	TableCreditCard     Table = "credit_cards"
	TableBudget         Table = "budgets"
	TableTransaction    Table = "transactions"
	TableCategory       Table = "categories"
	TableAccount        Table = "accounts"
)

// Columns for table Institution
const (
	ColumnInstitutionLogo       Column = "logo"
	ColumnInstitutionCreatedAt  Column = "created_at"
	ColumnInstitutionUpdatedAt  Column = "updated_at"
	ColumnInstitutionID         Column = "id"
	ColumnInstitutionExternalID Column = "external_id"
	ColumnInstitutionName       Column = "name"
)

// Columns for table User
const (
	ColumnUserName                  Column = "name"
	ColumnUserTier                  Column = "tier"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserID                    Column = "id"
	ColumnUserProvider              Column = "provider"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserUpdatedAt             Column = "updated_at"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserEmail                 Column = "email"
)

// Columns for table Transaction
const (
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionCreatedAt     Column = "created_at"
)

// Columns for table Account
const (
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountType          Column = "type"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountID            Column = "id"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
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

// Columns for table Investment
const (
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
)

// Columns for table Category
const (
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
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
