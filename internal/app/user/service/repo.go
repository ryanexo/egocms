package service

import (
    `context`
    
    `cms/internal/app/user/internal/dto`
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type UserRepo interface {
    contract.Repository[UserRepo]
    Create(ctx context.Context, data *model.User) error
    FindByUsername(ctx context.Context, username string) (*model.User, error)
    FirstByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.User, error)
    UpdatePassword(ctx context.Context, id datatype.SafeUint64, passwd string) (gen.ResultInfo, error)
    Delete(ctx context.Context, id datatype.SafeUint64) error
    UpdateProfile(ctx context.Context, data *model.UserProfile) (gen.ResultInfo, error)
    List(ctx context.Context, params *dto.UserListParams) ([]*model.User, int64, error)
}
