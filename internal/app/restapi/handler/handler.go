package handler

import (
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
	queryParamSearch           QueryParam = "search"
	queryParamPage             QueryParam = "page"
	queryParamPageSize         QueryParam = "page_size"
	queryParamDate             QueryParam = "date"
	queryParamStartDate        QueryParam = "start_date"
	queryParamEndDate          QueryParam = "end_date"
	queryParamInstitutionIDs   QueryParam = "institution_ids"
	queryParamCategoryIDs      QueryParam = "category_ids"
	queryParamUserIDs          QueryParam = "user_ids"
	queryParamIsExpense        QueryParam = "is_expense"
	queryParamIsIncome         QueryParam = "is_income"
	queryParamIsIgnored        QueryParam = "is_ignored"
	queryParamPaymentMethodIDs QueryParam = "payment_method_ids"
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
	page := c.QueryInt(queryParamPage, 1)
	pageSize := c.QueryInt(queryParamPageSize, 20)

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
	token := c.Locals(jwtutil.ClaimsKey).(*jwt.Token)
	if token == nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	return &jwtutil.UserClaims{
		Issuer:                (claims)["iss"].(string),
		IssuedAt:              time.Unix(int64((claims)["iat"].(float64)), 0),
		ExpiresAt:             time.Unix(int64((claims)["exp"].(float64)), 0),
		Tier:                  (claims)["tier"].(entity.Tier),
		SubscriptionExpiresAt: (claims)["subscription_expires_at"].(*time.Time),
	}
}

func prepareTransactionOptions(
	c *fiber.Ctx,
) (*repo.TransactionOptions, error) {
	search := c.Query(queryParamSearch)

	paymentMethodIDs, err := parseUUIDsParam(c, queryParamPaymentMethodIDs)
	if err != nil {
		return nil, errs.New(err)
	}

	institutionIDs, err := parseUUIDsParam(c, queryParamInstitutionIDs)
	if err != nil {
		return nil, errs.New(err)
	}

	categoryIDs, err := parseUUIDsParam(c, queryParamCategoryIDs)
	if err != nil {
		return nil, errs.New(err)
	}

	isExpense, err := parseBoolParam(c, queryParamIsExpense)
	if err != nil {
		return nil, errs.New(err)
	}

	isIncome, err := parseBoolParam(c, queryParamIsExpense)
	if err != nil {
		return nil, errs.New(err)
	}

	isIgnored, err := parseNillableBoolParam(c, queryParamIsIgnored)
	if err != nil {
		return nil, errs.New(err)
	}

	startDate, err := parseDateParam(c, queryParamStartDate)
	if err != nil {
		return nil, errs.New(err)
	}

	endDate, err := parseDateParam(c, queryParamEndDate)
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
