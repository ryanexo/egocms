package domain

import "errors"

type Status int8

const (
	StatusDraft Status = iota
	StatusPending
	StatusPublished
	StatusOffline
	StatusRejected
	StatusPendingRepublish
)

var (
	ErrInvalidStatus        = errors.New("invalid article status")
	ErrInvalidTransition    = errors.New("invalid article status transition")
	ErrMissingCurrentVersion = errors.New("article has no current version")
	ErrNoChangesToRepublish = errors.New("article has no changes to republish")
)

type Workflow struct {
	articleID          uint64
	status             Status
	currentVersionID   uint64
	publishedVersionID *uint64
}

func NewWorkflow(articleID, currentVersionID uint64) (Workflow, error) {
	return RestoreWorkflow(articleID, StatusDraft, currentVersionID, nil)
}

func RestoreWorkflow(articleID uint64, status Status, currentVersionID uint64, publishedVersionID *uint64) (Workflow, error) {
	if !status.Valid() {
		return Workflow{}, ErrInvalidStatus
	}
	if currentVersionID == 0 {
		return Workflow{}, ErrMissingCurrentVersion
	}
	return Workflow{
		articleID:          articleID,
		status:             status,
		currentVersionID:   currentVersionID,
		publishedVersionID: cloneUint64(publishedVersionID),
	}, nil
}

func (status Status) Valid() bool {
	return status >= StatusDraft && status <= StatusPendingRepublish
}

func (workflow Workflow) ArticleID() uint64 {
	return workflow.articleID
}

func (workflow Workflow) Status() Status {
	return workflow.status
}

func (workflow Workflow) CurrentVersionID() uint64 {
	return workflow.currentVersionID
}

func (workflow Workflow) PublishedVersionID() *uint64 {
	return cloneUint64(workflow.publishedVersionID)
}

func (workflow Workflow) CanEdit() error {
	switch workflow.status {
	case StatusDraft, StatusPublished, StatusOffline, StatusRejected:
		return nil
	default:
		return ErrInvalidTransition
	}
}

func (workflow *Workflow) SetCurrentVersion(versionID uint64) error {
	if versionID == 0 {
		return ErrMissingCurrentVersion
	}
	workflow.currentVersionID = versionID
	return nil
}

func (workflow *Workflow) Submit(canPublishDirect bool) error {
	if workflow.status != StatusDraft {
		return ErrInvalidTransition
	}
	if canPublishDirect {
		workflow.publishCurrent()
		return nil
	}
	workflow.status = StatusPending
	return nil
}

func (workflow *Workflow) Publish() error {
	if workflow.status != StatusPending && workflow.status != StatusPendingRepublish {
		return ErrInvalidTransition
	}
	workflow.publishCurrent()
	return nil
}

func (workflow *Workflow) Reject() error {
	if workflow.status != StatusPending && workflow.status != StatusPendingRepublish {
		return ErrInvalidTransition
	}
	workflow.status = StatusRejected
	return nil
}

func (workflow *Workflow) Offline() error {
	if workflow.status != StatusPublished {
		return ErrInvalidTransition
	}
	workflow.status = StatusOffline
	return nil
}

func (workflow *Workflow) Republish(canPublishDirect bool) error {
	switch workflow.status {
	case StatusRejected, StatusOffline:
	case StatusPublished:
		if workflow.publishedVersionID != nil && *workflow.publishedVersionID == workflow.currentVersionID {
			return ErrNoChangesToRepublish
		}
	default:
		return ErrInvalidTransition
	}

	if canPublishDirect {
		workflow.publishCurrent()
		return nil
	}
	workflow.status = StatusPendingRepublish
	return nil
}

func (workflow *Workflow) publishCurrent() {
	workflow.status = StatusPublished
	workflow.publishedVersionID = cloneUint64(&workflow.currentVersionID)
}

func cloneUint64(value *uint64) *uint64 {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}
