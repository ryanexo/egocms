package user

import (
    "context"
    "errors"
    
    "cms/internal/app/user/api"
    "cms/internal/app/user/errno"
    "cms/internal/app/user/model"
    "cms/internal/public/apitype"
    "cms/internal/public/jsontype"
    "cms/internal/public/utils"
    
    "gorm.io/gen"
    "gorm.io/gorm"
)

type Repository interface {
    Create(ctx context.Context, user *model.User) error
    FindByUsername(ctx context.Context, username string) (*model.User, error)
    FirstByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error)
    FindByID(ctx context.Context, id uint64) (*model.User, error)
    UpdatePassword(ctx context.Context, id uint64, passwordHash string) (gen.ResultInfo, error)
    Delete(ctx context.Context, id uint64) error
    UpdateProfile(ctx context.Context, profile *model.UserProfile) (gen.ResultInfo, error)
    List(ctx context.Context, params *api.UserListParams) ([]*model.User, int64, error)
}

type Service struct {
    repo Repository
}

func NewUserService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (service *Service) Create(ctx context.Context, params api.UserCreateParams) (jsontype.SafeUint64, error) {
    if params.Password != params.PasswordConfirm {
        return 0, errno.ErrPasswordMismatch
    }
    
    existing, err := service.repo.FirstByUsernameOrEmail(ctx, params.Username, params.Email)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return 0, err
    }
    if existing != nil {
        if existing.Username == params.Username {
            return 0, errno.ErrUsernameExists
        }
        if existing.Email == params.Email {
            return 0, errno.ErrEmailExists
        }
    }
    
    user, err := buildUser(params)
    if err != nil {
        return 0, err
    }
    if err = service.repo.Create(ctx, user); err != nil {
        return 0, err
    }
    return jsontype.SafeUint64(user.ID), nil
}

func (service *Service) FindByCredential(ctx context.Context, params api.UserCredentialParams) (*api.User, error) {
    user, err := service.repo.FindByUsername(ctx, params.Username)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errno.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    if !utils.Password(user.Password).Compare(params.Password) {
        return nil, errno.ErrWrongPassword
    }
    return toUserAPI(user), nil
}

func (service *Service) FindModelByID(ctx context.Context, id uint64) (*model.User, error) {
    user, err := service.repo.FindByID(ctx, id)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errno.ErrUserNotFound
    }
    return user, err
}

func (service *Service) FindByID(ctx context.Context, id uint64) (*api.User, error) {
    user, err := service.FindModelByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return toUserAPI(user), nil
}

func (service *Service) FindByName(ctx context.Context, username string) (*api.User, error) {
    user, err := service.repo.FindByUsername(ctx, username)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errno.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    return toUserAPI(user), nil
}

func (service *Service) ResetPassword(ctx context.Context, id uint64, password string) error {
    if _, err := service.FindModelByID(ctx, id); err != nil {
        return err
    }
    passwordHash, err := utils.Password(password).Generate()
    if err != nil {
        return err
    }
    _, err = service.repo.UpdatePassword(ctx, id, passwordHash)
    return err
}

func (service *Service) ChangePassword(ctx context.Context, current *model.User, params api.UserPasswordUpdateParams) error {
    if params.Password != params.PasswordConfirm {
        return errno.ErrPasswordMismatch
    }
    if params.Password == params.OldPassword {
        return errno.ErrPasswordUnchanged
    }
    if !utils.Password(current.Password).Compare(params.OldPassword) {
        return errno.ErrWrongPassword
    }
    return service.ResetPassword(ctx, current.ID, params.Password)
}

func (service *Service) Delete(ctx context.Context, id uint64) error {
    if _, err := service.FindModelByID(ctx, id); err != nil {
        return err
    }
    return service.repo.Delete(ctx, id)
}

func (service *Service) UpdateProfile(ctx context.Context, params api.UserProfileParams) error {
    if _, err := service.FindModelByID(ctx, params.UserID.Uint64()); err != nil {
        return err
    }
    profile, err := buildProfile(params)
    if err != nil {
        return err
    }
    _, err = service.repo.UpdateProfile(ctx, profile)
    return err
}

func (service *Service) List(ctx context.Context, params api.UserListParams) (*apitype.PaginatedResult[*api.User], error) {
    users, total, err := service.repo.List(ctx, &params)
    if err != nil {
        return nil, err
    }
    return &apitype.PaginatedResult[*api.User]{
        Pagination: params.Pagination,
        Total:      total,
        List:       toUserAPIList(users),
    }, nil
}
