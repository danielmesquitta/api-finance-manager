package usecase_test

import (
	"context"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestCalculateEmergencyReserveExecute(
	t *testing.T,
) {
	t.Parallel()
	asserts := assert.New(t)

	t.Run("should return the correct output", func(t *testing.T) {
		t.Parallel()

		v := validator.New()
		uc := usecase.NewCalculateEmergencyReserve(v)

		in := usecase.CalculateEmergencyReserveInput{
			JobType:                  entity.JobTypeEmployee,
			MonthlyExpenses:          100000, // R$1000
			MonthlyIncome:            200000, // R$2000
			MonthlySavingsPercentage: 2500,   // 25%
		}

		out, err := uc.Execute(context.Background(), in)

		asserts.Nil(err)
		asserts.NotNil(out)
		asserts.Equal(int64(6), out.RecommendedReserveInMonths)
		asserts.Equal(int64(600000), out.RecommendedReserveInValue)
		asserts.Equal(int64(12), out.MonthsToAchieveEmergencyReserve)
	})
}
