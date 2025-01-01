package usecase_test

import (
	"context"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCompoundInterestExecute(
	t *testing.T,
) {
	t.Parallel()
	asserts := assert.New(t)

	t.Run("should return the correct output", func(t *testing.T) {
		t.Parallel()

		v := validator.New()
		uc := usecase.NewCalculateCompoundInterest(v)

		in := usecase.CalculateCompoundInterestInput{
			InitialDeposit: 1000000, // R$10.000,00
			MonthlyDeposit: 100000,  // R$1.000,00
			Interest:       800,     // 8.00%
			InterestType:   entity.InterestType("ANNUAL"),
			PeriodInMonths: 5 * 12, // 5 years
		}

		out, err := uc.Execute(context.Background(), in)

		asserts.Nil(err)
		asserts.NotNil(out)
		asserts.Equal(
			int64(8763793),
			out.TotalAmount,
		) // R$87.637,93
		asserts.Equal(
			int64(1763793),
			out.TotalInterest,
		) // R$17.637,93
		asserts.Equal(
			int64(7000000),
			out.TotalDeposit,
		) // R$70.000,00

		asserts.Equal(
			int64(1600000),
			out.ByMonth[6].TotalDeposit,
		) // R$16.000,00

		asserts.Equal(
			int64(48965),
			out.ByMonth[6].TotalInterest,
		) // R$489,65

		asserts.Equal(
			int64(1648965),
			out.ByMonth[6].TotalAmount,
		) // R$16.489,65
	})
}
