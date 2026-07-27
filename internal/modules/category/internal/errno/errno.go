package errno

import (
    `cms/internal/httpx`
)

var (
    ErrCircular   = httpx.NewError(1201, "目标分类 %s 为当前分类 %s 的子级，无法移动")
    ErrPathFormat = httpx.NewError(1202, "分类路径格式不正确")
)
