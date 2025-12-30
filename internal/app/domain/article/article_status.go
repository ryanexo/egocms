package article

import `dpcms/internal/app/erroz`

type Status int8

const (
    // StatusDraft 草稿状态
    StatusDraft Status = iota
    // StatusPending 待审核
    StatusPending
    // StatusPublished 已发布
    StatusPublished
    // StatusOffline 已下线
    StatusOffline
    // StatusReject 审核拒绝
    StatusReject
    // StatusPendingRepublish 编辑后等待重新审核
    StatusPendingRepublish
)

var stateMachine = map[Status]struct {
    From  Status
    Error error
}{
    StatusPending:          {From: StatusDraft, Error: erroz.ArticleSubmitStatusNotAllowed.ToError()},
    StatusPublished:        {From: StatusPending, Error: erroz.ArticlePublishStatusNotAllowed.ToError()},
    StatusOffline:          {From: StatusPublished, Error: erroz.ArticleOfflineStatusNotAllowed.ToError()},
    StatusReject:           {From: StatusPending, Error: erroz.ArticleRejectStatusNotAllowed.ToError()},
    StatusPendingRepublish: {From: StatusPublished, Error: erroz.ArticleRepublishStatusNotAllowed.ToError()},
}
