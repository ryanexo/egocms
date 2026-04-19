package domain

import (
    `cms/internal/domain/article/internal/errno`
)

type ArticleStatus struct {
    status StatusValue
    actor  Actor
}

type StatusValue = int8

const (
    // StatusDraft 草稿状态
    StatusDraft StatusValue = iota
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

type Actor struct {
    CanPublishDirect bool
}

func NewStatus(status StatusValue) *ArticleStatus {
    return &ArticleStatus{status: status}
}

func (s *ArticleStatus) WithActor(actor Actor) *ArticleStatus {
    s.actor = actor
    return s
}

func (s *ArticleStatus) transitionStatus(expect StatusValue) error {
    stateMachine := map[StatusValue]struct {
        From  StatusValue
        Error error
    }{
        StatusPending:          {From: StatusDraft, Error: errno.ArticleSubmitStatusNotAllowed.ToError()},
        StatusPublished:        {From: StatusPending, Error: errno.ArticlePublishStatusNotAllowed.ToError()},
        StatusOffline:          {From: StatusPublished, Error: errno.ArticleOfflineStatusNotAllowed.ToError()},
        StatusReject:           {From: StatusPending, Error: errno.ArticleRejectStatusNotAllowed.ToError()},
        StatusPendingRepublish: {From: StatusPublished, Error: errno.ArticleRepublishStatusNotAllowed.ToError()},
    }
    
    nextState := stateMachine[expect]
    if nextState.From != s.status {
        return nextState.Error
    }
    s.status = expect
    return nil
}

func (s *ArticleStatus) Submit() error {
    if s.actor.CanPublishDirect {
        s.status = StatusPublished
        return nil
    }
    if s.status == StatusPending {
        return errno.ArticleAlreadySubmitted.ToError()
    }
    return s.transitionStatus(StatusPending)
}

func (s *ArticleStatus) Publish() error {
    if !s.actor.CanPublishDirect {
        if s.status == StatusPublished {
            return errno.ArticleAlreadyPublished.ToError()
        }
    }
    return s.transitionStatus(StatusPublished)
}

func (s *ArticleStatus) Offline() error {
    if s.status == StatusOffline {
        return errno.ArticleAlreadyOffline.ToError()
    }
    return s.transitionStatus(StatusOffline)
}

func (s *ArticleStatus) Reject() error {
    if s.status == StatusReject {
        return errno.ArticleAlreadyReject.ToError()
    }
    return s.transitionStatus(StatusReject)
}

func (s *ArticleStatus) Republish() error {
    if !s.actor.CanPublishDirect {
        if s.status == StatusPendingRepublish {
            return errno.ArticleAlreadyRepublish.ToError()
        }
    }
    return s.transitionStatus(StatusPendingRepublish)
}

func (s *ArticleStatus) Value() StatusValue {
    return s.status
}
