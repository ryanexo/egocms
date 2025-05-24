package category_error

import (
    `dpcms/internal/erroz`
)

var (
    ErrCircularReferenceWhenMove = erroz.New("category.circular", "目标分类 %s 为当前分类 %s 的子级")
)
