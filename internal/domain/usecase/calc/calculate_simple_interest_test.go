package calc

import (
	"context"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestCalculateSimpleInterestUseCaseExecute(
	t *testing.T,
) {
	t.Parallel()
	asserts := assert.New(t)

	t.Run("should return the correct output", func(t *testing.T) {
		t.Parallel()

		// Arrange
		v := validator.New()
		uc := NewCalculateSimpleInterestUseCase(v)

		in := CalculateSimpleInterestUseCaseInput{
			InitialDeposit: 1000, // R$10
			Interest:       1000, // 10%
			InterestType:   entity.InterestType("ANNUAL"),
			PeriodInMonths: 12,
		}

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		asserts.Nil(err)
		asserts.NotNil(out)
		asserts.Equal(int64(1100), out.TotalAmount)             // R$11
		asserts.Equal(int64(100), out.TotalInterest)            // R$1
		asserts.Equal(int64(1000), out.TotalDeposit)            // R$10
		asserts.Equal(out.ByMonth[6].TotalDeposit, int64(1000)) // R$10
		asserts.Equal(out.ByMonth[6].TotalInterest, int64(50))  // R$0.50
		asserts.Equal(out.ByMonth[6].TotalAmount, int64(1050))  // R$10.50
		asserts.Len(out.ByMonth, 12)
	})
}
