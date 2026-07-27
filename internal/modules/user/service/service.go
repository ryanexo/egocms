package service

import (
    "context"
    "errors"
    
    "cms/internal/infra/persistence"
    `cms/internal/infra/persistence/gorm/model`
    `cms/internal/infra/persistence/gorm/gquery`
    "cms/internal/pkg/datatype"
    "cms/internal/pkg/password"
    
    contract2 "cms/internal/modules/user/contract"
    "cms/internal/modules/user/internal/assembler"
    "cms/internal/modules/user/internal/dto"
    "cms/internal/modules/user/internal/errno"
    "cms/internal/util/types"
    
    "gorm.io/gorm"
)

type UserService struct {
    txManager persistence.Transactor
    repo      contract2.UserRepo
}

func NewUserService(txManager persistence.Transactor, repo contract2.UserRepo) *UserService {
    return &UserService{
        txManager: txManager,
        repo:      repo,
    }
}

func (s UserService) Create(ctx context.Context, params dto.UserCreateParams) (datatype.SafeUint64, error) {
    user, err := s.repo.FirstByUsernameOrEmail(ctx, params.Username, params.Email)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return 0, err
    }
    if user != nil {
        if user.Username == params.Username {
            return 0, errno.UserNameExists.ToError()
        }
        if user.Email == params.Email {
            return 0, errno.UserEmailExists.ToError()
        }
    }
    
    hashedPwd, err := password.Password(params.Password).Generate()
    if err != nil {
        return 0, err
    }
    data := assembler.ToUserCreateCommand(&params)
    data.Password = hashedPwd
    
    err = s.txManager.Transaction(func(tx *gquery.Query) error {
        return s.repo.CloneWithQuery(tx).Create(ctx, data)
    })
    if err != nil {
        return 0, err
    }
    
    return data.ID, nil
}

func (s UserService) FindByCredential(ctx context.Context, params dto.UserCredentialParams) (*dto.User, error) {
    data, err := s.repo.FindByUsername(ctx, params.Username)
    if err != nil {
        return nil, err
    }
    if !password.Password(data.Password).Compare(params.Password) {
        return nil, errno.UserWrongPasswd.ToError()
    }
    return assembler.ToUserDTO(data), nil
}

func (s UserService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.User, error) {
    data, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToUserDTO(data), nil
}

func (s UserService) FindByName(ctx context.Context, name string) (*dto.User, error) {
    data, err := s.repo.FindByUsername(ctx, name)
    if err != nil {
        return nil, err
    }
    return assembler.ToUserDTO(data), nil
}

func (s UserService) ResetPassword(ctx context.Context, id datatype.SafeUint64, pwd string) error {
    hashedPwd, err := password.Password(pwd).Generate()
    if err != nil {
        return err
    }
    _, err = s.repo.UpdatePassword(ctx, id, hashedPwd)
    return err
}

func (s UserService) ChangePassword(ctx context.Context, current *model.User, params dto.UserPasswdUpdateParams) error {
    if params.Password != params.PasswordConfirm {
        return errno.UserWrongConfirmPasswd.ToError()
    }
    if params.Password == params.OldPassword {
        return errno.UserEqualsOldPasswd.ToError()
    }
    _, err := s.FindByCredential(ctx, dto.UserCredentialParams{
        Username: current.Username,
        Password: params.OldPassword,
    })
    if err != nil {
        return err
    }
    return s.ResetPassword(ctx, current.ID, params.Password)
}

func (s UserService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *gquery.Query) error {
        return s.repo.CloneWithQuery(tx).Delete(ctx, id)
    })
}

func (s UserService) UpdateProfile(ctx context.Context, params dto.UserProfile) error {
    data := assembler.ToUserProfileModel(&params)
    _, err := s.repo.UpdateProfile(ctx, data)
    return err
}

func (s UserService) List(ctx context.Context, params dto.UserListParams) (*types.PaginatedResult[*dto.User], error) {
    data, total, err := s.repo.List(ctx, &params)
    if err != nil {
        return nil, err
    }
    
    return &types.PaginatedResult[*dto.User]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.ToUserListDTO(data),
    }, nil
}
