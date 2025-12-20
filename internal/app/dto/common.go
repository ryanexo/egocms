package srvparams

type QueryByResourceID struct {
    ID uint64 `validator:"required" json:"id"`
}
