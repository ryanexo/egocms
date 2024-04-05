package service

import (
    `context`
    `fmt`
    `testing`

    `GoBlog/internal/pkg/data`
    `GoBlog/internal/server/biz/repository`
    `GoBlog/internal/server/entity`
    `github.com/DATA-DOG/go-sqlmock`
    `github.com/stretchr/testify/assert`
)

func TestUser_Create(t *testing.T) {
    db, sqlm := data.NewMockDB()
    uname := "test_dev"
    email := "test@dev.com"
    repo := repository.NewUserRepo(db)
    srv := NewUserService(repo)
    sqlm.ExpectPrepare("^(?i)select .+? from `user`.+?where `username` = \\?").
        ExpectQuery().
        WillReturnRows(sqlmock.NewRows(nil))
    sqlm.ExpectPrepare("^(?i)select .+? from `user`.+?where `email` = \\?").
        ExpectQuery().
        WillReturnRows(sqlmock.NewRows(nil))

    sqlm.ExpectPrepare("^(?i)insert into `user`").
        ExpectExec().
        WillReturnResult(sqlmock.NewResult(1, 1))
    u := &entity.User{
        Username: uname,
        Email:    email,
    }
    err := srv.Create(context.Background(), u)
    assert.Nil(t, err)
    err = sqlm.ExpectationsWereMet()
    assert.Nil(t, err)
    assert.Equal(t, uint(1), u.ID)
}

func TestUser(t *testing.T) {
    x := []int{0, 0, 0}
    z := make([]int, 3)
    fmt.Println(len(x), len(z))
}
