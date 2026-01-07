package contextutil

import (
    `context`
    
    `dpcms/internal/erroz`
    `dpcms/internal/constant`
    `dpcms/internal/infra/persistence/model`
)

func GetAuthorizedUser(ctx context.Context) (*model.User, error) {
    u, ok := ctx.Value(constant.RequestUserKey).(*model.User)
    if !ok {
        return nil, erroz.Unauthorized.ToError()
    }
    return u, nil
}
