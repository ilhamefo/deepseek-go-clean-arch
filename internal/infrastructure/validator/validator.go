package validator

import (
	"errors"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/go-playground/locales/en"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	*validator.Validate
	translator ut.Translator
}

func NewValidator() *Validator {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, found := uni.GetTranslator("en")
	if !found {
		panic("translator not found")
	}

	validate := validator.New()
	if err := enTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic("failed to register translations: " + err.Error())
	}

	return &Validator{
		Validate:   validate,
		translator: trans,
	}
}

// ValidationErrors extracts translated error messages from a validation error.
func (v *Validator) ValidationErrors(err error) map[string]string {
	var ve validator.ValidationErrors
	out := make(map[string]string)
	if errors.As(err, &ve) {
		for _, fe := range ve {
			out[fe.Field()] = fe.Translate(v.translator)
		}
	} else if err != nil {
		out["unexpected_error"] = err.Error()
	}
	return out
}
