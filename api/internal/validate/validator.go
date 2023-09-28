package validate

import (
	"errors"
	"fmt"

	localeus "github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	globalValidator  *validator.Validate
	globalTranslator ut.Translator
)

func init() {
	translator := ut.New(localeus.New())

	validate := validator.New()
	if err := translations.RegisterDefaultTranslations(validate, translator.GetFallback()); err != nil {
		panic(fmt.Errorf("registering translator: %w", err))
	}

	globalValidator = validate
	globalTranslator = translator.GetFallback()
}

func Struct(s any) error {
	if err := globalValidator.Struct(s); err != nil {
		verrs := err.(validator.ValidationErrors)
		return errors.New(verrs[0].Translate(globalTranslator))
	}

	return nil
}
