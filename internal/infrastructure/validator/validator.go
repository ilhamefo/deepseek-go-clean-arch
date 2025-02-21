package validator

import (
	"errors"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/go-playground/locales/en"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	*validator.Validate
	translator ut.Translator
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		panic("translator not found")
	}

	validate := validator.New()

	return &Validator{
		Validate:   validate,
		translator: trans,
	}

}

func (v *Validator) ValidationErrors(err error) (errMsg map[string]string) {
	var ve validator.ValidationErrors
	out := make(map[string]string)
	if errors.As(err, &ve) {
		for _, fe := range ve {
			out[fe.Field()] = v.getErrorMsg(fe)
		}
	} else {
		out["unexpected_error"] = err.Error()
	}
	return out
}

func (v *Validator) getErrorMsg(fe validator.FieldError) string {
	trans := v.translator

	en_translations.RegisterDefaultTranslations(v.Validate, trans)

	return fe.Translate(trans)
}
