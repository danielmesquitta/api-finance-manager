package validator

import (
	"fmt"
	"log"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	esTranslations "github.com/go-playground/validator/v10/translations/es"
	ptBRTranslations "github.com/go-playground/validator/v10/translations/pt_BR"
)

type Validator struct {
	v   *validator.Validate
	uni *ut.UniversalTranslator
}

func New() *Validator {
	langsByKey := map[entity.Language]locales.Translator{
		entity.LanguagePortuguese: pt_BR.New(),
		entity.LanguageEnglish:    en.New(),
		entity.LanguageSpanish:    es.New(),
	}

	defaultLangKey := entity.LanguagePortuguese
	defaultLang := langsByKey[defaultLangKey]

	langs := []locales.Translator{}
	for _, lang := range langsByKey {
		langs = append(langs, lang)
	}

	uni := ut.New(defaultLang, langs...)

	validate := validator.New()

	t, ok := uni.GetTranslator(string(defaultLangKey))
	if !ok {
		log.Fatalf("translator for %s not found", defaultLangKey)
	}
	if err := enTranslations.RegisterDefaultTranslations(validate, t); err != nil {
		log.Fatalln(err)
	}

	return &Validator{
		v:   validate,
		uni: uni,
	}
}

// Validate validates the struct data using the provided language code ("pt_BR", "en", "es").
// It returns an error with translated messages if the data is invalid.
func (v *Validator) Validate(data any, lang ...entity.Language) error {
	if len(lang) == 0 {
		lang = append(lang, entity.LanguageEnglish)
	}
	l := lang[0]

	translator, ok := v.uni.GetTranslator(string(l))
	if !ok {
		return errs.New(
			fmt.Sprintf("translator not found for language: %s", l),
		)
	}

	// Register the default translations based on the language.
	switch l {
	case "pt_BR":
		if err := ptBRTranslations.RegisterDefaultTranslations(v.v, translator); err != nil {
			return err
		}
	case "en":
		if err := enTranslations.RegisterDefaultTranslations(v.v, translator); err != nil {
			return err
		}
	case "es":
		if err := esTranslations.RegisterDefaultTranslations(v.v, translator); err != nil {
			return err
		}
	default:
		// Fallback to English translations.
		if err := enTranslations.RegisterDefaultTranslations(v.v, translator); err != nil {
			return err
		}
	}

	err := v.v.Struct(data)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	strErrs := make([]string, len(validationErrs))
	errItems := make([]errs.ErrorItem, len(validationErrs))
	for i, validationErr := range validationErrs {
		strErrs[i] = validationErr.Translate(translator)
		errItems[i] = errs.ErrorItem{
			Name:   validationErr.Field(),
			Reason: validationErr.Error(),
		}
	}

	errMsg := strings.Join(strErrs, ", ") + "."
	return &errs.Err{
		Message: errMsg,
		Code:    errs.ErrCodeValidation,
		Errors:  errItems,
	}
}
