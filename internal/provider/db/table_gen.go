// Code generated by prisma-go-tools. DO NOT EDIT.

package db

import "fmt"

type tableAIChat string

func (t tableAIChat) String() string {
	return string(t)
}

func (t tableAIChat) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableAIChat) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableAIChat) ColumnTitle() string {
	return fmt.Sprintf("%s.title", t)
}

func (t tableAIChat) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableAIChat) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableAIChat) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableAIChat) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

const TableAiChat = tableAIChat("ai_chats")

type tableAIChatMessage string

func (t tableAIChatMessage) String() string {
	return string(t)
}

func (t tableAIChatMessage) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableAIChatMessage) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableAIChatMessage) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableAIChatMessage) ColumnCreatedByUserID() string {
	return fmt.Sprintf("%s.created_by_user_id", t)
}

func (t tableAIChatMessage) ColumnAiChatID() string {
	return fmt.Sprintf("%s.ai_chat_id", t)
}

func (t tableAIChatMessage) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableAIChatMessage) ColumnMessage() string {
	return fmt.Sprintf("%s.message", t)
}

func (t tableAIChatMessage) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

const TableAiChatMessage = tableAIChatMessage("ai_chat_messages")

type tableAccount string

func (t tableAccount) String() string {
	return string(t)
}

func (t tableAccount) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableAccount) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableAccount) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t tableAccount) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t tableAccount) ColumnType() string {
	return fmt.Sprintf("%s.type", t)
}

func (t tableAccount) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableAccount) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableAccount) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

func (t tableAccount) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableAccount) ColumnInstitutionID() string {
	return fmt.Sprintf("%s.institution_id", t)
}

const TableAccount = tableAccount("accounts")

type tableAccountBalance string

func (t tableAccountBalance) String() string {
	return string(t)
}

func (t tableAccountBalance) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableAccountBalance) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableAccountBalance) ColumnAccountID() string {
	return fmt.Sprintf("%s.account_id", t)
}

func (t tableAccountBalance) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

func (t tableAccountBalance) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableAccountBalance) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t tableAccountBalance) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableAccountBalance) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

const TableAccountBalance = tableAccountBalance("account_balances")

type tableBudget string

func (t tableBudget) String() string {
	return string(t)
}

func (t tableBudget) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableBudget) ColumnDate() string {
	return fmt.Sprintf("%s.date", t)
}

func (t tableBudget) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableBudget) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableBudget) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableBudget) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

func (t tableBudget) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableBudget) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

const TableBudget = tableBudget("budgets")

type tableBudgetCategory string

func (t tableBudgetCategory) String() string {
	return string(t)
}

func (t tableBudgetCategory) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableBudgetCategory) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableBudgetCategory) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t tableBudgetCategory) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableBudgetCategory) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableBudgetCategory) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableBudgetCategory) ColumnBudgetID() string {
	return fmt.Sprintf("%s.budget_id", t)
}

func (t tableBudgetCategory) ColumnCategoryID() string {
	return fmt.Sprintf("%s.category_id", t)
}

const TableBudgetCategory = tableBudgetCategory("budget_categories")

type tableFeedback string

func (t tableFeedback) String() string {
	return string(t)
}

func (t tableFeedback) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableFeedback) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableFeedback) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

func (t tableFeedback) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableFeedback) ColumnMessage() string {
	return fmt.Sprintf("%s.message", t)
}

func (t tableFeedback) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

const TableFeedback = tableFeedback("feedbacks")

type tableInstitution string

func (t tableInstitution) String() string {
	return string(t)
}

func (t tableInstitution) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableInstitution) ColumnLogo() string {
	return fmt.Sprintf("%s.logo", t)
}

func (t tableInstitution) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableInstitution) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableInstitution) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableInstitution) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableInstitution) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t tableInstitution) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

const TableInstitution = tableInstitution("institutions")

type tableInvestment string

func (t tableInvestment) String() string {
	return string(t)
}

func (t tableInvestment) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableInvestment) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableInvestment) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t tableInvestment) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t tableInvestment) ColumnRate() string {
	return fmt.Sprintf("%s.rate", t)
}

func (t tableInvestment) ColumnRateType() string {
	return fmt.Sprintf("%s.rateType", t)
}

func (t tableInvestment) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableInvestment) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableInvestment) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableInvestment) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t tableInvestment) ColumnCategoryID() string {
	return fmt.Sprintf("%s.category_id", t)
}

func (t tableInvestment) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

const TableInvestment = tableInvestment("investments")

type tableInvestmentCategory string

func (t tableInvestmentCategory) String() string {
	return string(t)
}

func (t tableInvestmentCategory) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableInvestmentCategory) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t tableInvestmentCategory) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t tableInvestmentCategory) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableInvestmentCategory) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableInvestmentCategory) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableInvestmentCategory) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

const TableInvestmentCategory = tableInvestmentCategory("investment_categories")

type tablePaymentMethod string

func (t tablePaymentMethod) String() string {
	return string(t)
}

func (t tablePaymentMethod) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tablePaymentMethod) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tablePaymentMethod) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tablePaymentMethod) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t tablePaymentMethod) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t tablePaymentMethod) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tablePaymentMethod) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

const TablePaymentMethod = tablePaymentMethod("payment_methods")

type tableTransaction string

func (t tableTransaction) String() string {
	return string(t)
}

func (t tableTransaction) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableTransaction) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableTransaction) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t tableTransaction) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t tableTransaction) ColumnAmount() string {
	return fmt.Sprintf("%s.amount", t)
}

func (t tableTransaction) ColumnDate() string {
	return fmt.Sprintf("%s.date", t)
}

func (t tableTransaction) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableTransaction) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableTransaction) ColumnUserID() string {
	return fmt.Sprintf("%s.user_id", t)
}

func (t tableTransaction) ColumnIsIgnored() string {
	return fmt.Sprintf("%s.is_ignored", t)
}

func (t tableTransaction) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableTransaction) ColumnPaymentMethodID() string {
	return fmt.Sprintf("%s.payment_method_id", t)
}

func (t tableTransaction) ColumnAccountID() string {
	return fmt.Sprintf("%s.account_id", t)
}

func (t tableTransaction) ColumnInstitutionID() string {
	return fmt.Sprintf("%s.institution_id", t)
}

func (t tableTransaction) ColumnCategoryID() string {
	return fmt.Sprintf("%s.category_id", t)
}

const TableTransaction = tableTransaction("transactions")

type tableTransactionCategory string

func (t tableTransactionCategory) String() string {
	return string(t)
}

func (t tableTransactionCategory) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableTransactionCategory) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableTransactionCategory) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableTransactionCategory) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableTransactionCategory) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

func (t tableTransactionCategory) ColumnExternalID() string {
	return fmt.Sprintf("%s.external_id", t)
}

func (t tableTransactionCategory) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

const TableTransactionCategory = tableTransactionCategory("transaction_categories")

type tableUser string

func (t tableUser) String() string {
	return string(t)
}

func (t tableUser) ColumnAll() string {
	return fmt.Sprintf("%s.*", t)
}

func (t tableUser) ColumnAuthID() string {
	return fmt.Sprintf("%s.auth_id", t)
}

func (t tableUser) ColumnProvider() string {
	return fmt.Sprintf("%s.provider", t)
}

func (t tableUser) ColumnTier() string {
	return fmt.Sprintf("%s.tier", t)
}

func (t tableUser) ColumnAvatar() string {
	return fmt.Sprintf("%s.avatar", t)
}

func (t tableUser) ColumnSubscriptionExpiresAt() string {
	return fmt.Sprintf("%s.subscription_expires_at", t)
}

func (t tableUser) ColumnSynchronizedAt() string {
	return fmt.Sprintf("%s.synchronized_at", t)
}

func (t tableUser) ColumnCreatedAt() string {
	return fmt.Sprintf("%s.created_at", t)
}

func (t tableUser) ColumnOpenFinanceID() string {
	return fmt.Sprintf("%s.open_finance_id", t)
}

func (t tableUser) ColumnName() string {
	return fmt.Sprintf("%s.name", t)
}

func (t tableUser) ColumnEmail() string {
	return fmt.Sprintf("%s.email", t)
}

func (t tableUser) ColumnVerifiedEmail() string {
	return fmt.Sprintf("%s.verified_email", t)
}

func (t tableUser) ColumnUpdatedAt() string {
	return fmt.Sprintf("%s.updated_at", t)
}

func (t tableUser) ColumnDeletedAt() string {
	return fmt.Sprintf("%s.deleted_at", t)
}

func (t tableUser) ColumnID() string {
	return fmt.Sprintf("%s.id", t)
}

const TableUser = tableUser("users")
