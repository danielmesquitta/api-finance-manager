package entity

type AccountType string

const (
	AccountTypeBank    AccountType = "BANK"
	AccountTypeCredit  AccountType = "CREDIT"
	AccountTypeUnknown AccountType = "UNKNOWN"
)
