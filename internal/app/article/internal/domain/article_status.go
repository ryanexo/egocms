package domain

import (
    "dpcms/internal/erroz"
)

type Status struct {
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

func NewStatus(status StatusValue) *Status {
    return &Status{status: status}
}

func (s *Status) WithActor(actor Actor) *Status {
    s.actor = actor
    return s
}

func (s *Status) transitionStatus(expect StatusValue) error {
    stateMachine := map[StatusValue]struct {
        From  StatusValue
        Error error
    }{
        StatusPending:          {From: StatusDraft, Error: erroz.ArticleSubmitStatusNotAllowed.ToError()},
        StatusPublished:        {From: StatusPending, Error: erroz.ArticlePublishStatusNotAllowed.ToError()},
        StatusOffline:          {From: StatusPublished, Error: erroz.ArticleOfflineStatusNotAllowed.ToError()},
        StatusReject:           {From: StatusPending, Error: erroz.ArticleRejectStatusNotAllowed.ToError()},
        StatusPendingRepublish: {From: StatusPublished, Error: erroz.ArticleRepublishStatusNotAllowed.ToError()},
    }
    
    nextState := stateMachine[expect]
    if nextState.From != s.status {
        return nextState.Error
    }
    s.status = expect
    return nil
}

func (s *Status) Submit() error {
    if s.actor.CanPublishDirect {
        s.status = StatusPublished
        return nil
    }
    if s.status == StatusPending {
        return erroz.ArticleAlreadySubmitted.ToError()
    }
    return s.transitionStatus(StatusPending)
}

func (s *Status) Publish() error {
    if !s.actor.CanPublishDirect {
        if s.status == StatusPublished {
            return erroz.ArticleAlreadyPublished.ToError()
        }
    }
    return s.transitionStatus(StatusPublished)
}

func (s *Status) Offline() error {
    if s.status == StatusOffline {
        return erroz.ArticleAlreadyOffline.ToError()
    }
    return s.transitionStatus(StatusOffline)
}

func (s *Status) Reject() error {
    if s.status == StatusReject {
        return erroz.ArticleAlreadyReject.ToError()
    }
    return s.transitionStatus(StatusReject)
}

func (s *Status) Republish() error {
    if !s.actor.CanPublishDirect {
        if s.status == StatusPendingRepublish {
            return erroz.ArticleAlreadyRepublish.ToError()
        }
    }
    return s.transitionStatus(StatusPendingRepublish)
}

func (s *Status) Value() StatusValue {
    return s.status
}
