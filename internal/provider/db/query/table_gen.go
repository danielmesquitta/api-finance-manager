package query

type Table string
type Column string

const (
	TableInvestment     Table = "investments"
	TableCategory       Table = "categories"
	TableAccount        Table = "accounts"
	TableBudget         Table = "budgets"
	TableUser           Table = "users"
	TableTransaction    Table = "transactions"
	TableCreditCard     Table = "credit_cards"
	TableInstitution    Table = "institutions"
	TableBudgetCategory Table = "budget_categories"
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
	ColumnTransactionPaymentMethod Column = "payment_method"
	ColumnTransactionIsIgnored     Column = "is_ignored"
	ColumnTransactionAccountID     Column = "account_id"
	ColumnTransactionDescription   Column = "description"
	ColumnTransactionAmount        Column = "amount"
	ColumnTransactionDate          Column = "date"
	ColumnTransactionCreatedAt     Column = "created_at"
	ColumnTransactionUpdatedAt     Column = "updated_at"
	ColumnTransactionID            Column = "id"
	ColumnTransactionExternalID    Column = "external_id"
	ColumnTransactionName          Column = "name"
	ColumnTransactionUserID        Column = "user_id"
	ColumnTransactionCategoryID    Column = "category_id"
)

// Columns for table Investment
const (
	ColumnInvestmentExternalID Column = "external_id"
	ColumnInvestmentAmount     Column = "amount"
	ColumnInvestmentType       Column = "type"
	ColumnInvestmentRateType   Column = "rateType"
	ColumnInvestmentID         Column = "id"
	ColumnInvestmentRate       Column = "rate"
	ColumnInvestmentCreatedAt  Column = "created_at"
	ColumnInvestmentUpdatedAt  Column = "updated_at"
	ColumnInvestmentUserID     Column = "user_id"
	ColumnInvestmentName       Column = "name"
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

// Columns for table User
const (
	ColumnUserID                    Column = "id"
	ColumnUserName                  Column = "name"
	ColumnUserEmail                 Column = "email"
	ColumnUserSynchronizedAt        Column = "synchronized_at"
	ColumnUserExternalID            Column = "external_id"
	ColumnUserProvider              Column = "provider"
	ColumnUserVerifiedEmail         Column = "verified_email"
	ColumnUserTier                  Column = "tier"
	ColumnUserAvatar                Column = "avatar"
	ColumnUserSubscriptionExpiresAt Column = "subscription_expires_at"
	ColumnUserCreatedAt             Column = "created_at"
	ColumnUserUpdatedAt             Column = "updated_at"
)

// Columns for table Account
const (
	ColumnAccountID            Column = "id"
	ColumnAccountExternalID    Column = "external_id"
	ColumnAccountUserID        Column = "user_id"
	ColumnAccountName          Column = "name"
	ColumnAccountBalance       Column = "balance"
	ColumnAccountType          Column = "type"
	ColumnAccountCreatedAt     Column = "created_at"
	ColumnAccountUpdatedAt     Column = "updated_at"
	ColumnAccountInstitutionID Column = "institution_id"
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

// Columns for table Institution
const (
	ColumnInstitutionID         Column = "id"
	ColumnInstitutionExternalID Column = "external_id"
	ColumnInstitutionName       Column = "name"
	ColumnInstitutionLogo       Column = "logo"
	ColumnInstitutionCreatedAt  Column = "created_at"
	ColumnInstitutionUpdatedAt  Column = "updated_at"
)
