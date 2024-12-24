package query

type Table string
type Column string

const (
	TableCreditCard     Table = "credit_cards"
	TableBudget         Table = "budgets"
	TableUser           Table = "users"
	TableTransaction    Table = "transactions"
	TableCategory       Table = "categories"
	TableBudgetCategory Table = "budget_categories"
	TableInvestment     Table = "investments"
	TableAccount        Table = "accounts"
	TableInstitution    Table = "institutions"
)

// Columns for table Transaction
const (
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionUpdatedAt     Column = "updated_at"
)

// Columns for table Investment
const (
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
)

// Columns for table CreditCard
const (
	ColumnCreditCardAvailableLimit Column = "available_limit"
	ColumnCreditCardCreatedAt      Column = "created_at"
	ColumnCreditCardUpdatedAt      Column = "updated_at"
	ColumnCreditCardAccountID      Column = "account_id"
	ColumnCreditCardID             Column = "id"
	ColumnCreditCardLevel          Column = "level"
	ColumnCreditCardBrand          Column = "brand"
	ColumnCreditCardLimit          Column = "limit"
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

// Columns for table User
const (
	ColumnUserID                    Column = "id"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserName                  Column = "name"
	ColumnUserEmail                 Column = "email"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserTier                  Column = "tier"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserUpdatedAt             Column = "updated_at"
)

// Columns for table Account
const (
	ColumnAccountID            Column = "id"
	ColumnAccountType          Column = "type"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
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

// Columns for table BudgetCategory
const (
	ColumnBudgetCategoryCategoryID Column = "category_id"
	ColumnBudgetCategoryID         Column = "id"
	ColumnBudgetCategoryAmount     Column = "amount"
	ColumnBudgetCategoryCreatedAt  Column = "created_at"
	ColumnBudgetCategoryUpdatedAt  Column = "updated_at"
	ColumnBudgetCategoryBudgetID   Column = "budget_id"
)

// Columns for table Category
const (
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
)
