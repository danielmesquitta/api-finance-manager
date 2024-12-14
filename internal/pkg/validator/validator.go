package validator

import (
	"log"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ptBRTranslations "github.com/go-playground/validator/v10/translations/pt_BR"
)

type Validator struct {
	v *validator.Validate
	t ut.Translator
}

func NewValidator() *Validator {
	validate := validator.New()
	portuguese := pt_BR.New()
	uni := ut.New(portuguese, portuguese)
	t, ok := uni.GetTranslator("pt_BR")
	if !ok {
		log.Fatalln("translator not found")
	}

	if err := ptBRTranslations.
		RegisterDefaultTranslations(validate, t); err != nil {
		log.Fatalln(err)
	}

	return &Validator{
		validate,
		t,
	}
}

// Validate validates the data (struct)
// returning an error if the data is invalid.
func (v *Validator) Validate(
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
