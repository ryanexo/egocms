package errmod

import (
    `dpcms/internal/app/erroz/internal/errcode`
)

type Code int64

func (e Code) String() string {
    return errcode.Code(e).String(3)
}

const (
    Server Code = iota
    Client
    User
    Role
    Category
    Menu
    Article
    ArticleModel
)
