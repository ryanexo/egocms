package repo

import (
	"context"
	"database/sql"
	"dpcms/internal/app/dto"
	"dpcms/internal/infra/persistence/datatype"
	"dpcms/internal/infra/persistence/dbscope"
	"dpcms/internal/infra/persistence/model"
	"dpcms/internal/infra/persistence/query"

	"gorm.io/gen"
)

type UserRepo struct {
	persist *query.Query
}

func NewUserRepo(persist *query.Query) UserRepo {
	return UserRepo{persist}
}

func (r UserRepo) Create(ctx context.Context, data *model.User) error {
	return r.persist.User.WithContext(ctx).Create(data)
}

func (r UserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	dao := r.persist.User
	return dao.WithContext(ctx).Where(dao.Username.Eq(username)).Preload(dao.Profile, dao.Role.Select(r.persist.Role.Name)).First()
}

func (r UserRepo) FirstByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error) {
	dao := r.persist.User
	return dao.WithContext(ctx).Where(dao.Username.Eq(username)).Or(dao.Email.Eq(email)).First()
}

func (r UserRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.User, error) {
	dao := r.persist.User
	return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
}

func (r UserRepo) UpdatePassword(ctx context.Context, id datatype.SafeUint64, passwd string) (gen.ResultInfo, error) {
	dao := r.persist.User
	return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Update(dao.Password, passwd)
}

func (r UserRepo) Delete(ctx context.Context, id datatype.SafeUint64) error {
	ud := r.persist.User
	_, err := ud.WithContext(ctx).Where(ud.ID.Eq(id.Raw())).Delete()
	if err != nil {
		return err
	}
	pd := r.persist.UserProfile
	_, err = pd.WithContext(ctx).Where(pd.UserID.Eq(id.Raw())).Delete()
	return err
}

func (r UserRepo) UpdateProfile(ctx context.Context, data *model.UserProfile) (gen.ResultInfo, error) {
	dao := r.persist.UserProfile
	return dao.WithContext(ctx).Where(dao.UserID.Eq(data.UserID.Raw())).Updates(data)
}

func (r UserRepo) List(ctx context.Context, params *dto.UserListParams) ([]*model.User, int64, error) {
	uo := r.persist.User
	po := r.persist.UserProfile
	q := uo.WithContext(ctx).LeftJoin(po, po.UserID.EqCol(uo.ID))
	if params.Username != nil {
		q = q.Where(uo.Username.Like("%" + *params.Username + "%"))
	}
	if params.Email != nil {
		q = q.Where(uo.Email.Eq(*params.Email))
	}
	if params.Status != nil {
		q = q.Where(uo.Status.Eq(*params.Status))
	}
	if params.VerifyStartTime != nil {
		q = q.Where(
			uo.VerifiedAt.Gte(sql.NullTime{
				Valid: true,
				Time:  *params.VerifyStartTime,
			}),
		)
	}
	if params.VerifyEndTime != nil {
		q = q.Where(
			uo.VerifiedAt.Lte(sql.NullTime{
				Valid: true,
				Time:  *params.VerifyEndTime,
			}),
		)
	}
	if params.IP != nil {
		q = q.Where(uo.IP.Like(*params.IP))
	}
	if params.RoleID != nil {
		q = q.Where(uo.RoleID.Eq(params.RoleID.Raw()))
	}
	if params.Nickname != nil {
		q = q.Where(po.Nickname.Eq(*params.Nickname))
	}
	if params.Gender != nil {
		q = q.Where(po.Gender.Eq(*params.Gender))
	}
	if params.Country != nil {
		q = q.Where(po.Country.Eq(*params.Country))
	}
	if params.Province != nil {
		q = q.Where(po.Province.Eq(*params.Province))
	}
	if params.City != nil {
		q = q.Where(po.City.Eq(*params.City))
	}

	count, err := q.Count()
	if err != nil {
		return nil, 0, err
	}
	users, err := q.Scopes(dbscope.Paginate(params.PageNo, params.PageSize)).Find()
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
