package article

import (
    `dpcms/internal/app/erroz`
)

type Draft struct {
    Status      int8
    Description string
    Content     string
}

type Article struct {
    actor       Actor
    status      Status
    description string
    content     string
}

type Actor struct {
    CanSubmit        bool
    CanPublish       bool
    CanPublishDirect bool
    CanOffline       bool
    CanReject        bool
}

func NewArticle(data Draft) *Article {
    return &Article{
        status:      Status(data.Status),
        content:     data.Content,
        description: data.Description,
    }
}

func (s *Article) WithActor(actor Actor) {
    s.actor = actor
}

func (s *Article) transitionStatus(expect Status) error {
    nextState := stateMachine[expect]
    if nextState.From != s.status {
        return nextState.Error
    }
    s.status = expect
    return nil
}

func (s *Article) Submit() error {
    if s.actor.CanPublishDirect {
        s.status = StatusPublished
        return nil
    }
    if !s.actor.CanSubmit {
        return erroz.ArticleMissingSubmitPerm.ToError()
    }
    if s.status == StatusPending {
        return erroz.ArticleAlreadySubmitted.ToError()
    }
    return s.transitionStatus(StatusPending)
}

func (s *Article) Publish() error {
    if !s.actor.CanPublishDirect {
        if !s.actor.CanPublish {
            return erroz.ArticleMissingPublishPerm.ToError()
        }
        if s.status == StatusPublished {
            return erroz.ArticleAlreadyPublished.ToError()
        }
    }
    return s.transitionStatus(StatusPublished)
}

func (s *Article) Offline() error {
    if !s.actor.CanOffline {
        return erroz.ArticleMissingOfflinePerm.ToError()
    }
    if s.status == StatusOffline {
        return erroz.ArticleAlreadyOffline.ToError()
    }
    return s.transitionStatus(StatusOffline)
}

func (s *Article) Reject() error {
    if !s.actor.CanReject {
        return erroz.ArticleMissingRejectPerm.ToError()
    }
    if s.status == StatusReject {
        return erroz.ArticleAlreadyReject.ToError()
    }
    return s.transitionStatus(StatusReject)
}

func (s *Article) Republish() error {
    if !s.actor.CanPublishDirect {
        if !s.actor.CanPublish {
            return erroz.ArticleMissingPublishPerm.ToError()
        }
        if s.status == StatusPendingRepublish {
            return erroz.ArticleAlreadyRepublish.ToError()
        }
    }
    return s.transitionStatus(StatusPendingRepublish)
}

func (s *Article) GetStatus() Status {
    return s.status
}

func (s *Article) GetDescription() string {
    if s.description != "" {
        return s.description
    }
    return s.content[:200]
}
