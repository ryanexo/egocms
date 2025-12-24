package domain

import (
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra/persistence/model`
)

type ArticleStatus int8

const (
    // ArticleDraft 草稿状态
    ArticleDraft ArticleStatus = iota
    // ArticlePending 待审核
    ArticlePending
    // ArticlePublished 已发布
    ArticlePublished
    // ArticleOffline 已下线
    ArticleOffline
    // ArticleReject 审核拒绝
    ArticleReject
    // ArticlePendingRepublish 编辑后等待重新审核
    ArticlePendingRepublish
)

var stateMachine = map[ArticleStatus]struct {
    From  ArticleStatus
    Error error
}{
    ArticlePending:          {From: ArticleDraft, Error: erroz.ArticleSubmitStatusNotAllowed.ToError()},
    ArticlePublished:        {From: ArticlePending, Error: erroz.ArticlePublishStatusNotAllowed.ToError()},
    ArticleOffline:          {From: ArticlePublished, Error: erroz.ArticleOfflineStatusNotAllowed.ToError()},
    ArticleReject:           {From: ArticlePending, Error: erroz.ArticleRejectStatusNotAllowed.ToError()},
    ArticlePendingRepublish: {From: ArticlePublished, Error: erroz.ArticleRepublishStatusNotAllowed.ToError()},
}

type Article struct {
    actor  ArticleActor
    id     uint64
    status ArticleStatus
}

type ArticleActor struct {
    CanSubmit        bool
    CanPublish       bool
    CanPublishDirect bool
    CanOffline       bool
    CanReject        bool
}

func NewArticle(data *model.Article) *Article {
    return &Article{id: data.ID, status: ArticleStatus(data.Status)}
}

func (s *Article) WithActor(actor ArticleActor) {
    s.actor = actor
}

func (s *Article) transitionStatus(expect ArticleStatus) error {
    nextState := stateMachine[expect]
    if nextState.From != s.status {
        return nextState.Error
    }
    s.status = expect
    return nil
}

func (s *Article) Submit() error {
    if s.actor.CanPublishDirect {
        s.status = ArticlePublished
        return nil
    }
    if !s.actor.CanSubmit {
        return erroz.ArticleMissingSubmitPerm.ToError()
    }
    if s.status == ArticlePending {
        return erroz.ArticleAlreadySubmitted.ToError()
    }
    return s.transitionStatus(ArticlePending)
}

func (s *Article) Publish() error {
    if !s.actor.CanPublishDirect {
        if !s.actor.CanPublish {
            return erroz.ArticleMissingPublishPerm.ToError()
        }
        if s.status == ArticlePublished {
            return erroz.ArticleAlreadyPublished.ToError()
        }
    }
    return s.transitionStatus(ArticlePublished)
}

func (s *Article) Offline() error {
    if !s.actor.CanOffline {
        return erroz.ArticleMissingOfflinePerm.ToError()
    }
    if s.status == ArticleOffline {
        return erroz.ArticleAlreadyOffline.ToError()
    }
    return s.transitionStatus(ArticleOffline)
}

func (s *Article) Reject() error {
    if !s.actor.CanReject {
        return erroz.ArticleMissingRejectPerm.ToError()
    }
    if s.status == ArticleReject {
        return erroz.ArticleAlreadyReject.ToError()
    }
    return s.transitionStatus(ArticleReject)
}

func (s *Article) Republish() error {
    if !s.actor.CanPublishDirect {
        if !s.actor.CanPublish {
            return erroz.ArticleMissingPublishPerm.ToError()
        }
        if s.status == ArticlePendingRepublish {
            return erroz.ArticleAlreadyRepublish.ToError()
        }
    }
    return s.transitionStatus(ArticlePendingRepublish)
}

func (s *Article) GetStatus() ArticleStatus {
    return s.status
}
