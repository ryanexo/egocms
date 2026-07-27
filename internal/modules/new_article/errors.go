package new_article

import (
    "errors"
    
    "cms/internal/httpx"
    "cms/internal/modules/new_article/domain"
)

var (
    ErrInvalidStatusTransition = httpx.NewError(1101, "当前状态不允许此操作")
    ErrMissingCurrentVersion   = httpx.NewError(1102, "文章缺少当前版本")
    ErrNoChangesToRepublish    = httpx.NewError(1103, "文章没有需要重新发布的修改")
    ErrInvalidArticleContent   = httpx.NewError(1104, "文章内容不合法：%s")
    ErrInvalidArticleRelation  = httpx.NewError(1105, "文章分类或标签不合法")
    ErrInvalidContentTypeData  = httpx.NewError(1106, "内容类型数据不合法：%s")
    ErrContentTypeRequired     = httpx.NewError(1107, "提交内容类型数据前必须选择内容类型")
)

func mapDomainError(err error) error {
    switch {
    case err == nil:
        return nil
    case errors.Is(err, domain.ErrInvalidTransition), errors.Is(err, domain.ErrInvalidStatus):
        return ErrInvalidStatusTransition
    case errors.Is(err, domain.ErrMissingCurrentVersion):
        return ErrMissingCurrentVersion
    case errors.Is(err, domain.ErrNoChangesToRepublish):
        return ErrNoChangesToRepublish
    case errors.Is(err, domain.ErrTitleRequired):
        return ErrInvalidArticleContent.Format("标题不能为空")
    case errors.Is(err, domain.ErrTitleTooLong):
        return ErrInvalidArticleContent.Format("标题不能超过 255 个字符")
    case errors.Is(err, domain.ErrSummaryTooLong):
        return ErrInvalidArticleContent.Format("摘要不能超过 500 个字符")
    case errors.Is(err, domain.ErrContentTooLong):
        return ErrInvalidArticleContent.Format("正文不能超过 65535 字节")
    case errors.Is(err, domain.ErrChangeLogTooLong):
        return ErrInvalidArticleContent.Format("修改说明不能超过 500 个字符")
    default:
        return err
    }
}
