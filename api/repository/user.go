package repository

import (
    "context"
    
    `dpcms/model`
    `dpcms/model/query`
    "gorm.io/gen"
    "gorm.io/gorm"
)

type User struct {
    query *query.Query
}

func NewUserRepo(db *gorm.DB) *User {
    return &User{query: query.Use(db)}
}

func (repo User) Create(ctx context.Context, data *model.User) error {
    return repo.query.User.WithContext(ctx).Create(data)
}

func (repo User) FindByID(ctx context.Context, id uint) (*model.User, error) {
    q := repo.query.User
    return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

func (repo User) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    q := repo.query.User
    return q.WithContext(ctx).Where(q.Username.Eq(username)).First()
}

func (repo User) FindByUsernameOrEmail(ctx context.Context, username string, email string) (*model.User, error) {
    q := repo.query.User
    return q.WithContext(ctx).Where(q.Username.Eq(username)).Or(q.Email.Eq(email)).First()
}

func (repo User) List(ctx context.Context, pageNo int, pageSize int) ([]*model.User, error) {
    q := repo.query.User
    userIDSet := q.WithContext(ctx).Select(q.ID).Limit(pageSize).Offset(pageSize * (pageNo - 1))
    return q.WithContext(ctx).Where(gen.Exists(userIDSet)).Find()
}

func (repo User) UpdateUserInfo(ctx context.Context, id uint, data map[string]any) (gen.ResultInfo, error) {
    q := repo.query.User
    return q.WithContext(ctx).Where(q.ID.Eq(id)).Updates(data)
}
