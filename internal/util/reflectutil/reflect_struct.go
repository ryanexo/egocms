package reflectutil

import (
    `errors`
    `reflect`
)

func InvokeImplementedStruct[T any](v any, callback func(reflect.Value, T) error) error {
    ref := reflect.ValueOf(v)
    
    if ref.Kind() == reflect.Ptr {
        ref = ref.Elem()
    }
    if ref.Kind() != reflect.Struct {
        return errors.New("v must be a struct")
    }
    
    for i := 0; i < ref.NumField(); i++ {
        iterateField := ref.Field(i)
        if !iterateField.CanInterface() {
            continue
        }
    INVOKE:
        controller, ok := iterateField.Interface().(T)
        if ok {
            err := callback(iterateField, controller)
            if err != nil {
                return err
            }
        } else if iterateField.Kind() == reflect.Ptr {
            iterateField = iterateField.Elem()
            goto INVOKE
        }
    }
    
    return nil
}
