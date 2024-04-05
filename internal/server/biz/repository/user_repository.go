package repository

import (
    `context`

    `GoBlog/internal/server/entity`
    `gorm.io/gorm`
)

type UserRepository interface {
    Create(ctx context.Context, data *entity.User) error
    Find(context.Context, Fields, Conditions) (entity.User, error)
    FindByID(ctx context.Context, id uint) (entity.User, error)
}
type userRepo struct {
    db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
    return &userRepo{db}
}

func (r userRepo) Find(ctx context.Context, fields Fields, conditions Conditions) (result entity.User, err error) {
    q := r.db.WithContext(ctx).Model(result)
    if len(fields) > 0 {
        q = q.Select(fields)
    }
    if conditions != nil {
        q = q.Where(conditions)
    }
    err = q.First(&result).Error
    return
}

func (r userRepo) Create(ctx context.Context, data *entity.User) error {
    return r.db.WithContext(ctx).Create(data).Error
}

func (r userRepo) FindByID(ctx context.Context, id uint) (result entity.User, err error) {
    return r.Find(ctx, nil, Conditions{"id": id})
}

func (r userRepo) List(ctx context.Context, pageNo int, pageSize int) (result []entity.User, err error) {
    table := new(entity.User)
    subQuery := r.db.Table("(?) AS idx", r.db.Model(table)).Select("id").Limit(pageSize).Offset(pageNo * pageSize)
    err = r.db.WithContext(ctx).Model(table).Where("id EXISTS ?", subQuery).Find(&result).Error
    return
}

func (r userRepo) Count(ctx context.Context, cond map[string]any) (count int64) {
    r.db.WithContext(ctx).Model(new(entity.User)).Where(cond).Count(&count)
    return
}
