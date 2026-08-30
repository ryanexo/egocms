package errno

import `cms/internal/public/erroz`

var (
    ErrInvalidNameLength        = erroz.NewError("CATEGORY_001", "分类名称长度超出限制")
    ErrInvalidPath              = erroz.NewError("CATEGORY_002", "分类路径仅允许数字、字母以及-_，且以数字、字母开头")
    ErrNameRequired             = erroz.NewError("CATEGORY_003", "分类名称必填")
    ErrPathRequired             = erroz.NewError("CATEGORY_004", "分类路径必填")
    ErrTitleInvalidLength       = erroz.NewError("CATEGORY_005", "分类SEO标题长度超出限制")
    ErrKeywordsInvalidLength    = erroz.NewError("CATEGORY_006", "分类SEO关键词长度超出限制")
    ErrDescriptionInvalidLength = erroz.NewError("CATEGORY_007", "分类SEO描述长度超出限制")
    ErrThumbInvalidLength       = erroz.NewError("CATEGORY_008", "分类缩略图URL长度超出限制")
    ErrThumbInvalidPath         = erroz.NewError("CATEGORY_009", "分类缩略图URL不合法")
)
