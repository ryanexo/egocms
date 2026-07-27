package contract

import (
    "context"
    
    "cms/internal/infra/persistence"
    `cms/internal/infra/persistence/gorm/model`
    "cms/internal/pkg/datatype"
    
    "cms/internal/modules/user/internal/dto"
    
    "gorm.io/gen"
)

type UserRepo interface {
    persistence.Repository[UserRepo]
    Create(ctx context.Context, data *model.User) error
    FindByUsername(ctx context.Context, username string) (*model.User, error)
    FirstByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.User, error)
    UpdatePassword(ctx context.Context, id datatype.SafeUint64, passwd string) (gen.ResultInfo, error)
    Delete(ctx context.Context, id datatype.SafeUint64) error
    UpdateProfile(ctx context.Context, data *model.UserProfile) (gen.ResultInfo, error)
    List(ctx context.Context, params *dto.UserListParams) ([]*model.User, int64, error)
}
