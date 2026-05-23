package domain

import (
	"cms/internal/app/article/internal/errno"
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
		StatusPending:          {From: StatusDraft, Error: errno.ErrSubmitStatusNotAllowed},
		StatusPublished:        {From: StatusPending, Error: errno.ErrPublishStatusNotAllowed},
		StatusOffline:          {From: StatusPublished, Error: errno.ErrOfflineStatusNotAllowed},
		StatusReject:           {From: StatusPending, Error: errno.ErrRejectStatusNotAllowed},
		StatusPendingRepublish: {From: StatusPublished, Error: errno.ErrRepublishStatusNotAllowed},
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
		return errno.ErrAlreadySubmitted
	}
	return s.transitionStatus(StatusPending)
}

func (s *Status) Publish() error {
	if !s.actor.CanPublishDirect {
		if s.status == StatusPublished {
			return errno.ErrAlreadyPublished
		}
	}
	return s.transitionStatus(StatusPublished)
}

func (s *Status) Offline() error {
	if s.status == StatusOffline {
		return errno.ErrAlreadyOffline
	}
	return s.transitionStatus(StatusOffline)
}

func (s *Status) Reject() error {
	if s.status == StatusReject {
		return errno.ErrAlreadyReject
	}
	return s.transitionStatus(StatusReject)
}

func (s *Status) Republish() error {
	if !s.actor.CanPublishDirect {
		if s.status == StatusPendingRepublish {
			return errno.ErrAlreadyRepublish
		}
	}
	return s.transitionStatus(StatusPendingRepublish)
}

func (s *Status) Value() StatusValue {
	return s.status
}
