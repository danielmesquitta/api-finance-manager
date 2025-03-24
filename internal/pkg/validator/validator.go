package validator

import (
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

const defaultLangKey = entity.LanguagePortuguese

var langsByKey = map[entity.Language]struct {
	translator               locales.Translator
	registerTranslationsFunc func(v *validator.Validate, trans ut.Translator) (err error)
}{
	entity.LanguagePortuguese: {
		translator:               pt_BR.New(),
		registerTranslationsFunc: ptBRTranslations.RegisterDefaultTranslations,
	},
	entity.LanguageEnglish: {
		translator:               en.New(),
		registerTranslationsFunc: enTranslations.RegisterDefaultTranslations,
	},
	entity.LanguageSpanish: {
		translator:               es.New(),
		registerTranslationsFunc: esTranslations.RegisterDefaultTranslations,
	},
}

type validatorLanguage struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

type Validator struct {
	vls map[entity.Language]validatorLanguage
}

func New() *Validator {
	defaultLang := langsByKey[defaultLangKey]

	translators := []locales.Translator{}
	for _, lang := range langsByKey {
		translators = append(translators, lang.translator)
	}

	uni := ut.New(defaultLang.translator, translators...)

	vls := map[entity.Language]validatorLanguage{}
	for k, v := range langsByKey {
		translator, ok := uni.GetTranslator(string(k))
		if !ok {
			log.Fatalf("translator not found for language: %s", k)
		}

		val := validator.New()
		if err := v.registerTranslationsFunc(val, translator); err != nil {
			log.Fatalln(err)
		}
		vls[k] = validatorLanguage{
			Validate:   val,
			Translator: translator,
		}
	}

	return &Validator{
		vls: vls,
	}
}

// Validate validates the struct data using the provided language code ("pt_BR", "en", "es").
// It returns an error with translated messages if the data is invalid.
func (v *Validator) Validate(data any, language ...entity.Language) error {
	if len(language) == 0 {
		language = append(language, defaultLangKey)
	}
	l := language[0]

	vl, ok := v.vls[l]
	if !ok {
		log.Printf(
			"validator not found for language %s, using fallback %s",
			l,
			defaultLangKey,
		)
		vl = v.vls[defaultLangKey]
	}

	err := vl.Validate.Struct(data)
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
		strErrs[i] = validationErr.Translate(vl.Translator)
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
