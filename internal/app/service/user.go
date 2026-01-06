package service

import (
	"context"
	"dpcms/internal/infra/persistence/model"
	"dpcms/internal/infra/persistence/repo"

	usrAssembler "dpcms/internal/app/assembler/user"
	"dpcms/internal/app/dto"
	"dpcms/internal/app/erroz"
	"dpcms/internal/infra"
	"dpcms/internal/infra/password"
	"dpcms/internal/infra/persistence/datatype"
	"dpcms/internal/infra/persistence/query"
)

type User struct {
	persist *query.Query
}

func NewUserService(i *infra.Infra) *User {
	return &User{persist: i.Query}
}

func (srv User) Create(ctx context.Context, params dto.UserCreateParams) (datatype.SafeUint64, error) {
	usrRepo := repo.NewUserRepo(srv.persist)
	user, err := usrRepo.FirstByUsernameOrEmail(ctx, params.Username, params.Email)
	if err != nil {
		return 0, err
	}
	if user.Username == params.Username {
		return 0, erroz.UserNameExists.ToError()
	}
	if user.Email == params.Email {
		return 0, erroz.UserEmailExists.ToError()
	}

	hashedPwd, err := password.Password(params.Password).Generate()
	if err != nil {
		return 0, err
	}
	data := usrAssembler.BuildUserCreateCommand(&params)
	data.Password = hashedPwd

	if err = usrRepo.Create(ctx, data); err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (srv User) FindByCredential(ctx context.Context, params dto.UserCredentialParams) (*dto.User, error) {
	data, err := repo.NewUserRepo(srv.persist).FindByUsername(ctx, params.Username)
	if err != nil {
		return nil, err
	}
	if !password.Password(data.Password).Compare(params.Password) {
		return nil, erroz.UserWrongPasswd.ToError()
	}
	return usrAssembler.BuildUserDTO(data), nil
}

func (srv User) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.User, error) {
	data, err := repo.NewUserRepo(srv.persist).FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return usrAssembler.BuildUserDTO(data), nil
}

func (srv User) FindByName(ctx context.Context, name string) (*dto.User, error) {
	data, err := repo.NewUserRepo(srv.persist).FindByUsername(ctx, name)
	if err != nil {
		return nil, err
	}
	return usrAssembler.BuildUserDTO(data), nil
}

func (srv User) ResetPassword(ctx context.Context, id datatype.SafeUint64, pwd string) error {
	hashedPwd, err := password.Password(pwd).Generate()
	if err != nil {
		return err
	}
	_, err = repo.NewUserRepo(srv.persist).UpdatePassword(ctx, id, hashedPwd)
	return err
}

func (srv User) ChangePassword(ctx context.Context, current *model.User, params dto.UserPasswdUpdateParams) error {
	if params.Password != params.PasswordConfirm {
		return erroz.UserWrongConfirmPasswd.ToError()
	}
	if params.Password == params.OldPassword {
		return erroz.UserEqualsOldPasswd.ToError()
	}
	_, err := srv.FindByCredential(ctx, dto.UserCredentialParams{
		Username: current.Username,
		Password: params.OldPassword,
	})
	if err != nil {
		return err
	}
	return srv.ResetPassword(ctx, current.ID, params.Password)
}

func (srv User) Delete(ctx context.Context, id datatype.SafeUint64) error {
	return srv.persist.Transaction(func(tx *query.Query) error {
		return repo.NewUserRepo(tx).Delete(ctx, id)
	})
}

func (srv User) UpdateProfile(ctx context.Context, params dto.UserProfile) error {
	data := usrAssembler.BuildUserProfileModel(&params)
	_, err := repo.NewUserRepo(srv.persist).UpdateProfile(ctx, data)
	return err
}

func (srv User) List(ctx context.Context, params dto.UserListParams) (*dto.PaginatedResult[*dto.User], error) {
	data, total, err := repo.NewUserRepo(srv.persist).List(ctx, &params)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedResult[*dto.User]{
		PageNo:   params.PageNo,
		PageSize: params.PageSize,
		Total:    total,
		List:     usrAssembler.BuildUserListDTO(data),
	}, nil
}
