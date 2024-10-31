package validator

import (
	"log"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	Validate(data any) error
}

type Validate struct {
	v *validator.Validate
	t ut.Translator
}

func NewValidate() *Validate {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	t, ok := uni.GetTranslator("en")
	if !ok {
		log.Fatalln("translator not found")
	}

	if err := enTranslations.
		RegisterDefaultTranslations(validate, t); err != nil {
		log.Fatalln(err)
	}

	return &Validate{
		validate,
		t,
	}
}

// Validate validates the data (struct)
// returning an error if the data is invalid.
func (v *Validate) Validate(
	data any,
) error {
	err := v.v.Struct(data)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	strErrs := make([]string, len(validationErrs))
	for i, validationErr := range validationErrs {
		strErrs[i] = validationErr.Translate(v.t)
	}

	errMsg := strings.Join(
		strErrs,
		", ",
	)

	return errs.NewWithType(errMsg, errs.ErrTypeValidation)
}

var _ Validator = (*Validate)(nil)
