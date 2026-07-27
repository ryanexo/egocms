package domain

var (
    // StatusDraft 草稿状态
    StatusDraft = statusValue{v: 0}
    // StatusPending 待审核
    StatusPending = statusValue{v: 1}
    // StatusPublished 已发布
    StatusPublished = statusValue{v: 2}
    // StatusOffline 已下线
    StatusOffline = statusValue{v: 3}
    // StatusReject 审核拒绝
    StatusReject = statusValue{v: 4}
    // StatusPendingRepublish 编辑后等待重新审核
    StatusPendingRepublish = statusValue{v: 5}
)

type statusValue struct {
    v int8
}

func (s statusValue) Value() int8 {
    return s.v
}
