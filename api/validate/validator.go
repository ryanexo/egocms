package validate

import (
    "reflect"
    `sync`
    
    "github.com/gin-gonic/gin/binding"
    `github.com/go-playground/locales/zh_Hans_CN`
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zhTranslation "github.com/go-playground/validator/v10/translations/zh"
    `github.com/google/wire`
)

var ProviderSet = wire.NewSet(New)

type PartialValidate interface {
    ValidationFields() []string
}

type customValidator struct {
    once       sync.Once
    validate   *validator.Validate
    translator ut.Translator
}

var _ binding.StructValidator = new(customValidator)

func (v *customValidator) lazyInit() error {
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
        trans, _ := uni.GetTranslator("zh_Hans_CN")
        err = zhTranslation.RegisterDefaultTranslations(validate, trans)
        
        if err != nil {
            return
        }
        
        v.validate = validate
        v.translator = trans
    })
    
    return err
}

func (v *customValidator) ValidateStruct(s any) error {
    err := v.lazyInit()
    if err != nil {
        return err
    }
    
    if obj, ok := s.(PartialValidate); ok {
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
        reason := e.Translate(v.translator)
        
        // if e.Tag() == "eqfield" {
        //     objType := reflect.TypeOf(s)
        //     if objType.Kind() == reflect.Ptr {
        //         objType = objType.Elem()
        //     }
        //     if objType.Kind() != reflect.Struct {
        //         goto AppendError
        //     }
        //     targetField := e.Param()
        //     if targetField == "" {
        //         goto AppendError
        //     }
        //     targetFieldType, found := objType.FieldByName(targetField)
        //     if !found {
        //         goto AppendError
        //     }
        //     targetFieldLabel := targetFieldType.Tag.Get("label")
        //     if targetFieldLabel != "" {
        //         reason = strings.Replace(reason, targetField, targetFieldLabel, 1)
        //     }
        // }
        
        // AppendError:
        finalErrors = append(finalErrors, ValidationError{Field: e.StructField(), Reason: reason})
    }
    
    return finalErrors
}

func (v *customValidator) Engine() any {
    err := v.lazyInit()
    if err != nil {
        return err
    }
    
    return v.validate
}

func New() binding.StructValidator {
    return &customValidator{}
}
