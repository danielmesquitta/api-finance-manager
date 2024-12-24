package query

type Table string
type Column string

const (
	TableUser           Table = "users"
	TableCategory       Table = "categories"
	TableAccount        Table = "accounts"
	TableBudget         Table = "budgets"
	TableBudgetCategory Table = "budget_categories"
	TableTransaction    Table = "transactions"
	TableInvestment     Table = "investments"
	TableCreditCard     Table = "credit_cards"
	TableInstitution    Table = "institutions"
)

// Columns for table Transaction
const (
	ColumnTransactionName          Column = "name"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionCategoryID    Column = "category_id"
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionCreatedAt     Column = "created_at"
)

// Columns for table Institution
const (
	ColumnInstitutionCreatedAt  Column = "created_at"
	ColumnInstitutionUpdatedAt  Column = "updated_at"
	ColumnInstitutionID         Column = "id"
	ColumnInstitutionExternalID Column = "external_id"
	ColumnInstitutionName       Column = "name"
	ColumnInstitutionLogo       Column = "logo"
)

// Columns for table CreditCard
const (
	ColumnCreditCardID             Column = "id"
	ColumnCreditCardLevel          Column = "level"
	ColumnCreditCardBrand          Column = "brand"
	ColumnCreditCardLimit          Column = "limit"
	ColumnCreditCardAvailableLimit Column = "available_limit"
	ColumnCreditCardCreatedAt      Column = "created_at"
	ColumnCreditCardUpdatedAt      Column = "updated_at"
	ColumnCreditCardAccountID      Column = "account_id"
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
	ColumnBudgetCategoryBudgetID   Column = "budget_id"
	ColumnBudgetCategoryCategoryID Column = "category_id"
	ColumnBudgetCategoryID         Column = "id"
	ColumnBudgetCategoryAmount     Column = "amount"
	ColumnBudgetCategoryCreatedAt  Column = "created_at"
	ColumnBudgetCategoryUpdatedAt  Column = "updated_at"
)

// Columns for table User
const (
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserTier                  Column = "tier"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserID                    Column = "id"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
	ColumnUserName                  Column = "name"
	ColumnUserEmail                 Column = "email"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserUpdatedAt             Column = "updated_at"
)

// Columns for table Investment
const (
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentCreatedAt  Column = "created_at"
)

// Columns for table Category
const (
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
	ColumnCategoryID         Column = "id"
)

// Columns for table Account
const (
	ColumnAccountType          Column = "type"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountInstitutionID Column = "institution_id"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountName          Column = "name"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountID            Column = "id"
)
