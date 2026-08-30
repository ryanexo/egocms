package user

import (
	"context"
	"errors"

	"cms/internal/app/user/api"
	"cms/internal/app/user/model"
	"cms/internal/infra/store/gorm/dbscope"
	"cms/internal/infra/store/gorm/gquery"

	"gorm.io/gen"
	"gorm.io/gorm"
)

type repo struct {
	q *gquery.Query
}

func NewUserRepo(q *gquery.Query) Repository {
	return &repo{q: q}
}

func (repo *repo) Create(ctx context.Context, user *model.User) error {
	return repo.q.Transaction(func(tx *gquery.Query) error {
		profile := user.Profile
		user.Profile = nil
		if err := tx.User.WithContext(ctx).Create(user); err != nil {
			user.Profile = profile
			return err
		}
		user.Profile = profile
		if profile == nil {
			return nil
		}
		profile.UserID = user.ID
		return tx.UserProfile.WithContext(ctx).Create(profile)
	})
}

func (repo *repo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	userQuery := repo.q.User
	user, err := userQuery.WithContext(ctx).Where(userQuery.Username.Eq(username)).First()
	if err != nil {
		return nil, err
	}
	return user, repo.hydrate(ctx, user)
}

func (repo *repo) FirstByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error) {
	userQuery := repo.q.User
	return userQuery.WithContext(ctx).
		Where(userQuery.Username.Eq(username)).
		Or(userQuery.Email.Eq(email)).
		First()
}

func (repo *repo) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	userQuery := repo.q.User
	user, err := userQuery.WithContext(ctx).Where(userQuery.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return user, repo.hydrate(ctx, user)
}

func (repo *repo) UpdatePassword(ctx context.Context, id uint64, passwordHash string) (gen.ResultInfo, error) {
	userQuery := repo.q.User
	return userQuery.WithContext(ctx).Where(userQuery.ID.Eq(id)).Update(userQuery.Password, passwordHash)
}

func (repo *repo) Delete(ctx context.Context, id uint64) error {
	return repo.q.Transaction(func(tx *gquery.Query) error {
		profileQuery := tx.UserProfile
		if _, err := profileQuery.WithContext(ctx).Where(profileQuery.UserID.Eq(id)).Delete(); err != nil {
			return err
		}
		userQuery := tx.User
		_, err := userQuery.WithContext(ctx).Where(userQuery.ID.Eq(id)).Delete()
		return err
	})
}

func (repo *repo) UpdateProfile(ctx context.Context, profile *model.UserProfile) (gen.ResultInfo, error) {
	profileQuery := repo.q.UserProfile
	return profileQuery.WithContext(ctx).
		Where(profileQuery.UserID.Eq(profile.UserID)).
		Select(
			profileQuery.Avatar,
			profileQuery.Nickname,
			profileQuery.Gender,
			profileQuery.Description,
			profileQuery.Country,
			profileQuery.Province,
			profileQuery.City,
		).
		Updates(profile)
}

func (repo *repo) List(ctx context.Context, params *api.UserListParams) ([]*model.User, int64, error) {
	userQuery := repo.q.User
	profileQuery := repo.q.UserProfile
	query := userQuery.WithContext(ctx).LeftJoin(profileQuery, profileQuery.UserID.EqCol(userQuery.ID))
	if params.Username != nil {
		query = query.Where(userQuery.Username.Like("%" + *params.Username + "%"))
	}
	if params.Email != nil {
		query = query.Where(userQuery.Email.Eq(*params.Email))
	}
	if params.Status != nil {
		query = query.Where(userQuery.Status.Eq(*params.Status))
	}
	if params.VerifyStartTime != nil {
		query = query.Where(userQuery.VerifiedAt.Gte(*params.VerifyStartTime))
	}
	if params.VerifyEndTime != nil {
		query = query.Where(userQuery.VerifiedAt.Lte(*params.VerifyEndTime))
	}
	if params.IP != nil {
		query = query.Where(userQuery.IP.Like(*params.IP))
	}
	if params.RoleID != nil {
		query = query.Where(userQuery.RoleID.Eq(params.RoleID.Uint64()))
	}
	if params.Nickname != nil {
		query = query.Where(profileQuery.Nickname.Eq(*params.Nickname))
	}
	if params.Gender != nil {
		query = query.Where(profileQuery.Gender.Eq(*params.Gender))
	}
	if params.Country != nil {
		query = query.Where(profileQuery.Country.Eq(*params.Country))
	}
	if params.Province != nil {
		query = query.Where(profileQuery.Province.Eq(*params.Province))
	}
	if params.City != nil {
		query = query.Where(profileQuery.City.Eq(*params.City))
	}

	total, err := query.Count()
	if err != nil {
		return nil, 0, err
	}
	users, err := query.Scopes(dbscope.Paginate(params.PageNo, params.PageSize)).Find()
	if err != nil {
		return nil, 0, err
	}
	for _, user := range users {
		if err = repo.hydrate(ctx, user); err != nil {
			return nil, 0, err
		}
	}
	return users, total, nil
}

func (repo *repo) hydrate(ctx context.Context, user *model.User) error {
	profileQuery := repo.q.UserProfile
	profile, err := profileQuery.WithContext(ctx).Where(profileQuery.UserID.Eq(user.ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		user.Profile = profile
	}
	if user.RoleID == nil {
		return nil
	}
	roleQuery := repo.q.Role
	role, err := roleQuery.WithContext(ctx).Where(roleQuery.ID.Eq(*user.RoleID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		user.RoleName = role.Name
	}
	return nil
}
