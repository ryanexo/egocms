package validate

import (
    "reflect"
    `strings`
    
    "github.com/gin-gonic/gin/binding"
    `github.com/go-playground/locales/zh_Hans_CN`
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zhTranslation "github.com/go-playground/validator/v10/translations/zh"
)

var (
    translator ut.Translator
)

func GetTranslator() ut.Translator {
    return translator
}

func ReplaceTranslator(trans ut.Translator) func() {
    rollback := translator
    translator = trans
    return func() {
        translator = rollback
    }
}

type ErrorSlice []validator.ValidationErrors

func (s ErrorSlice) Error() string {
    errStr := strings.Builder{}
    for _, validationErrors := range s {
        errStr.WriteString(validationErrors.Error())
        errStr.WriteString("\n")
    }
    return errStr.String()
}

type customValidator struct {
    validate *validator.Validate
}

var _ binding.StructValidator = (*customValidator)(nil)

func (v *customValidator) ValidateStruct(s any) error {
    if s == nil {
        return nil
    }
    val := reflect.ValueOf(s)
    switch val.Kind() {
    case reflect.Ptr:
        if val.Elem().Kind() != reflect.Struct {
            return v.ValidateStruct(val.Elem().Interface())
        }
        fallthrough
    case reflect.Struct:
        return v.validateStruct(s)
    case reflect.Slice, reflect.Array:
        count := val.Len()
        errs := make(ErrorSlice, 0, count)
        for i := 0; i < count; i++ {
            err := v.ValidateStruct(val.Index(i).Interface())
            if err == nil {
                continue
            }
            errs = append(errs, err.(validator.ValidationErrors))
        }
        if len(errs) == 0 {
            return nil
        }
        return nil
    
    default:
        return nil
    }
}

func (v *customValidator) validateStruct(s any) error {
    var err error
    if obj, ok := s.(partialValidate); ok {
        err = v.validateStructPartial(obj)
    } else {
        err = v.validate.Struct(s)
    }
    return err
}

func (v *customValidator) validateStructPartial(s partialValidate) error {
    fields := s.ValidationFields()
    if len(fields) == 0 {
        return v.validate.Struct(s)
    }
    return v.validate.StructPartial(s, fields...)
}

func (v *customValidator) Engine() any {
    return v.validate
}

func New() (binding.StructValidator, error) {
    v := validator.New()
    v.SetTagName("validate")
    v.RegisterTagNameFunc(func(field reflect.StructField) string {
        key, _ := strings.CutSuffix(field.Tag.Get("json"), ",")
        if key != "" {
            return key
        }
        return field.Name
    })
    
    zhLang := zh_Hans_CN.New()
    uni := ut.New(zhLang)
    trans, _ := uni.GetTranslator("zh")
    err := zhTranslation.RegisterDefaultTranslations(v, trans)
    if err != nil {
        return nil, err
    }
    translator = trans
    
    return &customValidator{v}, nil
}
