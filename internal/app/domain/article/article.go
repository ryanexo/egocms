package article

import (
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra/persistence/model`
)

type Article struct {
    actor  Actor
    id     uint64
    status Status
}

type Actor struct {
    CanSubmit        bool
    CanPublish       bool
    CanPublishDirect bool
    CanOffline       bool
    CanReject        bool
}

func NewArticle(data *model.Article) *Article {
    return &Article{id: data.ID, status: Status(data.Status)}
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
        s.status = Published
        return nil
    }
    if !s.actor.CanSubmit {
        return erroz.ArticleMissingSubmitPerm.ToError()
    }
    if s.status == Pending {
        return erroz.ArticleAlreadySubmitted.ToError()
    }
    return s.transitionStatus(Pending)
}

func (s *Article) Publish() error {
    if !s.actor.CanPublishDirect {
        if !s.actor.CanPublish {
            return erroz.ArticleMissingPublishPerm.ToError()
        }
        if s.status == Published {
            return erroz.ArticleAlreadyPublished.ToError()
        }
    }
    return s.transitionStatus(Published)
}

func (s *Article) Offline() error {
    if !s.actor.CanOffline {
        return erroz.ArticleMissingOfflinePerm.ToError()
    }
    if s.status == Offline {
        return erroz.ArticleAlreadyOffline.ToError()
    }
    return s.transitionStatus(Offline)
}

func (s *Article) Reject() error {
    if !s.actor.CanReject {
        return erroz.ArticleMissingRejectPerm.ToError()
    }
    if s.status == Reject {
        return erroz.ArticleAlreadyReject.ToError()
    }
    return s.transitionStatus(Reject)
}

func (s *Article) Republish() error {
    if !s.actor.CanPublishDirect {
        if !s.actor.CanPublish {
            return erroz.ArticleMissingPublishPerm.ToError()
        }
        if s.status == PendingRepublish {
            return erroz.ArticleAlreadyRepublish.ToError()
        }
    }
    return s.transitionStatus(PendingRepublish)
}

func (s *Article) GetStatus() Status {
    return s.status
}
