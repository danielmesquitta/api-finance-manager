package pluggy

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/google/uuid"
)

type TransactionsResponse struct {
	Total      int64    `json:"total"`
	TotalPages int64    `json:"totalPages"`
	Page       int64    `json:"page"`
	Results    []Result `json:"results"`
}

type Result struct {
	ID                      string              `json:"id"`
	Description             string              `json:"description"`
	DescriptionRaw          string              `json:"descriptionRaw"`
	CurrencyCode            CurrencyCode        `json:"currencyCode"`
	Amount                  float64             `json:"amount"`
	AmountInAccountCurrency *float64            `json:"amountInAccountCurrency"`
	Date                    time.Time           `json:"date"`
	Category                *string             `json:"category"`
	CategoryID              *string             `json:"categoryId"`
	AccountID               string              `json:"accountId"`
	Status                  Status              `json:"status"`
	PaymentData             *PaymentData        `json:"paymentData"`
	Type                    ResultType          `json:"type"`
	OperationType           *OperationType      `json:"operationType"`
	CreditCardMetadata      *CreditCardMetadata `json:"creditCardMetadata"`
	CreatedAt               time.Time           `json:"createdAt"`
	UpdatedAt               time.Time           `json:"updatedAt"`
}

type CreditCardMetadata struct {
	CardNumber        *string    `json:"cardNumber,omitempty"`
	PurchaseDate      *time.Time `json:"purchaseDate,omitempty"`
	TotalInstallments *int64     `json:"totalInstallments,omitempty"`
	InstallmentNumber *int64     `json:"installmentNumber,omitempty"`
	BillID            *string    `json:"billId,omitempty"`
	PayeeMCC          *int64     `json:"payeeMCC,omitempty"`
}

type PaymentData struct {
	Payer               *Payer         `json:"payer"`
	PaymentMethod       *PaymentMethod `json:"paymentMethod"`
	Receiver            *Payer         `json:"receiver"`
	ReceiverReferenceID *string        `json:"receiverReferenceId"`
}

type Payer struct {
	AccountNumber     *string         `json:"accountNumber"`
	BranchNumber      *string         `json:"branchNumber"`
	DocumentNumber    *DocumentNumber `json:"documentNumber"`
	Name              *string         `json:"name"`
	RoutingNumber     *string         `json:"routingNumber"`
	RoutingNumberISPB *string         `json:"routingNumberISPB"`
}

type DocumentNumber struct {
	Type  DocumentType `json:"type"`
	Value string       `json:"value"`
}

type CurrencyCode string

const (
	CurrencyCodeBRL CurrencyCode = "BRL"
	CurrencyCodeUSD CurrencyCode = "USD"
)

type OperationType string

const (
	OperationTypeCartao                        OperationType = "CARTAO"
	OperationTypeConvenioArrecadacao           OperationType = "CONVENIO_ARRECADACAO"
	OperationTypeDeposito                      OperationType = "DEPOSITO"
	OperationTypeEncargosJurosChequeEspecial   OperationType = "ENCARGOS_JUROS_CHEQUE_ESPECIAL"
	OperationTypeOperacaoCredito               OperationType = "OPERACAO_CREDITO"
	OperationTypeOperationTypeBOLETO           OperationType = "BOLETO"
	OperationTypeOperationTypePIX              OperationType = "PIX"
	OperationTypeOperationTypeTED              OperationType = "TED"
	OperationTypeOutros                        OperationType = "OUTROS"
	OperationTypePacoteTarifaServicos          OperationType = "PACOTE_TARIFA_SERVICOS"
	OperationTypeRendimentoAplicFinanceira     OperationType = "RENDIMENTO_APLIC_FINANCEIRA"
	OperationTypeResgateAplicFinanceira        OperationType = "RESGATE_APLIC_FINANCEIRA"
	OperationTypeSaque                         OperationType = "SAQUE"
	OperationTypeTarifaServicosAvulsos         OperationType = "TARIFA_SERVICOS_AVULSOS"
	OperationTypeTransferenciaMesmaInstituicao OperationType = "TRANSFERENCIA_MESMA_INSTITUICAO"
)

type DocumentType string

const (
	DocumentTypeCNPJ DocumentType = "CNPJ"
	DocumentTypeCPF  DocumentType = "CPF"
)

type PaymentMethod string

const (
	PaymentMethodOther      PaymentMethod = "OTHER"
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodBOLETO     PaymentMethod = "BOLETO"
	PaymentMethodDEBIT      PaymentMethod = "DEBIT"
	PaymentMethodPIX        PaymentMethod = "PIX"
	PaymentMethodTED        PaymentMethod = "TED"
	PaymentMethodTEF        PaymentMethod = "TEF"
)

type Status string

const (
	Pending Status = "PENDING"
	Posted  Status = "POSTED"
)

type ResultType string

const (
	Credit ResultType = "CREDIT"
	Debit  ResultType = "DEBIT"
)

func (c *Client) ListTransactions(
	ctx context.Context,
	accountID string,
	options ...openfinance.ListTransactionsOption,
) ([]openfinance.Transaction, error) {
	opts := openfinance.ListTransactionsOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	page := 1
	queryParams := map[string]string{
		"accountId": accountID,
		"pageSize":  "500",
		"page":      strconv.Itoa(page),
	}
	if !opts.StartDate.IsZero() {
		queryParams["from"] = opts.StartDate.Format(time.DateOnly)
	}
	if !opts.EndDate.IsZero() {
		queryParams["to"] = opts.EndDate.Format(time.DateOnly)
	}

	allTransactions := TransactionsResponse{}
	for {
		res, err := c.c.R().
			SetContext(ctx).
			SetQueryParams(queryParams).
			Get("/transactions")
		if err != nil {
			return nil, errs.New(err)
		}
		body := res.Body()
		if res.IsError() {
			return nil, errs.New(body)
		}

		transRes := TransactionsResponse{}
		if err := json.Unmarshal(body, &transRes); err != nil {
			return nil, errs.New(err)
		}

		allTransactions.Results = append(
			allTransactions.Results,
			transRes.Results...)

		if transRes.Page >= transRes.TotalPages {
			break
		}

		page++
		queryParams["page"] = strconv.Itoa(page)
	}

	transactions := []openfinance.Transaction{}
	for _, t := range allTransactions.Results {
		transaction, err := c.parseTransactionResultToEntity(t)
		if err != nil {
			slog.Error(
				"error parsing transaction result to entity",
				"transaction",
				t,
				"err",
				err,
			)
			continue
		}

		accountUUID := uuid.MustParse(accountID)
		transaction.AccountID = &accountUUID

		transactions = append(transactions, *transaction)
	}

	return transactions, nil
}

func (c *Client) parseTransactionResultToEntity(
	r Result,
) (*openfinance.Transaction, error) {
	transaction := openfinance.Transaction{
		Transaction: entity.Transaction{
			ExternalID: r.ID,
			Name:       r.Description,
			Date:       r.Date,
		},
	}

	if r.CategoryID != nil {
		transaction.CategoryExternalID = *r.CategoryID
	}

	if err := c.setTransactionAmount(&transaction, r); err != nil {
		return nil, errs.New(err)
	}

	c.setTransactionPaymentMethod(&transaction, r)

	return &transaction, nil
}

func (c *Client) setTransactionAmount(
	t *openfinance.Transaction,
	r Result,
) error {
	if r.CurrencyCode == CurrencyCodeBRL {
		t.Amount = money.ToCents(r.Amount)
		return nil
	}

	if r.AmountInAccountCurrency != nil {
		t.Amount = money.ToCents(*r.AmountInAccountCurrency)
		return nil
	}

	return errs.New("amount is nil")
}

func (c *Client) setTransactionPaymentMethod(
	t *openfinance.Transaction,
	r Result,
) {
	if r.CreditCardMetadata != nil {
		t.PaymentMethod = string(entity.PaymentMethodCreditCard)
		return
	}

	if strings.HasPrefix(r.Description, "COMPRA C/CARTAO") {
		t.PaymentMethod = string(entity.PaymentMethodCreditCard)
		return
	}

	if r.OperationType != nil && *r.OperationType == OperationTypeCartao {
		t.PaymentMethod = string(entity.PaymentMethodCreditCard)
		return
	}

	if r.PaymentData != nil && r.PaymentData.PaymentMethod != nil {
		t.PaymentMethod = string(*r.PaymentData.PaymentMethod)
		return
	}

	t.PaymentMethod = string(entity.PaymentMethodOther)
}
