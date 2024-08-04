package validate

import (
    `strings`
    
    `github.com/go-playground/validator/v10`
)

func ErrorToMap(errs validator.ValidationErrors) map[string]string {
    errorMap := make(map[string]string, len(errs))
    for _, err := range errs {
        /* TODO: 翻译临时使用registerTagNameFunc + StructField，需要在后续版本改为将字段翻译注册至翻译器 */
        key := strings.ToLower(err.StructField())
        errorMap[key] = err.Translate(translator)
    }
    return errorMap
}

func ErrorSliceToMap(errSlice ErrorSlice) []map[string]string {
    errorMaps := make([]map[string]string, len(errSlice))
    for _, errs := range errSlice {
        errorMaps = append(errorMaps, ErrorToMap(errs))
    }
    return errorMaps
}
