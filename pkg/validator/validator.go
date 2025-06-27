package validator

import (
    "fmt"

    "github.com/go-playground/locales/en"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Interface interface {
    Validate(data interface{}) error
}

type Validator struct {
    validator  *validator.Validate
    translator ut.Translator
}

func NewValidator() *Validator {
    v := validator.New()

    // Set up English translator
    eng := en.New()
    uni := ut.New(eng, eng)
    trans, _ := uni.GetTranslator("en")

    // Register default English translations
    enTranslations.RegisterDefaultTranslations(v, trans)

    // Register custom translations
    v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
        return ut.Add("required", "{0} is a required field", true)
    }, func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("required", fe.Field())
        return t
    })

    v.RegisterTranslation("min", trans, func(ut ut.Translator) error {
        return ut.Add("min", "{0} must be at least {1} characters long", true)
    }, func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("min", fe.Field(), fe.Param())
        return t
    })

    v.RegisterTranslation("max", trans, func(ut ut.Translator) error {
        return ut.Add("max", "{0} must not exceed {1} characters", true)
    }, func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("max", fe.Field(), fe.Param())
        return t
    })

    v.RegisterTranslation("numeric", trans, func(ut ut.Translator) error {
        return ut.Add("numeric", "{0} must contain only numeric characters", true)
    }, func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("numeric", fe.Field())
        return t
    })

    return &Validator{validator: v, translator: trans}
}

func (v *Validator) Validate(data interface{}) error {
    err := v.validator.Struct(data)
    if err != nil {
        if validationErrs, ok := err.(validator.ValidationErrors); ok {
            var errorMessages []string
            for _, e := range validationErrs {
                errorMessages = append(errorMessages, e.Translate(v.translator))
            }
            return fmt.Errorf("validation failed: %v", errorMessages)
        }
        return err
    }
    return nil
}