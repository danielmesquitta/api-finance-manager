package handler

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type QueryParam = string

const (
	QueryParamSearch           QueryParam = "search"
	QueryParamPage             QueryParam = "page"
	QueryParamPageSize         QueryParam = "page_size"
	QueryParamDate             QueryParam = "date"
	QueryParamStartDate        QueryParam = "start_date"
	QueryParamEndDate          QueryParam = "end_date"
	QueryParamInstitutionIDs   QueryParam = "institution_ids"
	QueryParamCategoryIDs      QueryParam = "category_ids"
	QueryParamUserIDs          QueryParam = "user_ids"
	QueryParamIsExpense        QueryParam = "is_expense"
	QueryParamIsIncome         QueryParam = "is_income"
	QueryParamIsIgnored        QueryParam = "is_ignored"
	QueryParamPaymentMethodIDs QueryParam = "payment_method_ids"
)

type PathParam = string

const (
	pathParamCategoryID    PathParam = "category_id"
	pathParamTransactionID PathParam = "transaction_id"
	pathParamAIChatID      PathParam = "ai_chat_id"
)

func parsePaginationParams(
	c *fiber.Ctx,
) usecase.PaginationInput {
	page := c.QueryInt(QueryParamPage, 1)
	pageSize := c.QueryInt(QueryParamPageSize, 20)

	return usecase.PaginationInput{
		Page:     uint(page),
		PageSize: uint(pageSize),
	}
}

func parseDateParam(
	c *fiber.Ctx,
	param QueryParam,
) (time.Time, error) {
	paramValue := c.Query(param)
	if paramValue == "" {
		return time.Time{}, nil
	}

	date, err := time.Parse(time.RFC3339, paramValue)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDate
	}

	log.Println("Date: ", date)

	return date, nil
}

func parseUUIDsParam(
	c *fiber.Ctx,
	param QueryParam,
) ([]uuid.UUID, error) {
	paramValue := c.Query(param)
	if paramValue == "" {
		return nil, nil
	}

	values := strings.Split(paramValue, ",")

	var uuids []uuid.UUID
	for _, value := range values {
		if value == "" {
			continue
		}
		id, err := uuid.Parse(strings.TrimSpace(value))
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}
		uuids = append(uuids, id)
	}

	return uuids, nil
}

func parseBoolParam(
	c *fiber.Ctx,
	param QueryParam,
) (bool, error) {
	paramValue := c.QueryBool(param)
	return paramValue, nil
}

func parseNillableBoolParam(
	c *fiber.Ctx,
	param QueryParam,
) (*bool, error) {
	paramValue := c.Query(param)
	if paramValue == "" {
		return nil, nil
	}

	b, err := strconv.ParseBool(paramValue)
	if err != nil {
		return nil, errs.ErrInvalidBool
	}

	return &b, nil
}

func GetUserClaims(
	c *fiber.Ctx,
) *jwtutil.UserClaims {
	token, ok := c.Locals(jwtutil.ClaimsKey).(*jwt.Token)
	if !ok {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	issuer, _ := claims.GetIssuer()
	issuedAt, _ := claims.GetIssuedAt()
	expiresAt, _ := claims.GetExpirationTime()
	tier, _ := claims["tier"].(string)

	var subscriptionExpiresAt *time.Time
	if expiresAtStr, ok := claims["subscription_expires_at"].(string); ok &&
		expiresAtStr != "" {
		if parsed, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
			subscriptionExpiresAt = &parsed
		}
	}

	return &jwtutil.UserClaims{
		Issuer:                issuer,
		IssuedAt:              issuedAt.Time,
		ExpiresAt:             expiresAt.Time,
		Tier:                  entity.Tier(tier),
		SubscriptionExpiresAt: subscriptionExpiresAt,
	}
}

func prepareTransactionOptions(
	c *fiber.Ctx,
) (*repo.TransactionOptions, error) {
	search := c.Query(QueryParamSearch)

	paymentMethodIDs, err := parseUUIDsParam(c, QueryParamPaymentMethodIDs)
	if err != nil {
		return nil, errs.New(err)
	}

	institutionIDs, err := parseUUIDsParam(c, QueryParamInstitutionIDs)
	if err != nil {
		return nil, errs.New(err)
	}

	categoryIDs, err := parseUUIDsParam(c, QueryParamCategoryIDs)
	if err != nil {
		return nil, errs.New(err)
	}

	isExpense, err := parseBoolParam(c, QueryParamIsExpense)
	if err != nil {
		return nil, errs.New(err)
	}

	isIncome, err := parseBoolParam(c, QueryParamIsIncome)
	if err != nil {
		return nil, errs.New(err)
	}

	isIgnored, err := parseNillableBoolParam(c, QueryParamIsIgnored)
	if err != nil {
		return nil, errs.New(err)
	}

	startDate, err := parseDateParam(c, QueryParamStartDate)
	if err != nil {
		return nil, errs.New(err)
	}

	endDate, err := parseDateParam(c, QueryParamEndDate)
	if err != nil {
		return nil, errs.New(err)
	}

	return &repo.TransactionOptions{
		Search:           search,
		CategoryIDs:      categoryIDs,
		InstitutionIDs:   institutionIDs,
		PaymentMethodIDs: paymentMethodIDs,
		IsExpense:        isExpense,
		IsIncome:         isIncome,
		IsIgnored:        isIgnored,
		StartDate:        startDate,
		EndDate:          endDate,
	}, nil
}
