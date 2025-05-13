package repository

import (
    `context`
    `errors`
    `time`
    
    `dpcms/model`
    `dpcms/model/query`
    `gorm.io/gen/field`
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type Token struct {
    query *query.Query
}

func NewTokenRepo(db *gorm.DB) *Token {
    return &Token{query: query.Use(db)}
}

func (repo *Token) IsRevoked(ctx context.Context, userId uint, uuid string, expires int) (bool, error) {
    isRevoked := true
    q := repo.query.WithContext(ctx).TokenBlacklist
    err := repo.query.Transaction(func(tx *query.Query) error {
        _, err := q.Clauses(clause.Locking{Strength: "UPDATE"}).Where(repo.query.TokenBlacklist.UserId.Eq(userId)).Select(field.NewUnsafeFieldRaw("1")).Find()
        if err != nil {
            return err
        }
        token, err := q.Where(repo.query.TokenBlacklist.UserId.Eq(userId), repo.query.TokenBlacklist.UUID.Eq(uuid)).First()
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                isRevoked = false
                return nil
            }
            return err
        }
        exp := token.CreatedAt.Add(time.Duration(expires) * time.Second)
        if exp.Before(time.Now()) {
            return nil
        }
        _, err = q.Delete(token)
        if err == nil {
            isRevoked = false
        }
        return err
    })
    if err != nil {
        return isRevoked, err
    }
    return isRevoked, nil
}

func (repo *Token) Revoke(ctx context.Context, userId uint, uuid string) error {
    return repo.query.WithContext(ctx).TokenBlacklist.Create(&model.TokenBlacklist{
        UserId: userId,
        UUID:   uuid,
    })
}
