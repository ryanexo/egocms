package article

import `dpcms/internal/app/erroz`

type Status int8

const (
    // Draft 草稿状态
    Draft Status = iota
    // Pending 待审核
    Pending
    // Published 已发布
    Published
    // Offline 已下线
    Offline
    // Reject 审核拒绝
    Reject
    // PendingRepublish 编辑后等待重新审核
    PendingRepublish
)

var stateMachine = map[Status]struct {
    From  Status
    Error error
}{
    Pending:          {From: Draft, Error: erroz.ArticleSubmitStatusNotAllowed.ToError()},
    Published:        {From: Pending, Error: erroz.ArticlePublishStatusNotAllowed.ToError()},
    Offline:          {From: Published, Error: erroz.ArticleOfflineStatusNotAllowed.ToError()},
    Reject:           {From: Pending, Error: erroz.ArticleRejectStatusNotAllowed.ToError()},
    PendingRepublish: {From: Published, Error: erroz.ArticleRepublishStatusNotAllowed.ToError()},
}
