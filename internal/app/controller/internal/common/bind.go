package common

import (
    `fmt`
    `reflect`
    
    `dpcms/internal/app/erroz`
    
    `github.com/gin-gonic/gin`
)

func BasicBind[T any](ctx *gin.Context, fn func(params T) (any, error)) {
    fnType := reflect.TypeOf(fn)
    if fnType.Kind() != reflect.Func {
        panic("logic only support func")
    }
    
    pType := fnType.In(0)
    zeroValue := reflect.Zero(pType)
    
    var params T
    
    switch zeroValue.Kind() {
    case reflect.Struct:
        params = zeroValue.Interface().(T)
    
    case reflect.Slice:
        params = reflect.MakeSlice(zeroValue.Type(), 0, 0).Interface().(T)
    
    default:
        panic(fmt.Errorf("unknown callback param type: %s", pType.Name()))
    }
    
    if err := ctx.ShouldBind(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    
    result, err := fn(params)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else if result == nil {
        erroz.OK.Write(ctx)
    } else {
        erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
    }
}
