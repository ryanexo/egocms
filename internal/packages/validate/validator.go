package validate

import (
    "reflect"
    `sync`
    
    "github.com/gin-gonic/gin/binding"
    `github.com/go-playground/locales/zh_Hans_CN`
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    "github.com/go-playground/validator/v10/translations/zh"
    `github.com/google/wire`
)

var ProviderSet = wire.NewSet(New)

type PartialValidation interface {
    ValidationFields() []string
}

type validatorEx struct {
    once       sync.Once
    validate   *validator.Validate
    translator ut.Translator
}

var _ binding.StructValidator = new(validatorEx)

func (v *validatorEx) lazyInit() error {
    var err error
    
    v.once.Do(func() {
        validate := validator.New()
        validate.SetTagName("validate")
        validate.RegisterTagNameFunc(func(field reflect.StructField) string {
            key := field.Tag.Get("label")
            if key == "" {
                return field.Name
            } else {
                return key
            }
        })
        
        locale := zh_Hans_CN.New()
        uni := ut.New(locale)
        translator, _ := uni.GetTranslator("zh_Hans_CN")
        err = zh.RegisterDefaultTranslations(validate, translator)
        
        if err != nil {
            return
        }
        
        v.validate = validate
        v.translator = translator
    })
    
    return err
}

func (v *validatorEx) ValidateStruct(s any) error {
    err := v.lazyInit()
    if err != nil {
        return err
    }
    
    if obj, ok := s.(PartialValidation); ok {
        fields := obj.ValidationFields()
        err = v.validate.StructPartial(s, fields...)
    } else {
        err = v.validate.Struct(s)
    }
    
    if err == nil {
        return nil
    }
    
    errs := err.(validator.ValidationErrors)
    finalErrors := make(ValidationErrors, 0, len(errs))
    
    for _, e := range errs {
        currentError := ValidationError{Field: e.StructField()}
        if msg, ok := s.(CustomValidationMessage); ok {
            reason, found := msg.ValidationMessage(e)
            if found {
                currentError.Reason = reason
            }
        }
        
        if currentError.Reason == "" {
            currentError.Reason = e.Translate(v.translator)
        }
        finalErrors = append(finalErrors, currentError)
    }
    
    return finalErrors
}

func (v *validatorEx) Engine() any {
    err := v.lazyInit()
    if err != nil {
        return err
    }
    
    return v.validate
}

func New() binding.StructValidator {
    return &validatorEx{}
}
