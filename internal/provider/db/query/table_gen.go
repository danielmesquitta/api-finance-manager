package query

type Table string
type Column string

const (
	TableTransaction    Table = "transactions"
	TableInstitution    Table = "institutions"
	TableAccount        Table = "accounts"
	TableCreditCard     Table = "credit_cards"
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

// Columns for table User
const (
	ColumnUserEmail                 Column = "email"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserUpdatedAt             Column = "updated_at"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserID                    Column = "id"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
	ColumnUserName                  Column = "name"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserTier                  Column = "tier"
)

// Columns for table Investment
const (
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
)

// Columns for table Category
const (
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
)

// Columns for table Budget
const (
	ColumnBudgetAmount    Column = "amount"
	ColumnBudgetDate      Column = "date"
	ColumnBudgetCreatedAt Column = "created_at"
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
	ColumnBudgetID        Column = "id"
)

// Columns for table BudgetCategory
const (
	ColumnBudgetCategoryUpdatedAt  Column = "updated_at"
	ColumnBudgetCategoryBudgetID   Column = "budget_id"
	ColumnBudgetCategoryCategoryID Column = "category_id"
	ColumnBudgetCategoryID         Column = "id"
	ColumnBudgetCategoryAmount     Column = "amount"
	ColumnBudgetCategoryCreatedAt  Column = "created_at"
)

// Columns for table Transaction
const (
	ColumnTransactionDate          Column = "date"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionID            Column = "id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionUserID        Column = "user_id"
)

// Columns for table Account
const (
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountID            Column = "id"
	ColumnAccountType          Column = "type"
	ColumnAccountCreatedAt     Column = "created_at"
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
