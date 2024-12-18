package query

type Table string
type Column string

const (
	TableCreditCard     Table = "credit_cards"
	TableInstitution    Table = "institutions"
	TableBudget         Table = "budgets"
	TableTransaction    Table = "transactions"
	TableInvestment     Table = "investments"
	TableCategory       Table = "categories"
	TableAccount        Table = "accounts"
	TableBudgetCategory Table = "budget_categories"
	TableUser           Table = "users"
)

// Columns for table Investment
const (
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
)

// Columns for table Category
const (
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
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
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
	ColumnUserName                  Column = "name"
	ColumnUserTier                  Column = "tier"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserID                    Column = "id"
	ColumnUserEmail                 Column = "email"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserUpdatedAt             Column = "updated_at"
)

// Columns for table Transaction
const (
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionName          Column = "name"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionAccountID     Column = "account_id"
)

// Columns for table Account
const (
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountID            Column = "id"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountType          Column = "type"
	ColumnAccountUserID        Column = "user_id"
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
	ColumnBudgetDate      Column = "date"
	ColumnBudgetCreatedAt Column = "created_at"
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
)
