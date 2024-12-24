package query

type Table string
type Column string

const (
	TableUser           Table = "users"
	TableInstitution    Table = "institutions"
	TableBudgetCategory Table = "budget_categories"
	TableCreditCard     Table = "credit_cards"
	TableBudget         Table = "budgets"
	TableTransaction    Table = "transactions"
	TableInvestment     Table = "investments"
	TableCategory       Table = "categories"
	TableAccount        Table = "accounts"
)

// Columns for table Transaction
const (
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionCategoryID    Column = "category_id"
)

// Columns for table Investment
const (
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentName       Column = "name"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentCreatedAt  Column = "created_at"
)

// Columns for table Account
const (
	ColumnAccountID            Column = "id"
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountType          Column = "type"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountInstitutionID Column = "institution_id"
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
	ColumnBudgetDate      Column = "date"
	ColumnBudgetCreatedAt Column = "created_at"
	ColumnBudgetUpdatedAt Column = "updated_at"
	ColumnBudgetUserID    Column = "user_id"
	ColumnBudgetID        Column = "id"
	ColumnBudgetAmount    Column = "amount"
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

// Columns for table User
const (
	ColumnUserProvider              Column = "provider"
	ColumnUserName                  Column = "name"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserTier                  Column = "tier"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserUpdatedAt             Column = "updated_at"
	ColumnUserID                    Column = "id"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserEmail                 Column = "email"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserCreatedAt             Column = "created_at"
)

// Columns for table Category
const (
	ColumnCategoryID         Column = "id"
	ColumnCategoryExternalID Column = "external_id"
	ColumnCategoryName       Column = "name"
	ColumnCategoryCreatedAt  Column = "created_at"
	ColumnCategoryUpdatedAt  Column = "updated_at"
)

// Columns for table CreditCard
const (
	ColumnCreditCardUpdatedAt      Column = "updated_at"
	ColumnCreditCardAccountID      Column = "account_id"
	ColumnCreditCardID             Column = "id"
	ColumnCreditCardLevel          Column = "level"
	ColumnCreditCardBrand          Column = "brand"
	ColumnCreditCardLimit          Column = "limit"
	ColumnCreditCardAvailableLimit Column = "available_limit"
	ColumnCreditCardCreatedAt      Column = "created_at"
)
